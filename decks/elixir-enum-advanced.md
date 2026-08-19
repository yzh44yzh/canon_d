# Elixir

## Enum (Advanced)

### Enum.chunk*
```
Enum.chunk_by/2
Enum.chunk_every/2
Enum.chunk_every/4
Enum.chunk_while/4
```

### Enum.count*
```
Enum.count/1
Enum.count/2
Enum.count_until/2
Enum.count_until/3
```

### Enum.dedup*
```
Enum.dedup/1
Enum.dedup_by/2
```

### Enum.drop*
```
Enum.drop/2
Enum.drop_every/2
Enum.drop_while/2
```

### Enum.find*
```
Enum.find/3
Enum.find_index/2
Enum.find_value/2
```

### Enum.flat_map*
```
Enum.flat_map/2
Enum.flat_map_reduce/3
```

### Enum.frequencies*
```
Enum.frequencies/1
Enum.frequencies_by/2
```

### Enum.group_by/3
```
group_by(enumerable, key_fun, value_fun \\ fn x -> x end)
@spec group_by(t(), (element() -> any()), (element() -> any())) :: map()
Enum.group_by(~w{ant buffalo cat dingo}, &String.length/1) # %{3 => ["ant", "cat"], 5 => ["dingo"], 7 => ["buffalo"]}
```

### Enum.intersperse/2
```
@spec intersperse(t(), element()) :: list()
Enum.intersperse([1, 2, 3], 0) # [1, 0, 2, 0, 3]
```

### Enum.into*
```
Enum.into/2
Enum.into/3
```

### Enum.map*
```
Enum.map/2
Enum.map_every/2
Enum.map_intersperse/2
Enum.map_join/3
Enum.map_reduce/3
```

### Enum.map_reduce/3
```
@spec map_reduce(t(), acc(), (element(), acc() -> {element(), acc()})) :: {list(), acc()}
Enum.map_reduce([1, 2, 3], 0, fn x, acc -> {x * 2, x + acc} end) # {[2, 4, 6], 6}
```

### Enum.min* max*
```
Enum.min/3
Enum.min_by/4
Enum.max/3
Enum.max_by/4
Enum.min_max/2
Enum.min_max/3
Enum.min_max_by/4
```

### Enum.product*
```
Enum.product/1
Enum.product_by/2
```

### Enum.reduce*
```
Enum.reduce/2
Enum.reduce/3
Enum.reduce_while/3
```

### Enum.reduce_while/3
```
@spec reduce_while(t(), acc(), (element(), acc() -> {:cont, acc()} | {:halt, acc()})) :: acc()
```

### Enum.reverse*
```
Enum.reverse/1
Enum.reverse/2
Enum.reverse_slice/3
```

### Enum.scan/2
```
@spec scan(t(), (element(), acc() -> acc())) :: list()
Enum.scan(["a", "b", "c"], fn x, acc -> acc <> x end) # ["a", "ab", "abc"]
```

### Enum.scan/3
```
@spec scan(t(), acc(), (element(), acc() -> acc())) :: list()
Enum.scan(["a", "b", "c"], "-", fn x, acc -> acc <> x end) # ["-a", "-ab", "-abc"]
```

### Enum.slice/2
```
@spec slice(t(), Range.t()) :: list()
Enum.slice([1, 2, 3, 4, 5], 1..3) # [2, 3, 4]
Enum.slice([1, 2, 3, 4, 5], 3..10) # [4, 5]
```

### Enum.slice/3
```
@spec slice(t(), index(), non_neg_integer()) :: list()
Enum.slice(1..100, 5, 4) # [6, 7, 8, 9]
```

### Enum.slide/3
```
@spec slide(t(), Range.t() | index(), index()) :: list()
Enum.slide([:a, :b, :c, :d], 2, 1) # [:a, :c, :b, :d]
Enum.slide([:a, :b, :c, :d], 2..4, 0) # [:c, :d, :a, :b]
```

### Enum.sort/3
```
sort_by(enumerable, mapper, sorter \\ :asc)
@spec sort_by(
  t(),
  (element() -> mapped_element),
  (element(), element() -> boolean())
  | :asc
  | :desc
  | module()
  | {:asc | :desc, module()}
) :: list()
when mapped_element: element()
```

### Enum.take_every/2
```
@spec take_every(t(), non_neg_integer()) :: list()
Enum.take_every(1..10, 2) # [1, 3, 5, 7, 9]
```

### Enum.take_random/2
```
@spec take_random(t(), non_neg_integer()) :: list()
Enum.take_random(1..10, 2)
```

### Enum.uniq_by/2
```
@spec uniq_by(t(), (element() -> term())) :: list()
```

