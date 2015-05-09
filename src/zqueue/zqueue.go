package zqueue
import (
    "fmt"
    "os" 
    "io"
    "sync"
    "time"
    "strings"
    "strconv"
    "sync/atomic"
    "configure"
    "zqueue/zhttp"
)
//-------------------------------------------------
//-------------------结构体-----------------------------
//-----------------------------------------------
//session结构体
type Client struct{
    qid   string
    token string           //认证
    delay time.Duration     //延迟时间
    dip  chan string      //目的端ip队列
    mpt  map[string]int  //目的端ip散列
    ch    *Queue
    sync.RWMutex           //互斥锁
}
//队列结构体
type Queue struct{
    deep  int32
    pipe  chan []string         //队列
}
//--------------------------------------
//-------------通用函数--------------------
func min(a int,b int)int{
    if a<b{
        return a
    }else{
        return b
    }
}
//-------------队列结构函数--------------------

func (this *Queue) Init(){
    this.deep=0
    this.pipe=make(chan []string,configure.MAX_DEEP+8)
}

func (this *Queue) Push(mesg []string) bool{    
    select {
    	case <-time.After(time.Second* time.Duration(configure.WRITE_MAX_LATE)):
            return false 
        
        case this.pipe<-mesg:
            atomic.AddInt32(&this.deep, 1)
            return true 
    }       
    return false
}

func (this *Queue) Pop() ([]string,bool){   
    mesg,ok:=<-this.pipe
    atomic.AddInt32(&this.deep,-1)
    if ok==false{
        return nil,false
    }else{
        return mesg,true
    }
}
 
//------------------------------------
//----------------客户端session------------------------

//初始化
func (this *Client) Init(){
    this.dip=make(chan string,configure.QUEUE_MAX_CLIENT)
    this.mpt=make(map [string]int)
    //初始化队列
    this.ch=new(Queue)
    this.ch.Init()
}
//添加目的端地址
func (this *Client) Append_obj(ip string,fmax int,force bool) bool {
    //fmt.Println(ip,fmax)
    if fmax<=0 {
        return false
    }
    this.Lock()
    if _,ok:=this.mpt[ip];ok{
        this.Unlock()         //解锁，读的锁——注意
        if force==true{      //强制执行
            this.Reset(ip,fmax)
            return true
        }else{
            return false
        }
    }
    this.dip <- ip
    this.mpt[ip]=fmax
    this.Unlock()             //解锁，函数结束
    return true
}
//重置机器的错误计数，如果被踢出，则加回。
func (this *Client) Reset(ip string,fmax int) bool {
    if fmax<=0 {
        return false
    }
    this.Lock()
    defer this.Unlock()
    if _,ok:=this.mpt[ip];ok==false{
        return false
    }
     //如果机器已经被踢出，需要重新加回channel
    if(this.mpt[ip]<=0){
        this.dip <- ip
    }
    this.mpt[ip]=fmax
    return true
}
func (this *Client) Get_opi_list() []string{
    var ans []string
    this.RLock()
    defer this.RUnlock() 
    for k,v:=range(this.mpt){
        ans=append(ans,k+","+strconv.Itoa(v))
    }
    return ans
}

//输出队列状态     目标IP列表
func (this *Client) Show_oip() string{
    var ans string
    this.RLock()
    defer this.RUnlock()
    for k,v:=range(this.mpt){
        ans=ans+"["+k+","+strconv.Itoa(v)+"];"
    }

    return ans
}
//输出查询的项
func (this *Client) Show_what(keys []string) string{
    //var ans string
    this.RLock()
    defer this.RUnlock()
    var tmp string
    var ans string
    for _,k:=range(keys){
        switch k{
            case "qid" :
                tmp=this.qid
                break
            case "token" :
                tmp=this.token
                break
            //延迟时间
            case "delay" :
                tmp=func (delay time.Duration) string {
                       //纳秒换算成秒
                       var tmp float64
                       tmp=float64(int64(delay)/1000000000)
                       return fmt.Sprintf("%f",tmp) 
                    }(this.delay)
                break   
            default:
                tmp="nil"
        }
        if ans!=""{
            ans=ans+"|"
        }
        ans+=tmp
    }
    return ans
}
//公开信息
func (this *Client) Show_free() string{
    var ans string
    this.RLock()
    defer this.RUnlock()
    ans=this.qid+"|"+strconv.Itoa(int(this.ch.deep))
    return ans
}
//鉴权
func (this *Client)Auth(token []string) bool{
    if len(token)>0 && this.token==token[0]{
        return true
    }else{
        return false
    }
}
//-----------------------------------------------------
func (this *Client) Push_Message_In(mesg []string)bool{
    if this.ch.Push(mesg)==true{
         return true
    }else{
         return false     //队列满了
    }
}
//推送消息到目的端
func (this *Client) Push_Message_Out(){
//    line:=len(this.dip)
//    fmt.Println(line)
    httpd:=zhttp.HttpRequestNew(time.Duration(configure.WRITE_MAX_LATE)*time.Second)
    for {
         //提取消息
         mesg,ok:=this.ch.Pop();
         if ok==false {break}	 
         //选择目的端
         for item := range this.dip { 
               // ch关闭时，for循环会自动结束
               ok=zhttp.HttpPostForm(httpd,item,mesg,this.token)	       
               if ok==true{
                   this.dip <- item
                   break
               }else{
               //===============加锁，保护mpt，写锁====================
                   this.Lock()
                   this.mpt[item]-=1
                   if(this.mpt[item]<=0){
                       delete(this.mpt,item)
                       fmt.Println(item+"挂了",this.mpt[item])
                   }else{		       
                       this.dip <- item
                   }
                   this.Unlock()
               }
         }
//         now := time.Now()         
//         fmt.Println("delay",delay)
         time.Sleep(time.Nanosecond*this.delay)
//         end_time := time.Now()
//         var dur_time time.Duration = end_time.Sub(now)
//         fmt.Println(now,end_time,dur_time)

    }
}

