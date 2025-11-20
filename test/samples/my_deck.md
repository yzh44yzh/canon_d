# My Deck

## Enum

### Enum.map/2 

```
@spec map(t(), (element() -> any())) :: list()
Enum.map([1, 2, 3], fn x -> x * 2 end)
```

### Enum.filter/2

```
@spec filter(t(), (element() -> as_boolean(term()))) :: list()
Enum.filter([1, 2, 3], fn x -> rem(x, 2) == 0 end)
```
