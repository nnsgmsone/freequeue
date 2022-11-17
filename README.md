freequeue
==========

This package provies the array based lock free queue with fixed size.

Example
=======
Using the Queue is very simple:
```go
q := queue.New(1024)
for i := 0; i < 1024; i++{
    q.Push(i)
}
for i := 0; i < 1024; i++{
    fmt.Printf("%v\n", q.Pop())
}
```
