defmodule CanonD do
  alias CanonD.Parser
  alias CanonD.Model.{Deck, Card, LearnState}

  def main() do
    priv_dir = :code.priv_dir(:canon_d)

    # TODO get path from args
    Path.join([priv_dir, "decks", "erlang-lib.md"])
    |> Parser.parse()
    |> learn_deck()
  end
  
  @spec learn_deck(Deck.t()) :: :ok
  def learn_deck(%Deck{} = deck) do
    IO.puts("Learn deck \"#{deck.name}\"")
    total_cards = length(deck.cards)
    init_state = %LearnState{total_cards: total_cards}
    # TODO shuffle cards
    cards = Enum.zip(1..total_cards, deck.cards)
    result = Enum.reduce(cards, init_state, &learn_card/2)
    IO.puts("Correct answers: #{result.correct_answers}\nIncorrect answers:#{result.incorrect_answers}")
    # TODO: repeate incorrect cards
    :ok
  end

  @spec learn_card(Card.t(), LearnState.t()) :: LearnState.t()
  def learn_card({idx, %Card{lines: card_lines} = card}, %LearnState{} = state) do
    %LearnState{
      total_cards: tc, 
      correct_answers: ca, 
      incorrect_answers: ia, 
      incorrect_cards: ic
    } = state

    IO.puts("  Card #{idx}/#{tc}: #{card.header}")
    your_lines = read_lines([])
    case card_lines do
      ^your_lines -> 
        IO.puts("  Correct\n")
        %LearnState{state | correct_answers: ca + 1}
      _ -> 
        IO.puts("  Incorrect\n")
        IO.puts("  Correct Answer:")
        Enum.each(card_lines, fn(line) -> IO.puts("    " <> line) end)
        %LearnState{state | incorrect_answers: ia + 1, incorrect_cards: ic ++ [card]}
    end
  end

  @spec read_lines([String.t()]) :: [String.t()]
  def read_lines(acc) do
    case String.trim(IO.gets("  > ")) do
      "" -> Enum.reverse(acc)
      line -> read_lines([line | acc])
    end
  end
end
