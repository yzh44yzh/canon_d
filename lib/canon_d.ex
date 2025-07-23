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
    parse_deck(deck)
    |> dbg()
    # TODO: 
    # - learn card
    :ok
  end

  # TODO spec and test
  def parse_deck(deck) do
    deck
    |> String.split("\n")
    |> Enum.map(&String.trim/1)
    |> Enum.filter(fn line -> line != "" and line != "```" end)
    |> group_by_header()
    # TODO filter by Card.has_lines?
    |> Enum.filter(fn 
      ([_single_line]) -> false;
      ([_|_]) -> true
    end)
    # TODO not needed
    |> Enum.map(fn lines -> 
      [header | lines] = Enum.reverse(lines)
      {header, lines}
    end)
  end

  # TODO spec and test
  def group_by_header([first_line | lines]) do
    Enum.reduce(
      lines,
      {[], [first_line]}, 
      fn line, {cards, curr_card} ->
        if String.starts_with?(line, "#") do
          # TODO create Card struct
          {[curr_card | cards], [line]}
        else
          # TODO add line to Card struct
          {cards, [line | curr_card]}
        end
      end)
    |> then(fn {cards, last_card} -> [last_card | cards] end)
  end
end
