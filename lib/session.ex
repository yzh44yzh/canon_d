defmodule CanonD.Session do

  alias CanonD.Model
  alias CanonD.Model.Deck
  alias CanonD.Model.Card
  alias CanonD.Model.LearnState
  
  
  @spec learn_deck(Deck.t(), Model.mode(), integer() | :unlimited) :: :ok
  def learn_deck(%Deck{} = deck, mode, limit \\ :unlimited) do
    IO.puts("Learn deck \"#{deck.name}\", mode: #{mode}, limit: #{limit}")
    
    limit = case limit do
      n when is_integer(n) -> n
      _ -> length(deck.cards)
    end
    
    cards = deck.cards
      |> Enum.shuffle()
      |> Enum.take(limit)
      |> Enum.zip(1..limit)
    
    init_state = %LearnState{mode: mode, total_cards: limit}
    result = Enum.reduce(cards, init_state, &learn_card/2)

    IO.puts("Correct answers: #{result.correct_answers}\nIncorrect answers:#{result.incorrect_answers}")
    ask_for_repeat(result.incorrect_cards, deck, result) 
  end
  
  @spec ask_for_repeat([Card.t()], Deck.t, LearnState.t()) :: :ok
  defp ask_for_repeat([], %Deck{}, _state) do
    :ok
  end

  defp ask_for_repeat(incorrect_cards, %Deck{} = deck, state) do
    case String.trim(IO.gets("Repeat incorrect cards? y/n ")) do
      "y" -> 
        learn_deck(%Deck{deck | cards: incorrect_cards}, state.mode)
      _ -> :ok
    end
  end

  # TODO implement mode: :repeat
  @spec learn_card({Card.t(), integer()}, LearnState.t()) :: LearnState.t()
  def learn_card({%Card{lines: card_lines} = card, idx}, %LearnState{} = state) do
    %LearnState{
      total_cards: tc, 
      correct_answers: ca, 
      incorrect_answers: ia, 
      incorrect_cards: ic
    } = state

    IO.puts("#{idx}/#{tc} | #{card.header}")
    your_lines = read_lines([])
    case card_lines do
      ^your_lines -> 
        IO.puts("Correct\n")
        %LearnState{state | correct_answers: ca + 1}
      _ -> 
        IO.puts("Incorrect")
        IO.puts("Correct Answer:")
        Enum.each(card_lines, fn(line) -> IO.puts("  " <> line) end)
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
end
