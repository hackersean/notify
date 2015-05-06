package configure
import (
    "fmt"
    "os" 
    "io"
    "strings"
    "strconv"
)
//服务ID，防止发给自己
var SERVER_ID string=""
var WRITE_MAX_LATE int=1
var CLIENT_MAX_FAIL int=3
var MESSAGE_MAX_LENTGH int=256
var MAX_DEEP int=100000
var QUEUE_MAX_CLIENT int=100
var SESSION_PATH string="cnf/session.cnf"
var FAILOVER_PATH string="cnf/failover.cnf"
var PORT string="9090"

var path="cnf/configure.txt"
func do_read_configure(){
    fd,err := os.Open(path)  
    if err != nil {
        fmt.Println(path, err)
        return
    }
    defer fd.Close()
    var tmp string
    for {
        _,err:=fmt.Fscanf(fd,"%s",&tmp)	
	if err==io.EOF{break}
	if err != nil {	      
              fmt.Println(path, err)
              return
        }
//把文件的配置写入变量中

        tmp_lis:=strings.Split(tmp,"=")
        if(len(tmp_lis)==2){
            switch tmp_lis[0] {
                case "SERVER_ID" :
                    SERVER_ID=tmp_lis[1]
                case "WRITE_MAX_LATE" :
                    WRITE_MAX_LATE,_=strconv.Atoi(tmp_lis[1])
                    break
                case "CLIENT_MAX_FAIL" :
                    CLIENT_MAX_FAIL,_=strconv.Atoi(tmp_lis[1])
                    break
                case "MAX_DEEP" :
                    MAX_DEEP,_=strconv.Atoi(tmp_lis[1])
                    break
                case "QUEUE_MAX_CLIENT" :
                    QUEUE_MAX_CLIENT,_=strconv.Atoi(tmp_lis[1])
                    break
                case "MESSAGE_MAX_LENTGH" :
                    MESSAGE_MAX_LENTGH,_=strconv.Atoi(tmp_lis[1])
                case "SESSION_PATH" :
                    SESSION_PATH=tmp_lis[1]
                    break
                case "FAILOVER_PATH" :
                    FAILOVER_PATH=tmp_lis[1]
                case "PORT" :
                    PORT=tmp_lis[1]
                    break
            }
	    }
	}
}
func init(){
    do_read_configure()
    fmt.Println("import configure ok")
}