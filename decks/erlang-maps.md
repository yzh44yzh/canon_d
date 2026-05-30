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

## maps:groups_from_list/2
```
groups_from_list(KeyFun, List) -> GroupsMap
KeyFun :: fun((Elem) -> Key)
GroupsMap :: #{Key => Group}
```

## maps:is_key/2
```
is_key(Key, Map) -> boolean()
```

## maps:map/2
```
map(Fun, MapOrIter) -> Map
Fun :: fun((Key, Value1) -> Value2)
```

## maps:merge_with/3
```
merge_with(Combiner, Map1, Map2) -> Map3
Combiner :: fun((Key, Value1, Value2) -> Value3)
```

## maps:put/3
```
put(Key, Value, Map1) -> Map2
```

## maps:remove/2
```
remove(Key, Map1) -> Map2
```

## maps:take/2
```
take(Key, Map1) -> {Value, Map2} | error
```

## maps:update/3
```
update(Key, Value, Map1) -> Map2
exception {badkey, Key}
```

## maps:update_with/3
```
update_with(Key, Fun, Map1) -> Map2
Fun :: fun((Value1) -> Value2)
```

## maps:update_with/4
```
update_with(Key, Fun, Init, Map1) -> Map2
Fun :: fun((Value1) -> Value2)
```

## maps:with/2
```
with(Ks, Map1) -> Map2
```

## maps:without/2
```
without(Ks, Map1) -> Map2
```
