package tui

import (
	"errors"
	"testing"

	domain "com.github.yzh44yzh/canon_d/internal/model"
	"com.github.yzh44yzh/canon_d/internal/parser"
	"com.github.yzh44yzh/canon_d/internal/session"
	"github.com/google/go-cmp/cmp"
)

func testPractice(t *testing.T, content string) *practice {
	t.Helper()

	s, err := session.NewSession(parser.MakeDeck("test.md", content), 10, nil)
	if err != nil {
		t.Fatal(err)
	}

	return &practice{session: s}
}

func TestPracticeQuitSavesUnfinishedCard(t *testing.T) {
	p := testPractice(t, "# Card\n`one` and `two`")
	var saved []map[domain.CardIdentity]domain.CardStatistics
	p.save = func(s *session.Session) error {
		saved = append(saved, s.Statistics())
		return nil
	}
	p.present(false)
	if len(p.session.Statistics()) != 0 {
		t.Fatal("hidden cards must not count as shows")
	}

	p.present(true)
	p.present(true)
	if got := p.advance("one"); got != answerSubmitted {
		t.Fatalf("submission transition: %v", got)
	}

	if len(saved) != 0 || !p.session.AwaitingContinue() {
		t.Fatal("submission must await acknowledgement without saving")
	}

	p.quit()
	want := []map[domain.CardIdentity]domain.CardStatistics{
		{{Deck: "test.md", Header: "Card"}: {Shows: 1, CorrectAnswers: 1}},
	}
	if diff := cmp.Diff(want, saved); diff != "" {
		t.Fatalf("partial checkpoint (-want +got):\n%s", diff)
	}
}

func TestPracticeCompletedAndRetryCheckpoints(t *testing.T) {
	p := testPractice(t, "# Card\n`answer`")
	var saved []map[domain.CardIdentity]domain.CardStatistics
	p.save = func(s *session.Session) error {
		if !s.Done() {
			t.Fatal("round checkpoint must follow completion")
		}

		saved = append(saved, s.Statistics())
		return nil
	}
	p.present(true)
	if got := p.repeat(); got != practiceUnchanged {
		t.Fatal("unfinished rounds cannot repeat")
	}

	p.advance("wrong")
	if got := p.advance(""); got != roundCompleted || len(saved) != 1 {
		t.Fatal("round completion must checkpoint before returning")
	}

	p.present(true)
	if got := p.advance(""); got != practiceUnchanged || len(saved) != 1 {
		t.Fatal("completed rounds must not advance or checkpoint again")
	}

	if got := p.repeat(); got != retryStarted || len(saved) != 1 {
		t.Fatal("retry must restart missed cards without a checkpoint")
	}

	p.present(true)
	p.present(true)
	p.advance("answer")
	if got := p.advance(""); got != roundCompleted || len(saved) != 2 {
		t.Fatal("retry completion must checkpoint before returning")
	}

	if got := p.repeat(); got != practiceUnchanged {
		t.Fatal("a round without misses cannot repeat")
	}

	p.quit()
	id := domain.CardIdentity{Deck: "test.md", Header: "Card"}
	want := []map[domain.CardIdentity]domain.CardStatistics{
		{id: {Shows: 1, IncorrectAnswers: 1}},
		{id: {Shows: 2, CorrectAnswers: 1, IncorrectAnswers: 1}},
		{id: {Shows: 2, CorrectAnswers: 1, IncorrectAnswers: 1}},
	}
	if diff := cmp.Diff(want, saved); diff != "" {
		t.Fatalf("cumulative checkpoints (-want +got):\n%s", diff)
	}
}

func TestPracticeReportsAcknowledgementAndCardAdvancement(t *testing.T) {
	p := testPractice(t, "# Card\n`one` and `two`")
	p.present(true)
	p.advance("one")
	if p.acceptingAnswer() {
		t.Fatal("typing must be disabled during feedback")
	}

	if got := p.advance("ignored"); got != feedbackAcknowledged || !p.acceptingAnswer() {
		t.Fatal("acknowledgement must expose the next answer without submitting it")
	}

	if p.session.AnswerIndex() != 1 || p.session.AwaitingContinue() {
		t.Fatal("acknowledgement must advance exactly one answer")
	}

	p = testPractice(t, "# One\nReading\n# Two\nMore reading")
	p.present(true)
	if p.acceptingAnswer() || p.advance("") != cardAdvanced {
		t.Fatal("reading cards must advance without answer submission")
	}

	p.present(true)
	if got := p.advance(""); got != roundCompleted {
		t.Fatalf("last reading card must complete the round: %v", got)
	}

	for _, total := range p.session.Statistics() {
		if total != (domain.CardStatistics{Shows: 1}) {
			t.Fatalf("reading card statistics: %+v", total)
		}
	}
}

func TestPracticeSaveFailureBlocksFurtherActions(t *testing.T) {
	for _, completed := range []bool{false, true} {
		name := "partial quit"
		if completed {
			name = "round completion"
		}

		t.Run(name, func(t *testing.T) {
			p := testPractice(t, "# Card\n`answer`")
			want := errors.New("disk unavailable")
			calls := 0
			p.save = func(*session.Session) error {
				calls++
				return want
			}

			if completed {
				p.present(true)
				p.advance("wrong")
				p.advance("")
			} else {
				p.quit()
			}

			before := p.session.Statistics()
			original := p.session
			p.present(true)
			advanced, repeated := p.advance("answer"), p.repeat()
			if advanced != practiceUnchanged || repeated != practiceUnchanged {
				t.Fatal("failed practice must reject further transitions")
			}

			p.start(testPractice(t, "# Other\nReading").session)
			p.quit()
			if !errors.Is(p.err, want) || calls != 1 || p.session != original {
				t.Fatal("failure must block saves and session replacement")
			}

			if p.acceptingAnswer() {
				t.Fatal("failed practice must disable answer input")
			}

			if diff := cmp.Diff(before, p.session.Statistics()); diff != "" {
				t.Fatalf("failed practice changed statistics: %s", diff)
			}
		})
	}
}

func TestPracticeWithoutSessionHasNoEffects(t *testing.T) {
	p := &practice{save: func(*session.Session) error {
		t.Fatal("practice before selection must not checkpoint")
		return nil
	}}
	p.present(true)
	p.quit()
	if p.advance("") != practiceUnchanged || p.repeat() != practiceUnchanged {
		t.Fatal("practice before selection must ignore transitions")
	}

	p.start(testPractice(t, "# Card\nReading").session)
	p.present(true)
	if len(p.session.Statistics()) != 1 {
		t.Fatal("starting selected practice must enable presentation recording")
	}
}