//--------------------session散列----------------------

var queue_session map[string]*Client
var global sync.RWMutex
//判断qid是否存在


//-------------------session 维护，内部接口-------------------------
//列出所有session的指针
func List() []*Client{
    global.RLock()                         //全局锁保护
    defer global.RUnlock()
    var ans []*Client
    for _,v:=range(queue_session){
        ans=append(ans,v)
    }
    return ans
}

//校验
func Check(qid string) bool{
    global.RLock()                         //全局锁保护
    defer global.RUnlock()
    if _,ok:=queue_session[qid];ok{        
        return true
    }else{
        return false
    }
}

//获取session
func Find(qid string) *Client{
    global.RLock()                   //全局锁保护
    defer global.RUnlock()
    if v,ok:=queue_session[qid];ok{
        return v
    }else{
        return v
    }
}

//删除队列
func Delete_Session(qid string,token []string) int{
    tmp:=Find(qid)
    if len(token)>0 && tmp.Auth(token){
        if len(tmp.ch.pipe)==0{
            //关闭channel
            close(tmp.ch.pipe)
            //删除session            
            close(tmp.dip) 
            Del(qid)                     
            return 0
        }else{
            return 2                  //队列不为空
        }
    }
    return 1
}

//添加队列--------------写锁---------------
func Add(session_str string,force bool) int{
   var qid,token,str_delay,oip string
    //----------解析-------------
   session_lis:=strings.Split(session_str,"|")
    if(len(session_lis)!=4){
        return 1     //session不合法
     }else{
         qid=session_lis[0]
         token=session_lis[1]
         str_delay=session_lis[2]
         oip=session_lis[3]
     }
    //-------------------------
    //把秒转换为纳秒
    var client *Client
    global.Lock()               //全局锁保护
    
    client,ok:=queue_session[qid]      //判断session是否存在
//    fmt.Println("sesson",ok)
    if ok==false{
        client=new(Client)      //如果不存在，则初始化一个
        //初始化session
        client.Init()
        //赋值token
        client.token=token
        //读入队列编号
        client.qid=qid
        //延时delay设定,匿名函数
        client.delay=func (str_delay string)time.Duration{var tmp float64
                   fmt.Sscanf(str_delay,"%f",&tmp)
                   //秒 换算成纳秒
                   return time.Duration(tmp*1000000000)
               }(str_delay)
        //把指针保存到map中
        queue_session[qid]=client
    }
    global.Unlock()
    //读入ip列表
    for _,item:=range(strings.Split(oip,";")){   
       if len(item)<=3 {continue}
       item=item[1:len(item)-1] 
       item_lis:=strings.Split(item,",")

       //判断切块
       if(len(item_lis)!=2){continue}
       fmax,ok:=strconv.Atoi(item_lis[1])
       if(ok!=nil){
           fmax=configure.CLIENT_MAX_FAIL
       }else{
           fmax=min(configure.CLIENT_MAX_FAIL,fmax)
       }
//       fmt.Println(item_lis[0],fmax)
       client.Append_obj(item_lis[0],fmax,force)
    }  
    if ok==false{
        go client.Push_Message_Out()
    }
    return 0
}

//删除session结构体--------------写锁--------------------
func Del(qid string) bool{
    global.Lock()               //全局锁保护
    delete(queue_session, qid)  //删除
    global.Unlock() 
    fmt.Println(Find(qid))
    fmt.Println("删除成功")
    return true
}

//-------------------------------------------------------------------------------
//-----------------------------前端管理接口-------------------------------------
//-------------------------------------------------------------------------------
//写入消息
func Push_Message(qid string,token []string,mesg []string) int{
    if qid=="failover"{
        return 4
    }else if Check(qid)==false{
        return 3
    }
    tmp:=Find(qid)
    if len(token)>0 && tmp.Auth(token){       
        //消息入队列,如果写队列失败，直接返回失败
        if tmp.Push_Message_In(mesg)==true{
             return 0
        }else{
             return 2     //队列满了
        }	       
    }else{
        return 1   //鉴权失败
    }
}
//获取队列列表
func Get_Queue_List()[]string{
    global.RLock()
    defer global.RUnlock()
    var ans []string
    for _,v:=range(queue_session){
        ans=append(ans,v.Show_free())
    }
    return ans
}
//获取指定队列状态
func Get_Status(qid string,token []string)(string,int){
    if Check(qid)==false{
        return "",1
    }
    tmp:=Find(qid)
    if len(token)>0 && tmp.Auth(token){
        var ans string
        ans=ans+tmp.Show_what([]string{"qid","delay"})+"|"+tmp.Show_oip()
        return ans,0
    }
    return "",2
}

//---------------------------------------session初始化--------------------------------------------
//读取文件中的session
//文本格式解析


func init_session(){
    fd,err := os.Open(configure.SESSION_PATH)  
    if err != nil {
        fmt.Println(configure.SESSION_PATH, err)
        return
    }
    defer fd.Close()
    var session_str string
    for {
        _,err:=fmt.Fscanf(fd,"%s",&session_str)	
	    if err==io.EOF{break}
	    if err != nil {	      
              fmt.Println(configure.SESSION_PATH, err)
              return
        }
        //把文件的session写入数据结构中
        if(Add(session_str,true)!=0){
             fmt.Println(configure.SESSION_PATH,"session error")          
         }
     //   fmt.Println(line,err)
     }
}




func init(){
     queue_session=make(map[string]*Client)
     init_session()
     fmt.Println("队列模块载入完毕")
}

func main() {
      
}