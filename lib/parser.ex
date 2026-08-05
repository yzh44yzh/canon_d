defmodule CanonD.Parser do
  alias CanonD.Model.{Card, Deck}

  @spec parse(Path.t()) :: Deck.t()
  def parse(path) do
    {:ok, data} = File.read(path)
    name = Path.basename(path, ".md")
    deck = %Deck{name: name, cards: []}
    cards = parse_data(data)
    %Deck{deck | cards: cards}
  end

  @spec parse_data(String.t()) :: [Card.t()]
  def parse_data(data) do
    data
    |> String.split("\n")
    |> Enum.map(&String.trim/1)
    |> Enum.filter(fn line -> line != "" and line != "```" end)
    |> group_lines_by_header()
    |> Enum.filter(fn group -> length(group) > 1 end)
    |> Enum.map(&make_card/1)
  end

  @spec group_lines_by_header([String.t()]) :: [[String.t()]]
  def group_lines_by_header([]), do: []

  def group_lines_by_header([first_line | lines]) do
    Enum.reduce(
      lines,
      {[], [first_line]},
      fn line, {groups, curr_group} ->
        if String.starts_with?(line, "#") do
          {groups ++ [curr_group], [line]}
        else
          {groups, curr_group ++ [line]}
        end
      end
    )
    |> then(fn {groups, last_group} -> groups ++ [last_group] end)
  end

  @spec make_card([String.t()]) :: Card.t()
  def make_card([header | lines]) do
    header = header |> String.trim("#") |> String.trim()
    %Card{header: header, lines: lines}
  end
end
