# Erlang maps module

## maps iterator consumed by
```
maps:next/1
maps:filter/2
maps:filtermap/2
maps:fold/3
maps:foreach/2
maps:map/2
maps:to_list/1
```

## maps:filter/2
```
filter(Pred, MapOrIter) -> Map
Pred :: fun((Key, Value) -> boolean())
```

## maps:filtermap/2
```
filtermap(Fun, MapOrIter) -> Map
Fun :: fun((Key, Value1) -> boolean() | {true, Value2})
```

## maps:find/2
```
find(Key, Map) -> {ok, Value} | error
```

## maps:fold/3
```
fold(Fun, Init, MapOrIter) -> Acc
Fun :: fun((Key, Value, AccIn) -> AccOut)
```

## maps:foreach/2
```
foreach(Fun, MapOrIter) -> ok
Fun :: fun((Key, Value) -> term())
```

## maps:from_list/1
```
from_list(List) -> Map
List :: [{Key, Value}]
```

## maps:get/2
```
get(Key, Map) -> Value
exception {badkey, Key}
```

## maps:get/3
```
get(Key, Map, Default) -> Value
```
