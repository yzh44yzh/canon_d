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

## strings

### strings.EqualFold
```
func EqualFold(s, t string) bool
```

### strings.Index
```
func Index(s, substr string) int
```

### strings.HasPrefix
```
func HasPrefix(s, prefix string) bool
```

### strings.Fields
```
func Fields(s string) []string
```

### strings.Split
```
func Split(s, sep string) []string
```

### strings.Replace
```
func Replace(s, old, new string, n int) string
```

### strings.SplitAfter
```
func SplitAfter(s, sep string) []string
```

### strings.Contains
```
func Contains(s, substr string) bool
```

### strings.TrimSpace
```
func TrimSpace(s string) string
```

### strings.Join
```
func Join(elems []string, sep string) string
```

## slice
```
a := make([]int, 4)
a[0] = 42

b := []int{1, 2, 3, 4}
b = append(b, 5)
```

### iterate slice
```
a := []int{1, 2, 3, 4}
for i, v := range a {
	fmt.Println(i, " ", v)
}
```
