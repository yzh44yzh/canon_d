defmodule CanonD.Card do

  @type t() :: %__MODULE__{
    header: String.t(),
    lines: [String.t()]
  }

  @enforce_keys [:header, :lines]
  defstruct [:header, :lines]

end
