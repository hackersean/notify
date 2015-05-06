package main

import (
    "runtime"
    "fmt"
//    "time"
    "net/http"
    "strings"
    "log"
    "zqueue"
    "failover"
    "zqueue/configure"
)


//-----------------------------------------------------
//URL判断，只是简单的URL判断
func go_simple_check(url_string string) (string,int) {
    url_lis:=strings.Split(url_string,"/")
    if(len(url_lis)==4 && url_lis[0]=="" && url_lis[1]=="queue"){
//        fmt.Printf("go_url %d: %s|%s|%s\n",len(url_lis),url_lis[0],url_lis[1],url_lis[2])
        if(url_lis[3]=="in"){
            return url_lis[2],1
        }else if(url_lis[3]=="status"){
            return url_lis[2],2
        }	
    }else{
        return "",-1
    }
    return "",-1
}
//-------------------------------------------------------
//常规消息展示
func go_index(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()  //解析参数，默认是不会解析的
//    fmt.Println(r.Form)  //这些信息是输出到服务器端的打印信息
/*    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
 */ 
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }    
    //列出当前队列
    if(r.URL.Path=="/list"){
        fmt.Fprintf(w,"当前队列状态：\n队列号|队列堆积 \n")
        ans:=zqueue.Get_Queue_List()
        for _,v:=range(ans){
            fmt.Fprintf(w,v+"\n")
        }
        return
    }
    if(r.URL.Path=="/failover"){
        fmt.Fprintf(w,"容灾列表\n")
        ans:=failover.Get_Adress_List()
 
        fmt.Fprintf(w,ans)
       
        return
    }
    fmt.Fprintf(w, "这是基于http的类notify中间件") //这个写入到w的是输出到客户端的
}

//---------------容灾队列---------------------------------------
func go_failover(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()  //解析参数，默认是不会解析的
 //   fmt.Println(r.Method)
    if r.Method == "GET"{
        go_index(w,r)
        return
    }
    //判断server_id
    if r.Form["server_id"]!=nil && len(r.Form["server_id"])>=1 {
        //排除自循环
        if r.Form["server_id"][0]==configure.SERVER_ID {
            fmt.Fprintf(w, "忽略自己发给自己")
            return
        }
    }
    ans,ok:=failover.Sync_Session_Request(r.Form["token"],r.Form["mesg"])
    if ok==0{
        //返回列表
        for _,v:=range(ans){
            fmt.Fprintf(w,v+"\n")
        }
    }else if ok==1{
        fmt.Fprintf(w, "鉴权失败")
    }
    return
}
//------------------------常规队列-------------------------------------
//队列的格式 /queue/{id}/{do}
func go_queue(w http.ResponseWriter, r *http.Request) {
//    fmt.Println(strings.Split(r.URL.Path,"/"))   
  //  rt := time.Now().UnixNano()

    //定向到go_failover函数去
    if r.URL.Path=="/queue/" {
        go_failover(w,r)
        return
    }
    r.ParseForm()
    var sum int=configure.MESSAGE_MAX_LENTGH
    for _,v:=range(r.Form["mesg"]){
        sum-=len(v)
    }
    if sum<0{
        fmt.Fprintf(w,"消息过长")
        return
    }
    qid,stat:=go_simple_check(r.URL.Path)
    if(stat==-1){
         fmt.Fprintf(w,"url is fail")
    }else if(stat==1){
          //调用函数，发送消息
          ok:=zqueue.Push_Message(qid,r.Form["token"],r.Form["mesg"])             
          if ok==0 {
              fmt.Fprintf(w,"%s ok",qid)
          }else if ok==1{
              fmt.Fprintf(w,"%s token fail",qid)
          }else if ok==2 {
              fmt.Fprintf(w,"队列满")
          }else if ok==3 {
              fmt.Fprintf(w,"%s 不存在",qid)
          }else if ok==4{
              fmt.Fprintf(w,"%s 为保留队列",qid)
          }
    }else if(stat==2){
         stat_str,ok:=zqueue.Get_Status(qid,r.Form["token"])
         if ok==0 {     
             fmt.Fprintf(w,"%s",stat_str)
         }else if ok==1{
             fmt.Fprintf(w,"%s 不存在",qid)
         }else if ok==2{
             fmt.Fprintf(w,"%s 鉴权失败",qid)
         }
         
    }
   // rw := time.Now().UnixNano()

//    fmt.Println(rt,rw,(rw-rt)/1000000)
//    fmt.Fprintf(w, "this is queue") //这个写入到w的是输出到客户端的
}


func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()*2)
    http.HandleFunc("/", go_index) //设置访问的路由
     http.HandleFunc("/queue/failover/in", go_failover) //设置访问的路由
     http.HandleFunc("/queue",go_failover)
    http.HandleFunc("/queue/",go_queue)
    fmt.Printf("the port is: "+configure.PORT+"\n")
    err := http.ListenAndServe(":"+configure.PORT, nil) //设置监听的端口
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}