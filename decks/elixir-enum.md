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

### Enum.count/1
```
@spec count(t()) :: non_neg_integer()
Enum.count([1, 2, 3]) # 3
```






### Enum.map/2 

```
@spec map(t(), (element() -> any())) :: list()
Enum.map([1, 2, 3], fn x -> x * 2 end) # [2, 4, 6]
```

### Enum.filter/2

```
@spec filter(t(), (element() -> as_boolean(term()))) :: list()
Enum.filter([1, 2, 3], fn x -> rem(x, 2) == 0 end) # [2]
```
