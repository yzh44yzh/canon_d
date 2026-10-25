# Golang

## tooling
go `build`
go `fmt`
go `mod`
go `test`
go `vet`

## map
a := `make(map[string]int)`
`a["key1"] = 1`
`a["key2"] = 2`

b := `map[string]int `{
    `"key1": 1`,
    `"key2": 2`,
} 

`v1, ok1 := b["key1"]` // 1, true
`v2, ok2 := b["key3"]` // 0, false

## struct
type Contact struct {
	Name string
	Surname string
	Phone string
}

var book = `[]Contact{}`
book = append(`book, Contact{"Bob", "Bobov", "123"}`)

## slice
a := make(`[]int, 4`)
`a[0] = 42`

b := `[]int{1, 2, 3, 4}`
b = append(`b, 5`)

## delete element from slice
s = append(`s[:i], s[i+1:]...`)

## defer
f, err := `os.Open(filepath)`
defer `f.Close()`

## check interface
_, ok := `interface{}(a).(MyType)`

## bufio.NewReaderSize
reader := `bufio.NewReaderSize(os.Stdin, 1<<20)`
line, _ := `reader.ReadString('\n')`
line = `strings.TrimSpace(line)`
words := `strings.Fields(line)`
n, _ := `strconv.Atoi(words[0])`

## bufio.NewWriterSize
writer := `bufio.NewWriterSize(os.Stdout, 1<<20)`
defer `writer.Flush()`
res := 42
writer.`WriteString(strconv.Itoa(res))`
writer.`WriteByte('\n')`

## Append to file
f, err1 := os.`OpenFile`("/tmp/f.txt", `os.O_APPEND | os.O_CREATE | os.O_WRONLY`, 0o644)
defer f.Close()

n, err2 := f.Write(`[]byte("some data")`)

## Error
func MyDiv(a, b int) `(int, error)` {
	if b == 0 {
		return `0, errors.New("division by zero")`
	}
	return `a / b, nil`
}

## Errorf
fmt.Errorf("`%d` isn't an even number", i)

## string length
import `"unicode/utf8"`
utf8.`RuneCountInString(str)`

## go test
go test `-v ./...`

## new project
mkdir project_name
cd project_name
go `mod init com.github.author/project_name`
`git init`

## type switch
type A struct {
	F int
}
type B struct {
	F string
}

func check(`v interface{}`) bool {
	switch `v2 = v.(type)` {
		case A:
			return true
		case B:
			return true
		default:
			return false
	}
}

## non-comparable types
- `map`
- `slice`
- `function`
- `channel`

## define type
type Score int
type Converter func`(string)Score`
type TeamScores map`[string]Score`

## type assertion
var x interface{} = 7
i := `x.(int)` // i has type int and value 7
v1, ok1 := `x.(int)` // 7, true
v2, ok2 := `x.(float64)` // 0, false

## method
type Person struct {
	FirstName string
	LastName string
	Age int
}

func `(p Person) String() string` {
	return fmt.Sprintf("%s %s, age %d", p.FirstName, p.LastName, p.Age)
}

bob := Person{"Bob", "Bobov", 18)
output := bob.String()

## string interpolation
msg := `fmt.Sprintf`("Hello, `%s`!", name)

## strings builder 
var builder `strings.Builder`
builder.`WriteString`("Hello")
builder.`WriteByte`(' ')
builder.`WriteRune`('世')
result := builder.`String`()
fmt.Println(result) 

