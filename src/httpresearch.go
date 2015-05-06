package main
import(
    "net"
    "net/http"
    "fmt"
    "time"
    "io/ioutil"
)



func main(){
    c := http.Client{
        Transport: &http.Transport{
	    Dial: func(netw, addr string) (net.Conn, error) {
//	    deadline := time.Now().Add(25 * time.Second)   //
	    c, err := net.DialTimeout(netw, addr, time.Second*1)    //链接建立的超时
	    if err != nil {
	        return nil, err
	    }
//	    c.SetDeadline(deadline)        //送接收数据超时，绝对过期时间。多次请求的话会累计。
	    return c, nil
	    },
        },
    }
    
    fmt.Println("ok",c)
    t := time.Now().Unix()
//    response,err := c.Get("http://192.0.0.1:4040/asdfsadf/asdf")
    response,err := c.Get("http://www.qq.com")
    w := time.Now().Unix()
    fmt.Println(t,w,w-t)
    if(err!=nil){
        fmt.Println(err)
        return
    }
    defer response.Body.Close()
    body,_ := ioutil.ReadAll(response.Body)
    fmt.Println(string(body))
}