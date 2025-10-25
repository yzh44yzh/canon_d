defmodule CanonD.Deck do

  alias CanonD.Card
 
  @type t() :: %__MODULE__{
    name: String.t(),
    cards: [Card.t()]
  }
  
  @enforce_keys [:name, :cards]
  defstruct [:name, :cards]

end