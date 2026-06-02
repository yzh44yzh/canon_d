# Golang

## switch arg
```
switch arg {
case 0:
	fmt.Println("Zero")
case 1, 2:
	fmt.Println("One Two")
	fallthrough
default:
	fmt.Println("Arg is", arg)
}
```

## switch no arg
```
switch {
case arg == 0:
	fmt.Println("Zero")
case arg > 0:
	fmt.Println("Positive")
case arg < 0:
	fmt.Println("Negative")
}
```

## for
```
for i := 0; i < 10; i++ {
	fmt.Println(i)
}
```

## read stdin
```
fmt.Printf("Enter your name: ")
var name string
fmt.Scanln(&name)
fmt.Println("Your name is", name)
```

## map
```
a := make(map[string]int)
a["key1"] = 1
a["key2"] = 2

b := map[string]int {
    "key1": 1,
    "key2": 2,
} 

v1, ok1 := b["key1"] // 1, true
v2, ok2 := b["key3"] // 0, false
```

## struct
```
type Contact struct {
	Name string
	Surname string
	Phone string
}

var book = []Contact{}
book = append(book, Contact{"Bob", "Bobov", "123"})
```

## slice
```
a := make([]int, 4)
a[0] = 42

b := []int{1, 2, 3, 4}
b = append(b, 5)
```

## iterate slice
```
a := []int{1, 2, 3, 4}
for i, v := range a {
	fmt.Println(i, " ", v)
}
```

## delete element from slice
```
s = append(s[:i], s[i+1:]...)
```

## anonymous func
```
af := func(p int) int {
    return p * 2
}
res := af(42) 
```

## sort.Slice
```
sort.Slice(points, func(i, j int) bool {
    return points[i].X < points[j].X
})
```

## defer
```
f, err := os.Open(filepath)
defer f.Close()
```

## check interface
```
_, ok := interface{}(a).(MyType)
```
