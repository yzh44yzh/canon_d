compile:
	mix compile
	mix escript.build
	
run:
	./canon_d priv/decks/erlang-lib.md
	
test:
	mix test
