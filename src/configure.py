#-*- coding:utf-8 -*-
import readline
import urllib
import urllib2
import os
import threading

main_url="http://www.kryptosx.info:"
requrl="http://www.kryptosx.info:9090/queue/1/in"



#print req
def go_curl(myurl,mesg,token):
    test_data = {'token':token,'mesg':mesg}
    test_data_urlencode = urllib.urlencode(test_data)
    req = urllib2.Request(url = myurl,data =test_data_urlencode)
    res_data = urllib2.urlopen(req)
    res = res_data.read()
    return res
            

def get_queue_stat(myurl,token,mesg):
#    print myurl
    ans=go_curl(myurl,mesg,token)
    return ans

token="3435234"
myport="9090"
while True:
    
    print "--------------------------------"
    print "qstat 队列号"
    print "play 队列号 消息内容"
    print "set port/token xxx"
    print "当前myport： "+myport
    print "当前token: "+token
    print "--------------------------------"
    cmmd = raw_input("请输入你的命令: ")
    cmmd_list=cmmd.split()
    #work="qstat"
    if cmmd_list[0]=="qstat":
        #qid=raw_input("请输入队列号: ")
        #token=raw_input("请输入token: ")
        #qid="2"
        qid=cmmd_list[1]
        myurl=main_url+myport+"/queue/"+qid+"/status"
        ans=get_queue_stat(myurl,token,"")
        print ans

    if cmmd_list[0]=="play":
        #qid=raw_input("请输入队列号: ")
        #token=raw_input("请输入token: ")
        #qid="2"
        qid=cmmd_list[1]
        mesg=cmmd_list[2]
#        token="3435234"
        myurl=main_url+myport+"/queue/"+qid+"/in"
        ans=get_queue_stat(myurl,token,mesg)
        print ans
    
    if cmmd_list[0]=="set":
        if cmmd_list[1]=="token":
            token=cmmd_list[2]    
        if cmmd_list[1]=="myport":
            myport=cmmd_list[2]

 



 
