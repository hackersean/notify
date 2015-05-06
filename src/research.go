package main

import (
    "fmt"
    "sync"
)
var w sync.WaitGroup

func son(){
    for{
        fmt.Println("jhejhe")
    }
}

func father(){
    fmt.Println("father start")
    go son()
    fmt.Println("father end")
    w.Done()
}
func main() {

    w.Add(1)
    father()  
    w.Wait()
    for{
        fmt.Println("testt")
    }
}