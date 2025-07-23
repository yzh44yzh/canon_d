## maps

### maps:filter/2
```
filter(Pred, MapOrIter) -> Map
Pred = fun((Key, Value) -> boolean()
```

### maps:filtermap/2
```
filtermap(Fun, MapOrIter) -> Map
Fun = fun((Key, Value1) -> boolean() | {true, Value2})
```

## lists

### lists:concat/1
```
concat([Thing]) -> string()
Thing = atom() | integer() | float() | string()

> lists:concat([doc, '/', file, '.', 3]).
"doc/file.3"
```

### lists:enumerate/1,2
```
enumerate(List1) -> List2

> lists:enumerate([a,b,c]).
[{1,a},{2,b},{3,c}]

enumerate(Index, List1) -> List2

> lists:enumerate(10, [a,b,c]).
[{10,a},{11,b},{12,c}]
```
