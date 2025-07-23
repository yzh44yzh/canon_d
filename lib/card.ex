defmodule CanonD.Card do

  @type t() :: %__MODULE__{
    header: String.t(),
    lines: [String.t()]
  }

  defstruct [:header, :lines]

  @spec new(String.t()) :: t()
  def new(header) do
    header = header |> String.trim("#") |> String.trim()
    %__MODULE__{header: header, lines: []}
  end

  @spec add_line(t(), String.t()) :: t()
  def add_line(card, line) do
    %__MODULE__{card | lines: card.lines ++ [line]}
  end

  @spec has_lines?(t()) :: boolean()
  def has_lines?(card) do
    not Enum.empty?(card.lines)
  end

end
