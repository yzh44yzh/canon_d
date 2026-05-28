compile:
	mix compile
	mix escript.build
	
run:
	./canon_d ./decks/erlang-lists.md
	
check:
	mix test
	mix dialyzer
