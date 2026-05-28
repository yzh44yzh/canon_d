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
      {_params, [deck_file], []} -> 
        # TODO apply num_cards limit
        # params: [limit: 10]
      
        deck_file
        |> Parser.parse()
        |> Session.learn_deck()
      _ -> 
        help()
    end
  end

  def options do
    [
      strict: [limit: :integer, version: :boolean, help: :boolean],
      aliases: [l: :limit, v: :version, h: :help]
    ]
  end

  def help() do
    IO.puts("""
    USAGE:
        canon_d [OPTIONS] <path/to/deck.md>
    OPTIONS:
        -l, --limit <N>  Limit number of cards to learn
        -v, --version    Show version
        -h, --help       Show this help message
    """)
  end

  def version() do
    IO.puts(@name <> " v" <> @version)
  end
end
