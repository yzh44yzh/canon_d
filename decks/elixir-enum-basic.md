# Elixir

## Enum (Basic)

### Enum.all?/1
@spec all?(t()) :: boolean()
Enum.all?([1, nil, 3]) # false

### Enum.all?/2
@spec all?(t(), (element() -> as_boolean(term()))) :: boolean()
Enum.all?([2, 3, 4], fn x -> rem(x, 2) == 0 end) # false

### Enum.any?/1
@spec any?(t()) :: boolean()
Enum.any?([false, true, false]) # true

### Enum.any?/2
@spec any?(t(), (element() -> as_boolean(term()))) :: boolean()
Enum.any?([2, 3, 4], fn x -> rem(x, 2) == 0 end) # true

### Enum.at/3
@spec at(t(), index(), default()) :: element() | default()
Enum.at([2, 4, 6], 0) # 2
Enum.at([2, 4, 6], 4) # nil
Enum.at([2, 4, 6], 4, :none) # :none

### Enum.concat/1
@spec concat(t()) :: t()
list1 = [1, 2]
list2 = [3, 4]
Enum.concat([list1, list2]) # [1, 2, 3, 4]

### Enum.concat/2
@spec concat(t(), t()) :: t()
list1 = [1, 2]
list2 = [3, 4]
Enum.concat(list1, list2) # [1, 2, 3, 4]
list1 ++ list2 # [1, 2, 3, 4]

### Enum.count/1
@spec count(t()) :: non_neg_integer()
Enum.count([1, 2, 3]) # 3

### Enum.dedup/1
@spec dedup(t()) :: list()
Enum.dedup([1, 2, 3, 3, 2, 1]) # [1, 2, 3, 2, 1]

### Enum.drop/2
@spec drop(t(), integer()) :: list()
Enum.drop([1, 2, 3], 2) # [3]
Enum.drop([1, 2, 3], -1) # [1, 2]

### Enum.drop_while/2
@spec drop_while(t(), (element() -> as_boolean(term()))) :: list()
Enum.drop_while([1, 2, 3, 2, 1], fn x -> x < 3 end) # [3, 2, 1]

### Enum.each/2
@spec each(t(), (element() -> any())) :: :ok

### Enum.empty?
@spec empty?(t()) :: boolean()

### Enum.fetch/2
@spec fetch(t(), index()) :: {:ok, element()} | :error
Enum.fetch([2, 4, 6], 0) # {:ok, 2}
Enum.fetch([2, 4, 6], -1) # {:ok, 6}

### Enum.fetch!/2
@spec fetch!(t(), index()) :: element()
Enum.fetch!([2, 4, 6], 0) # 2
Enum.fetch!([2, 4, 6], -1) # 6
Enum.fetch!([2, 4, 6], 3) # Enum.OutOfBoundsError

### Enum.filter/2
@spec filter(t(), (element() -> as_boolean(term()))) :: list()
Enum.filter([1, 2, 3], fn x -> rem(x, 2) == 0 end) # [2]

### Enum.find/3
find(enumerable, default \\ nil, fun)
@spec find(t(), default(), (element() -> any())) :: element() | default()
Enum.find([2, 3, 4], fn x -> rem(x, 2) == 1 end) # 3
Enum.find([2, 4, 6], fn x -> rem(x, 2) == 1 end) # nil
Enum.find([2, 4, 6], 0, fn x -> rem(x, 2) == 1 end) # 0

### Enum.flat_map/2
@spec flat_map(t(), (element() -> t())) :: list()
Enum.flat_map([:a, :b, :c], fn x -> [x, x] end) # [:a, :a, :b, :b, :c, :c]

### Enum.into/2
@spec into(Enumerable.t(), Collectable.t()) :: Collectable.t()
Enum.into([a: 2], %{a: 1, b: 3}) # %{a: 2, b: 3}

### Enum.join/2
@spec join(t(), binary()) :: binary()
Enum.join([1, 2, 3], ",") # "1,2,3"

### Enum.map/2 
@spec map(t(), (element() -> any())) :: list()
Enum.map([1, 2, 3], fn x -> x * 2 end) # [2, 4, 6]

### Enum.member?/2
@spec member?(t(), element()) :: boolean()

### Enum.random/1
@spec random(t()) :: element()

### Enum.reduce/2
@spec reduce(t(), (element(), acc() -> acc())) :: acc()
Enum.reduce([1, 2, 3], fn x, acc -> x + acc end) # 6
Enum.reduce([], fn x, acc -> x + acc end) # Enum.EmptyError

### Enum.reduce/3
@spec reduce(t(), acc(), (element(), acc() -> acc())) :: acc()
Enum.reduce([1, 2, 3], 0, fn x, acc -> x + acc end) # 6

### Enum.reject/2
@spec reject(t(), (element() -> as_boolean(term()))) :: list()
Enum.reject([1, 2, 3], fn x -> rem(x, 2) == 0 end) # [1, 3]

### Enum.reverse/1
@spec reverse(t()) :: list()

### Enum.shuffle/1
@spec shuffle(t()) :: list()

### Enum.sort/1
@spec sort(t()) :: list()
Enum.sort([3, 2, 1]) # [1, 2, 3]

### Enum.sort/2
@spec sort(t(), (element(), element() -> boolean()) | :asc | :desc | more...) :: list()
Enum.sort([1, 2, 3], &(&1 >= &2)) # [3, 2, 1]
Enum.sort([2, 3, 1], :desc) # [3, 2, 1]

### Enum.split/2
@spec split(t(), integer()) :: {list(), list()}
Enum.split([1, 2, 3], 2) # {[1, 2], [3]}

### Enum.split_while/2
@spec split_while(t(), (element() -> as_boolean(term()))) :: {list(), list()}
Enum.split_while([1, 2, 3, 4], fn x -> x < 3 end) # {[1, 2], [3, 4]}

### Enum.split_with/2
@spec split_with(t(), (element() -> as_boolean(term()))) :: {list(), list()}
Enum.split_with([1, 2, 3, 4], fn x -> rem(x, 2) == 0 end) # {[2, 4], [1, 3]}

### Enum.sum/1
@spec sum(t()) :: number()

### Enum.sum_by/2
@spec sum_by(t(), (element() -> number())) :: number()

### Enum.take/2
@spec take(t(), integer()) :: list()
Enum.take([1, 2, 3], 2) # [1, 2]

