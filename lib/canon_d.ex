defmodule CanonD do
  @moduledoc """
  Documentation for `CanonD`.
  """
  alias CanonD.Parser
  alias CanonD.Model.Deck

  def main() do
    priv_dir = :code.priv_dir(:canon_d)
    path = Path.join([priv_dir, "decks", "erlang-lib.md"])
    deck = Parser.parse(path)
    dbg(deck)
    learn_deck(deck)
  end
  
  @spec learn_deck(Deck.t()) :: :ok
  def learn_deck(_deck) do
    # TODO: 
    # - learn card
    :ok
  end
end
