# Erlang lists module

## lists:append/1
```
append(ListOfLists) -> List1

lists:append([[1,2,3], [a,b], [4,5,6]]).
[1,2,3,a,b,4,5,6]
```

## lists:append/2
```
append(List1, List2) -> List3
lists:append("abc", "def").
"abcdef"
```

## lists:concat/1
```
concat([Thing]) -> string()
Thing :: atom() | integer() | float() | string()

lists:concat([doc, '/', file, '.', 3]).
"doc/file.3"
```

## lists:delete/2
```
delete(Elem, List1) -> List2
```

## lists:dropwhile/2
```
dropwhile(Pred, List1) -> List2
```

## lists:duplicate/2
```
duplicate(N, Elem) -> List

lists:duplicate(3, elem).
[elem,elem,elem]
```

## lists:enumerate/1
```
enumerate(List1) -> List2

lists:enumerate([a,b,c]).
[{1,a},{2,b},{3,c}]
```

## lists:enumerate/2
```
enumerate(Index, List1) -> List2

lists:enumerate(10, [a,b,c]).
[{10,a},{11,b},{12,c}]
```

## lists:filter/2
```
filter(Pred, List1) -> List2
Pred :: fun((Elem) -> boolean())
```

## lists:filtermap/2
```
filtermap(Fun, List1) -> List2
Fun :: fun((Elem) -> boolean() | {true, Value})
```

## lists:flatmap/2
```
flatmap(Fun, List1) -> List2
Fun :: fun((A) -> [B])

lists:flatmap(fun(X) -> [X,X] end, [a,b,c]).
[a,a,b,b,c,c]
```

## lists:foldl/3
```
foldl(Fun, Acc0, List) -> Acc1
Fun :: fun((Elem, AccIn) -> AccOut)
```

## lists:join/2
```
join(Sep, List1) -> List2

lists:join(x, [a,b,c]).
[a,x,b,x,c]
```

## lists:map/2
```
map(Fun, List1) -> List2
Fun :: fun((A) -> B)
```

## lists:mapfoldl/3
```
mapfoldl(Fun, Acc0, List1) -> {List2, Acc1}
Fun :: fun((A, AccIn) -> {B, AccOut})
```

## lists:member/2
```
member(Elem, List) -> boolean()
```

## lists:nth/2
```
nth(N, List) -> Elem

lists:nth(3, [a, b, c, d, e]).
c
```

## lists:nthtail/2
```
nthtail(N, List) -> Tail

lists:nthtail(3, [a, b, c, d, e]).
[d,e]
```

## lists:partition/2
```
partition(Pred, List) -> {Satisfying, NotSatisfying}
```

## lists:prefix/2
```
prefix(List1, List2) -> boolean()
```

## lists:search/2
```
search(Pred, List) -> {value, Value} | false
```

## lists:seq/2
```
seq(From, To) -> Seq

lists:seq(1, 5).
[1,2,3,4,5]
```

## lists:seq/3
```
seq(From, To, Incr) -> Seq

lists:seq(1, 10, 3).
[1,4,7,10]
```

## lists:sort/3
```
sort(Fun, List1) -> List2
Fun :: fun((A, B) -> boolean())
true if A <= B
```

## lists:zipwith/3
```
zipwith(Combine, List1, List2) -> List3
Combine :: fun((X, Y) -> T)
```
