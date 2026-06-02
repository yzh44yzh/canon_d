# Golang

## strings

### strings.Compare
```
func Compare(a, b string) int
```

### strings.Contains
```
func Contains(s, substr string) bool
```

### strings.ContainsAny
```
func ContainsAny(s, chars string) bool
```

### strings.ContainsRune
```
func ContainsRune(s string, r rune) bool
```

### strings.Cut
```
func Cut(s, sep string) (before, after string, found bool)
```

### strings.CutPrefix
```
func CutPrefix(s, prefix string) (after string, found bool)
```

### strings.CutSuffix
```
func CutSuffix(s, suffix string) (before string, found bool)
```

### strings.EqualFold
```
func EqualFold(s, t string) bool
```

### strings.Fields
```
func Fields(s string) []string
```

### strings.HasPrefix
```
func HasPrefix(s, prefix string) bool
```

### strings.HasSuffix
```
func HasSuffix(s, suffix string) bool
```

### strings.Index
```
func Index(s, substr string) int
```

### strings.Join
```
func Join(elems []string, sep string) string
```

### strings.Lines
```
func Lines(s string) iter.Seq[string]
```

### strings.Replace
```
func Replace(s, old, new string, n int) string
```

### strings.ReplaceAll
```
func ReplaceAll(s, old, new string) string
```

### strings.Split
```
func Split(s, sep string) []string
```

### strings.SplitAfter
```
func SplitAfter(s, sep string) []string
```

### strings.Trim
```
func Trim(s, cutset string) string
```

### strings.TrimLeft
```
func TrimLeft(s, cutset string) string
```

### strings.TrimPrefix
```
func TrimPrefix(s, prefix string) string
```

### strings.TrimSpace
```
func TrimSpace(s string) string
```
