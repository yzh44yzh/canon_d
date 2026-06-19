defmodule CanonD do

  alias CanonD.Parser
  alias CanonD.Session

  @name "Canon D"
  @version "1.0.0"

  @spec main([String.t()]) :: :ok
  def main(args) do
    case OptionParser.parse(args, options()) do
      {[help: true], [], []} -> help()
      {[version: true], [], []} -> version()
      {params, [deck_file], []} -> 
        mode = case Keyword.get(params, :mode) do
          "repeat" -> :repeat
          "by_heart" -> :by_heart
          _ -> :repeat
        end
        limit = Keyword.get(params, :limit, :unlimited)
      
        deck_file
        |> Parser.parse()
        |> Session.learn_deck(mode, limit)
      _ -> 
        help()
    end
  end

  def options do
    [
      strict: [mode: :string, limit: :integer, version: :boolean, help: :boolean],
      aliases: [m: :mode, l: :limit, v: :version, h: :help]
    ]
  end

  def help() do
    IO.puts("""
    USAGE:
        canon_d [OPTIONS] <path/to/deck.md>
    OPTIONS:
        -m, --mode <M>   Working mode: 'repeat' or 'by_heart'
        -l, --limit <N>  Limit number of cards to learn
        -v, --version    Show version
        -h, --help       Show this help message
    """)
  end

  def version() do
    IO.puts(@name <> " v" <> @version)
  end
end
