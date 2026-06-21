### sync go-routines
```
var waitGroup sync.WaitGroup

func main() {
	for i := 0; i < 10; i++ {
		waitGroup.Add(1)
		go actor(i)
	}
	waitGroup.Wait()
}
	
func actor(id int) {
	defer waitGroup.Done()
	fmt.Println("actor", id)
}
```
	
### channels
```
var ch = make(chan int)

func main() {
	for i := 0; i < 10; i++ {
		go actor(i)
	}
	
	for i := 0; i < 10; i++ {
		val := <- ch
		fmt.Println(val)
	}
	close(ch)
}
		
func actor(id int) {
	fmt.Println("actor", id)
	ch <- id
}
```

### channel direction
```
func MyFun(in chan<- int, out <-chan int) {
	val := <- in // read from in
	out <- 42 // write to out
}
```

### select
```
select {
	case res1 := <- ch1:
		fmt.Println("result from ch1", res1)
	case res2 := <- ch2:
		fmt.Println("result from ch2", res2)
	case <- time.After(2 * time.Second):
		fmt.Println("timeout")
}
```
