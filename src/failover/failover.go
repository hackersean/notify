package failover
import (
    "fmt"
 //   "os" 
//    "io"
 //   "sync"
//    "strings"
//      "strconv"
      "zqueue"
      "time"
 //   "zqueue/configure"
  //  "zqueue/zhttp"
      
)
//容灾队列指针
var client *zqueue.Client

func Get_Session_List() []string{
    var ans []string
    for _,v:=range(zqueue.List()){
  //      fmt.Println(v.Show_what([]string{"qid","delay","token"})+"|"+v.Show_oip())
        ans=append(ans,v.Show_what([]string{"qid","token","delay"})+"|"+v.Show_oip())
    }
    return ans
}

func Auto_Sync(){
    for{
        mesg:=Get_Session_List()
 //       fmt.Println(mesg,"\n--------------------")
        client.Push_Message_In(mesg)
        time.Sleep(3*time.Second)
    }
}

func Get_Adress_List()string{
    tmp_str:=client.Show_oip()
    return tmp_str
}

func Sync_Session_Request(token []string,mesg []string)([]string,int){
    if client.Auth(token)==false{
        return nil,1               //鉴权失败
    } 
    for _,v:=range(mesg){
        if ok:=zqueue.Add(v,true);ok!=0{
            fmt.Println("session error")
        }
    }  
    return Get_Session_List(),0
}

func init(){
    client=zqueue.Find("failover") 
    go Auto_Sync()
    fmt.Println("容灾模块载入完毕")
//     ans:=Get_failover_list()
//     fmt.Println(ans)
}