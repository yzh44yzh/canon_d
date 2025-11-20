defmodule CanonD.ParserTest do
  use ExUnit.Case

  alias CanonD.Parser
  alias CanonD.Model.{Card, Deck}

  test "make card" do
    lines = [
      "### maps:filter/2",
      "filter(Pred, MapOrIter) -> Map",
      "Pred = fun((Key, Value) -> boolean()"
    ]
    assert Parser.make_card(lines) == %Card{
      header: "maps:filter/2", 
      lines: [
        "filter(Pred, MapOrIter) -> Map",
        "Pred = fun((Key, Value) -> boolean()"
      ]
    }
  end

  test "group lines by header" do
    lines = [
      "## line 1",
      "line 2",
      "line 3",
      "line 4",
      "## line 5",
      "line 6",
      "line 7",
      "## line 8",
      "line 9",
      "## line 10",
      "## line 11",
      "line 12",
      "line 13",
      "line 14",
      "line 15"
    ]
    assert Parser.group_lines_by_header(lines) == [
      [
        "## line 1",
        "line 2",
        "line 3",
        "line 4"
      ],
      [
        "## line 5",
        "line 6",
        "line 7"
      ],
      [
        "## line 8",
        "line 9"
      ],
      [
        "## line 10"
      ],
      [
        "## line 11",
        "line 12",
        "line 13",
        "line 14",
        "line 15"
      ]
    ]
  end

  test "parse data" do
    data = """
    ## line 1
    line 2

    ```
    line 3
    ```
      line 4

    ## line 5
    line 6
    line 7

    ## line 8
    line 9
    line 10
    line 11
    line 12

    ## line 13
    """
    assert Parser.parse_data(data) == [
      %Card{header: "line 1", lines: ["line 2", "line 3", "line 4"]},
      %Card{header: "line 5", lines: ["line 6", "line 7"]},
      %Card{header: "line 8", lines: ["line 9", "line 10", "line 11", "line 12"]}
    ]
  end

  test "parse" do
    path = File.cwd!() |> Path.join("test/samples/my_deck.md") 
    assert Parser.parse(path) == %Deck{
      name: "my_deck",
      cards: [
        %Card{header: "Enum.map/2", lines: [
          "@spec map(t(), (element() -> any())) :: list()",
          "Enum.map([1, 2, 3], fn x -> x * 2 end)"
        ]},
        %Card{header: "Enum.filter/2", lines: [
          "@spec filter(t(), (element() -> as_boolean(term()))) :: list()",
          "Enum.filter([1, 2, 3], fn x -> rem(x, 2) == 0 end)"
        ]}
      ]
    }
  end
end
