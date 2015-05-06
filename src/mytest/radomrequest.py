#-*- coding:utf-8 -*-
import readline
import urllib
import urllib2
import os
import sys
import threading
import string
import random

def GenPassword(length):
    chars=string.ascii_letters+string.digits
    return ''.join([random.choice(chars) for i in range(length)])#得出的结果中字符会有重复的
    #return ''.join(random.sample(chars, 15))#得出的结果中字符不会有重复的

myport="9090"
#print req
class go_curl(threading.Thread):
    def __init__(self,cyc,qid):
        super(go_curl, self).__init__()
        self.mesg=GenPassword(10)
        self.cyc=cyc
        self.requrl="http://www.kryptosx.info:"+myport+"/queue/"+qid+"/in"
    def run(self):
        for item in range(self.cyc):
            test_data = {'token':'3435234','mesg':self.mesg+str(item)}
            test_data_urlencode = urllib.urlencode(test_data)
            req = urllib2.Request(url = self.requrl,data =test_data_urlencode) 
            res_data = urllib2.urlopen(req)
            res = res_data.read()
            print res
def play(qid,num,cyc):
    curl=[]
    for item in range(num):
        curl.append(go_curl(cyc,qid))
    for item in range(num):
        curl[item].start()
    for item in range(num):
        curl[item].join()

if(len(sys.argv)==2):
    myport=sys.argv[1]
print myport
while True:
    cmmd = raw_input("请输入你的命令(play 队列号 线程 周期) :")
    cmmd_list=cmmd.split()
    #work="qstat"
    if cmmd_list[0]=="play":
        #qid=raw_input("请输入队列号: ")
        #token=raw_input("请输入token: ")
        #qid="2"
        qid=cmmd_list[1]
        num=cmmd_list[2]
        cyc=cmmd_list[3]
        #token="3435234"
        play(qid,int(num),int(cyc))
