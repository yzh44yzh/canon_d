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



### Enum.map/2 
```
@spec map(t(), (element() -> any())) :: list()
Enum.map([1, 2, 3], fn x -> x * 2 end) # [2, 4, 6]
```
