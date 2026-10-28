# Bordered terminal UI and semantic color scheme

Canon D uses Bubble Tea for terminal interaction, Bubbles for text input and scrolling, and Lip Gloss for styling. The UI presents the card's context together, separates it from answer entry with rounded panels, and assigns colors by semantic role. This record describes the current implementation as of 2026-10-28.

## Screen layout

The TUI runs in the terminal's alternate screen. During practice, it renders:

1. The deck name and current card number above the panels.
2. A rounded card panel containing the hierarchical card heading and a scrollable body. A horizontal separator joins the side borders below the heading; both the heading and separator stay fixed while the body scrolls.
3. A rounded answer panel containing the input or continuation prompt, followed by the answer counter, reading status, or feedback.
4. Keyboard hints below the panels.

Panels span the terminal width and have thin borders with no added interior padding. At width `W` and height `H`, the card viewport is `max(1, W-2)` columns by `max(1, H-10)` rows; the input width is `max(1, W-4)`. Titles, status, and hints are truncated with an ellipsis when necessary. Card content wraps within the viewport.

The minimum supported size is 20 columns by 11 rows. Below that size, a resize message replaces the layout and practice/selection key handling is suspended. Ctrl+C remains available. Resizing rebuilds the viewport content and brings the active answer's starting row into view. Refreshing after submission or continuation also follows the active answer. Page Up and Page Down allow manual scrolling; changing cards resets scroll position.

The completion screen uses one rounded panel containing totals and, when needed, repeat/quit choices. Its keyboard hints remain below the panel.

## Card rendering and controls

Unrevealed answers appear as `[1: ???]`, `[2: ???]`, and so on, numbered within each card. The current blank uses the active style. Submitting reveals the expected answer with correct or incorrect styling; other answers remain hidden until submitted. Reading cards display their text without an answer field.

Rendering removes backticks, expands each tab into four spaces, and filters control characters. This is display processing, not a Markdown renderer: other Markdown markers remain literal. Source formatting and expected answers stay in the parsed model, as described in the [parsing record](0001-deck-syntax-and-parsing.md).

| Control | Behavior |
|---|---|
| Enter | Submit an answer, acknowledge feedback, or advance a reading card |
| Page Up / Page Down | Scroll the card body |
| Up / Down / Tab / Shift+Tab | Toggle repeat/quit selection on the summary |
| Enter on summary | Activate the selected choice, or quit if there are no misses |
| Ctrl+C | Quit from any screen |

The text input is focused at startup and after continuing or repeating, and uses a virtual cursor. Ordinary letters, including `q`, remain answer input rather than quit shortcuts. The [session](0002-user-session.md) owns answer matching and round transitions.

## Color scheme

`ColorScheme` defines five `color.Color` fields. `DefaultColorScheme` supplies one terminal-palette scheme; styles are built when the UI model is constructed and reused by the input, card renderer, panels, and summary.

| Role | Default palette index | Applied to |
|---|---|---|
| `Border` | 8, gray | All card, answer, and summary borders, including the card header separator |
| `Active` | 6, cyan | Current blank, selected summary choice, input prompt and cursor |
| `Muted` | 8, gray | Keyboard hints, placeholder and suggestions, answer counter, reading status |
| `Correct` | 2, green | Correct revealed answers and feedback |
| `Incorrect` | 1, red | Incorrect revealed answers and feedback |

These indices use the user's terminal palette; their exact appearance depends on that palette. Ordinary card text, typed answers, headings, progress, and backgrounds inherit terminal defaults. Focused and blurred input states use the same role assignments.

Color supplements existing indicators: active blanks and selected choices use bold reverse video, incorrect answers and feedback use underlining, summary selection retains its `>` marker, and feedback explicitly says `Correct` or `Incorrect`. The scheme is defined in code. There is one default instance and no scheme flag, configuration file, runtime switch, or separate panel-specific color roles.

## Implementation

- [UI model, rendering, sizing, and controls](../../internal/tui/model.go)
- [Color roles and derived styles](../../internal/tui/colors.go)
- [Interaction, scrolling, resize, and layout tests](../../internal/tui/model_test.go)
- [Scheme propagation tests](../../internal/tui/colors_test.go)
