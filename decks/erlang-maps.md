# Erlang Lib

## maps:filter/2
```
filter(Pred, MapOrIter) -> Map
Pred = fun((Key, Value) -> boolean())
```

## maps:filtermap/2
```
filtermap(Fun, MapOrIter) -> Map
Fun = fun((Key, Value1) -> boolean() | {true, Value2})
```
