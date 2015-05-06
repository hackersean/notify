package zhttp
import (
    "fmt"
    "net"
    "time"
    "net/http"
    "net/url"
    "zqueue/configure"
//    "io/ioutil"
    
)

//解决超时问题
func HttpRequestNew(deadtime time.Duration) http.Client{
    c := http.Client{
        Transport: &http.Transport{
	    Dial: func(netw, addr string) (net.Conn, error) {
//		    deadline := time.Now().Add(time.Second*3)   //
		    c, err := net.DialTimeout(netw, addr, time.Second*deadtime)    //链接建立的超时
		    if err != nil {
			return nil, err
		    }
//		    c.SetDeadline(deadline)        //送接收数据超时，绝对过期时间。多次请求的话会累计。
		    return c, nil
	    },	    
 //           MaxIdleConnsPerHost :  2,
 //           ResponseHeaderTimeout: time.Second * 2,
        },
	Timeout:time.Second * 3,
    }
    return c

}

//消息推送
func HttpPostForm(http http.Client,push_url string,mesg []string,token string) bool {
//  fmt.Println(push_url,mesg)
    resp, err := http.PostForm(push_url,url.Values{"mesg": mesg,"token":{token},"server_id":{configure.SERVER_ID}})
    
    if err != nil {
        // handle error  
        fmt.Println("fail",err) 
        return false
    }
    

    defer resp.Body.Close()
  //  body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
  // handle error
	return false
    }
 
//    fmt.Println(string(body))
   
    return true
 
}