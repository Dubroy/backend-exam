package main

type Employee struct{}

type Item1 struct{}

type Item2 struct{}

type Item3 struct{}

type Item interface {
	// Process 這是一個耗時操作
	Process()
}

func main() {}
