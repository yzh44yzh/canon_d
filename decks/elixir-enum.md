# Elixir

## Enum

### Enum.all?/1
```
@spec all?(t()) :: boolean()
Enum.all?([1, nil, 3]) # false
```

### Enum.all?/2
```
@spec all?(t(), (element() -> as_boolean(term()))) :: boolean()
Enum.all?([2, 3, 4], fn x -> rem(x, 2) == 0 end) # false
```

### Enum.any?/1
```
@spec any?(t()) :: boolean()
Enum.any?([false, true, false]) # true
```

### Enum.any?/2
```
@spec any?(t(), (element() -> as_boolean(term()))) :: boolean()
Enum.any?([2, 3, 4], fn x -> rem(x, 2) == 0 end) # true
```

### Enum.at/3
```
@spec at(t(), index(), default()) :: element() | default()
Enum.at([2, 4, 6], 0) # 2
Enum.at([2, 4, 6], 4) # nil
Enum.at([2, 4, 6], 4, :none) # :none
```

### Enum.chunk*
```
Enum.chunk_by/2
Enum.chunk_every/2
Enum.chunk_every/4
Enum.chunk_while/4
```

### Enum.concat/1
```
@spec concat(t()) :: t()
list1 = [1, 2]
list2 = [3, 4]
Enum.concat([list1, list2]) # [1, 2, 3, 4]
```

### Enum.concat/2
```
@spec concat(t(), t()) :: t()
list1 = [1, 2]
list2 = [3, 4]
Enum.concat(list1, list2) # [1, 2, 3, 4]
list1 ++ list2 # [1, 2, 3, 4]
```

### Enum.count*
```
Enum.count/1
Enum.count/2
Enum.count_until/2
Enum.count_until/3
```

### Enum.count/1
```
@spec count(t()) :: non_neg_integer()
Enum.count([1, 2, 3]) # 3
```

### Enum.dedup*
```
Enum.dedup/1
Enum.dedup/2
```

### Enum.dedup/1
```
@spec dedup(t()) :: list()
Enum.dedup([1, 2, 3, 3, 2, 1]) # [1, 2, 3, 2, 1]
```

### Enum.drop*
```
Enum.drop/2
Enum.drop_every/2
Enum.drop_while/2
```

### Enum.drop/2
```
@spec drop(t(), integer()) :: list()
Enum.drop([1, 2, 3], 2) # [3]
Enum.drop([1, 2, 3], -1) # [1, 2]
```

### Enum.drop_while/2
```
@spec drop_while(t(), (element() -> as_boolean(term()))) :: list()
Enum.drop_while([1, 2, 3, 2, 1], fn x -> x < 3 end) # [3, 2, 1]
```

### Enum.each/2
```
@spec each(t(), (element() -> any())) :: :ok
```

### Enum.empty?
```
@spec empty?(t()) :: boolean()
```

### Enum.fetch/2
```
@spec fetch(t(), index()) :: {:ok, element()} | :error
Enum.fetch([2, 4, 6], 0) # {:ok, 2}
Enum.fetch([2, 4, 6], -1) # {:ok, 6}
```

### Enum.fetch!/2
```
@spec fetch!(t(), index()) :: element()
Enum.fetch!([2, 4, 6], 0) # 2
Enum.fetch!([2, 4, 6], -1) # 6
Enum.fetch!([2, 4, 6], 3) # Enum.OutOfBoundsError
```

### Enum.filter/2
```
@spec filter(t(), (element() -> as_boolean(term()))) :: list()
Enum.filter([1, 2, 3], fn x -> rem(x, 2) == 0 end) # [2]
```

### Enum.find*
```
Enum.find/3
Enum.find_index/2
Enum.find_value/2
```

### Enum.find/3
```
find(enumerable, default \\ nil, fun)
@spec find(t(), default(), (element() -> any())) :: element() | default()
Enum.find([2, 3, 4], fn x -> rem(x, 2) == 1 end) # 3
Enum.find([2, 4, 6], fn x -> rem(x, 2) == 1 end) # nil
Enum.find([2, 4, 6], 0, fn x -> rem(x, 2) == 1 end) # 0
```

### Enum.flat_map*
```
Enum.flat_map/2
Enum.flat_map_reduce/3
```

### Enum.flat_map/2
```
@spec flat_map(t(), (element() -> t())) :: list()
Enum.flat_map([:a, :b, :c], fn x -> [x, x] end) # [:a, :a, :b, :b, :c, :c]
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
Enum.group_by(~w{ant buffalo cat dingo}, &String.length/1)
# %{3 => ["ant", "cat"], 5 => ["dingo"], 7 => ["buffalo"]}
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

### Enum.into/2
```
@spec into(Enumerable.t(), Collectable.t()) :: Collectable.t()
Enum.into([a: 2], %{a: 1, b: 3}) # %{a: 2, b: 3}
```

### Enum.join/2
```
@spec join(t(), binary()) :: binary()
Enum.join([1, 2, 3], ",") # "1,2,3"
```

### Enum.map*
```
Enum.map/2
Enum.map_every/2
Enum.map_intersperse/2
Enum.map_join/3
Enum.map_reduce/3
```

### Enum.map/2 
```
@spec map(t(), (element() -> any())) :: list()
Enum.map([1, 2, 3], fn x -> x * 2 end) # [2, 4, 6]
```

### Enum.map_reduce/3
```
map_reduce(enumerable, acc, fun)
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

### Enum.member?/2
```
@spec member?(t(), element()) :: boolean()
```

### Enum.product*
```
Enum.product/1
Enum.product_by/2
```

### Enum.random/1
```
@spec random(t()) :: element()
```

