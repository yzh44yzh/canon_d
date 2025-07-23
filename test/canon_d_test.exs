defmodule CanonDTest do
  use ExUnit.Case
  doctest CanonD

  test "greets the world" do
    assert CanonD.hello() == :world
  end
end
