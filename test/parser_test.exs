defmodule CanonD.ParserTest do
  use ExUnit.Case

  alias CanonD.Parser
  alias CanonD.Model.Card

  test "make card" do
    lines = [
      "### maps:filter/2",
      "filter(Pred, MapOrIter) -> Map",
      "Pred = fun((Key, Value) -> boolean()"
    ]
    card = Parser.make_card(lines)
    assert card == %Card{
      header: "maps:filter/2", 
      lines: [
        "filter(Pred, MapOrIter) -> Map",
        "Pred = fun((Key, Value) -> boolean()"
      ]
    }
  end
end
