# Canon D

Canon D is a flashcard trainer for practicing cards from Markdown decks.

## Language

**Deck**:
A named collection of cards supplied by one Markdown file. A card retains its source
deck when practiced together with cards from other decks.

**Card**:
A learning item with a hierarchical heading and a body containing reading material,
hidden answers, or both.

**Round**:
One pass through a set of card occurrences. A retry round contains the occurrences
missed in the preceding round.

**Card identity**:
The source deck's filename together with the card's full hierarchical header.
Changing either identifies a different card for learning history.

**Answer totals**:
The cumulative numbers of correct and incorrect submitted answers for a card
identity, across sessions and retry rounds.

**Card statistics**:
The cumulative number of shows and answer totals for a card identity, shared across
sessions and retry rounds.

**Show**:
One presentation of a card during practice, including a retry presentation or a
reading card. Redrawing or interacting with the same presentation is not a new show.
