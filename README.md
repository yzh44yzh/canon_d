# Canon D

> The goal of learning is not to acquire or accumulate knowledge, but to develop the ability to act with understanding.
> Pyotr Yakovlevich Galperin.

> Conceptual knowledge is not enough. You need conviction that comes from your own experience.
> Maxim Dorofeev.

A command-line tool for active [spaced repetition](https://en.wikipedia.org/wiki/Spaced_repetition).

Most spaced repetition systems use "passive" flashcards. This means you simply read a card and then decide how well you remember the information on it. In fact, you need to make an effort to recall it, and that effort is an essential part of spaced repetition.

However, the brain tries to avoid this effort, and "passive" flashcards can easily become a formality with no real benefit. An "active" flashcard requires an action—in this case, typing the missing text. This forces the brain to make an effort and recall the information.


## Meaning of the Name

> Canon (art): a set of norms and rules in art.

> Canon (science): a set of fundamental laws, norms, and methods in scientific knowledge.

Canon D -- Canon for Development -- stands for a programmer's fundamental skills. (_Perhaps this name is too ambitious for the current version. But I plan to develop the application into something more_).

Canon D (Canon in D major) is also a well-known piece by the German composer Johann Pachelbel. Thousands of different covers have been made of it. For example, [this rock version](https://www.youtube.com/watch?v=XZgiNnGB8m4).


## Why Use This When AI Knows Everything, Remembers Everything, and Can Give You All the Answers?

We train our brains much like athletes train their muscles -- through exercises that create an artificial workload. Athletes run and lift weights. Knowledge workers memorize facts and solve problems.

Having a car does not replace going for a run. Having movers does not replace going to the gym. AI does not remove the need to train your natural intelligence.


## Usage

```
$ ./canon_d --deck ./decks/network.md
```

Flashcards are ordinary Markdown files. See [decks/network.md](./decks/network.md).

You can and should create your own flashcards to suit your needs.


## Future Plans

Spaced repetition itself has not been implemented yet, and no information is saved at the end of a session. This will be added later, along with learning progress tracking and reports.

There may also be a multi-user web version, a mobile version, and an Obsidian plugin.
