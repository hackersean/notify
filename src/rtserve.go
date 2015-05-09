package main

import (
    "fmt"
    "net/http"
    "strings"
    "log"

    "time"
    "os"
)
var vstart,vstop int64
func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()   //解析参数，默认是不会解析的
    fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func test(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()  //解析参数，默认是不会解析的
//    fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
    
    time.Sleep(time.Second*0)
    fmt.Fprintf(w, "receive") //这个写入到w的是输出到客户端的
    
    path:="/proc/uptime"
    fd,err := os.Open(path)  
    if err != nil {
        fmt.Println(path, err)
        return
    }
    defer fd.Close()
    var tmp string    
    fmt.Fscanf(fd,"%s",&tmp)	    
    var tx,ty float32
    fmt.Sscanf(r.Form["mesg"][0],"%f",&tx)
    fmt.Sscanf(tmp,"%f",&ty)
    fmt.Printf("%0.2f %0.2f %0.2f\n",tx,ty,ty-tx)
}

func main() {
    fmt.Println("start")
    http.HandleFunc("/", sayhelloName) //设置访问的路由
    http.HandleFunc("/test/", test) //设置访问的路由
    err :=http.ListenAndServe(":"+os.Args[1], nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}