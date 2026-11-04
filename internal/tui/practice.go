package tui

import "com.github.yzh44yzh/canon_d/internal/session"

// practice coordinates session transitions and checkpoints independently of terminal input.
// Rendering reads the session; all practice actions go through this coordinator.
type practice struct {
	session *session.Session
	save    func(*session.Session) error
	err     error
}

type practiceTransition uint8

const (
	practiceUnchanged practiceTransition = iota
	answerSubmitted
	feedbackAcknowledged
	cardAdvanced
	roundCompleted
	retryStarted
)

func (p *practice) start(s *session.Session) {
	if p.err == nil {
		p.session = s
	}
}

// present accepts display eligibility from the terminal adapter, not terminal dimensions.
func (p *practice) present(visible bool) {
	if visible && p.session != nil && p.err == nil {
		p.session.RecordShow()
	}
}

func (p *practice) acceptingAnswer() bool {
	return p.err == nil && p.session != nil && !p.session.Done() &&
		!p.session.AwaitingContinue() && p.session.AnswerCount() > 0
}

func (p *practice) advance(answer string) practiceTransition {
	if p.err != nil || p.session == nil || p.session.Done() {
		return practiceUnchanged
	}

	if p.acceptingAnswer() {
		p.session.Submit(answer)

		return answerSubmitted
	}

	before := p.session.CardIndex()
	p.session.Continue()
	if p.session.Done() {
		p.checkpoint()

		return roundCompleted
	}

	if p.session.CardIndex() != before {
		return cardAdvanced
	}

	return feedbackAcknowledged
}

func (p *practice) repeat() practiceTransition {
	if p.err != nil || p.session == nil || !p.session.Done() || p.session.MissedCount() == 0 {
		return practiceUnchanged
	}

	p.session.Repeat()

	return retryStarted
}

// quit checkpoints partial-card results as well as completed rounds.
func (p *practice) quit() {
	p.checkpoint()
}

func (p *practice) checkpoint() {
	if p.err == nil && p.session != nil && p.save != nil {
		p.err = p.save(p.session)
	}
}
