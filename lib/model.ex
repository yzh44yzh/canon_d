defmodule CanonD.Model do
  defmodule Card do

    @type t() :: %__MODULE__{
      header: String.t(),
      lines: [String.t()]
    }

    @enforce_keys [:header, :lines]
    defstruct [:header, :lines]

  end

  defmodule Deck do

    alias CanonD.Model.Card
   
    @type t() :: %__MODULE__{
      name: String.t(),
      cards: [Card.t()]
    }
    
    @enforce_keys [:name, :cards]
    defstruct [:name, :cards]

  end
end
