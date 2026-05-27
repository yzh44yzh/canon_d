compile:
	mix compile
	mix escript.build
	
run:
	./canon_d ./decks/erlang-lists.md
	
test:
	mix test
