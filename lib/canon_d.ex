defmodule CanonD do
  @moduledoc """
  Documentation for `CanonD`.
  """

  def main() do
    priv_dir = :code.priv_dir(:canon_d)
    path = Path.join([priv_dir, "decks", "erlang-lib.md"])
    {:ok, deck} = File.read(path)
    learn_deck(deck)
  end

  @spec learn_deck(String.t()) :: :ok
  def learn_deck(deck) do
    IO.puts(deck)
    # TODO: 
    # - parse cards
    # - learn card
    :ok
  end
end
