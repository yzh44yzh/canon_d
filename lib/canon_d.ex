defmodule CanonD do
  alias CanonD.Parser
  alias CanonD.Model.{Deck, Card, LearnState}

  def main() do
    priv_dir = :code.priv_dir(:canon_d)

    Path.join([priv_dir, "decks", "erlang-lib.md"])
    |> Parser.parse()
    |> learn_deck()
  end
  
  @spec learn_deck(Deck.t()) :: :ok
  def learn_deck(deck) do
    result = Enum.reduce(deck.cards, LearnState.new(), &learn_card/2)
    dbg(result)
    :ok
  end

  @spec learn_card(Card.t(), LearnState.t()) :: LearnState.t()
  def learn_card(card, state) do
    %LearnState{correct_answers: ca, incorrect_answers: _ia, incorrect_cards: _ic} = state
    IO.puts(card.header)
    %LearnState{state | correct_answers: ca + 1}
  end
end
