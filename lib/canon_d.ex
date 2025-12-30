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
  def learn_deck(deck) do
    IO.puts("Learn deck \"#{deck.name}\"")
    total_cards = length(deck.cards)
    init_state = %LearnState{total_cards: total_cards}
    cards = Enum.zip(1..total_cards, deck.cards)
    result = Enum.reduce(cards, init_state, &learn_card/2)
    IO.puts("Correct answers: #{result.correct_answers}\nIncorrect answers:#{result.incorrect_answers}")
    # TODO: repeate incorrect cards
    :ok
  end

  @spec learn_card(Card.t(), LearnState.t()) :: LearnState.t()
  def learn_card({idx, card}, state) do
    %LearnState{
      total_cards: tc, 
      correct_answers: ca, 
      incorrect_answers: _ia, 
      incorrect_cards: _ic
    } = state
    IO.puts("  Card #{idx}/#{tc}: #{card.header}")
    # TODO: read and check input
    %LearnState{state | correct_answers: ca + 1}
  end
end
