defmodule CanonD do
  @moduledoc """
  Documentation for `CanonD`.
  """

  alias CanonD.Card

  def main() do
    priv_dir = :code.priv_dir(:canon_d)
    path = Path.join([priv_dir, "decks", "erlang-lib.md"])
    {:ok, deck} = File.read(path)
    learn_deck(deck)
  end

  @spec learn_deck(String.t()) :: :ok
  def learn_deck(deck) do
    parse_deck(deck)
    |> dbg()
    # TODO: 
    # - learn card
    :ok
  end

  # TODO test
  @spec parse_deck(String.t()) :: [Card.t()]
  def parse_deck(deck) do
    deck
    |> String.split("\n")
    |> Enum.map(&String.trim/1)
    |> Enum.filter(fn line -> line != "" and line != "```" end)
    |> make_cards()
    |> Enum.filter(fn card -> Card.has_lines?(card) end)
  end

  # TODO test
  @spec make_cards([String.t()]) :: [Card.t()]
  def make_cards([first_line | lines]) do
    Enum.reduce(
      lines,
      {[], Card.new(first_line)}, 
      fn line, {cards, curr_card} ->
        if String.starts_with?(line, "#") do
          {[curr_card | cards], Card.new(line)}
        else
          {cards, Card.add_line(curr_card, line)}
        end
      end)
    |> then(fn {cards, last_card} -> [last_card | cards] end)
  end
end
