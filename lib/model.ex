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

  defmodule LearnState do
    alias CanonD.Model.Card

    @type t() :: %__MODULE__{
      correct_answers: non_neg_integer(),
      incorrect_answers: non_neg_integer(),
      incorrect_cards: [Card.t()]
    }

    @enforce_keys [:correct_answers, :incorrect_answers, :incorrect_cards]
    defstruct [:correct_answers, :incorrect_answers, :incorrect_cards]

    def new() do
      %__MODULE__{
        correct_answers: 0,
        incorrect_answers: 0,
        incorrect_cards: []
      }
    end
  end
end
