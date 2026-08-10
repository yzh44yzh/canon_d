defmodule CanonD.Session do
  alias CanonD.Model
  alias CanonD.Model.Deck
  alias CanonD.Model.Card
  alias CanonD.Model.LearnState

  @spec learn_deck(Deck.t(), Model.mode(), integer()) :: :ok
  def learn_deck(%Deck{} = deck, mode, limit) do
    puts(:blue, "Learn deck \"#{deck.name}\", mode: #{mode}, limit: #{limit}")
    IO.puts(" ")

    cards =
      deck.cards
      |> Enum.shuffle()
      |> Enum.take(limit)
      |> Enum.zip(1..limit)

    init_state = %LearnState{mode: mode, total_cards: min(limit, length(deck.cards))}
    result = Enum.reduce(cards, init_state, &learn_card/2)

    IO.puts(
      "Correct answers: #{result.correct_answers}\nIncorrect answers:#{result.incorrect_answers}"
    )

    ask_for_repeat(result.incorrect_cards, deck, result)
  end

  @spec ask_for_repeat([Card.t()], Deck.t(), LearnState.t()) :: :ok
  defp ask_for_repeat([], %Deck{}, _state) do
    :ok
  end

  defp ask_for_repeat(incorrect_cards, %Deck{} = deck, state) do
    case String.trim(IO.gets("Repeat incorrect cards? y/n ")) do
      "y" ->
        learn_deck(%Deck{deck | cards: incorrect_cards}, state.mode, length(incorrect_cards))

      _ ->
        :ok
    end
  end

  @spec learn_card({Card.t(), integer()}, LearnState.t()) :: LearnState.t()
  def learn_card(
        {%Card{lines: card_lines} = card, idx},
        %LearnState{
          mode: :repeat,
          total_cards: tc,
          correct_answers: ca,
          incorrect_answers: ia,
          incorrect_cards: ic
        } = state
      ) do
    puts(:blue, "#{idx}/#{tc} | #{card.header}")

    if length(card_lines) != 1 do
      Enum.each(card_lines, fn line -> IO.puts("  " <> line) end)
      IO.puts(" ")
      IO.puts("Repeat line by line:")
    end

    correct =
      Enum.reduce(card_lines, true, fn line, acc ->
        IO.puts("  " <> line)
        your_line = String.trim(IO.gets("> "))

        if line == your_line do
          puts(:green, "ok")
          true && acc
        else
          puts(:red, "differ")
          false
        end
      end)

    case correct do
      true ->
        puts(:green, "Correct")
        %LearnState{state | correct_answers: ca + 1}

      false ->
        puts(:red, "Incorrect")
        %LearnState{state | incorrect_answers: ia + 1, incorrect_cards: ic ++ [card]}
    end
  end

  def learn_card(
        {%Card{lines: card_lines} = card, idx},
        %LearnState{
          mode: :by_heart,
          total_cards: tc,
          correct_answers: ca,
          incorrect_answers: ia,
          incorrect_cards: ic
        } = state
      ) do
    puts(:blue, "#{idx}/#{tc} | #{card.header}")

    your_lines = read_lines([])

    case card_lines do
      ^your_lines ->
        puts(:green, "Correct")
        %LearnState{state | correct_answers: ca + 1}

      _ ->
        puts(:red, "Incorrect")
        IO.puts("Correct Answer:")
        Enum.each(card_lines, fn line -> IO.puts("  " <> line) end)
        IO.puts(" ")
        %LearnState{state | incorrect_answers: ia + 1, incorrect_cards: ic ++ [card]}
    end
  end

  @spec read_lines([String.t()]) :: [String.t()]
  def read_lines(acc) do
    case String.trim(IO.gets("> ")) do
      "" -> Enum.reverse(acc)
      line -> read_lines([line | acc])
    end
  end

  @spec puts(atom(), String.t()) :: :ok
  defp puts(color, str) do
    :io_ansi.fwrite([color, str <> "\n"])
    :ok
  end
end
