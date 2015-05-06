import urllib
import urllib2
import os
import threading

requrl="http://www.kryptosx.info:9090/queue/1/in"



#print req
class go_curl(threading.Thread):
    def __init__(self, mesg):
        test_data = {'token':'3435234','mesg':mesg}
        test_data_urlencode = urllib.urlencode(test_data)
        req = urllib2.Request(url = requrl,data =test_data_urlencode)
        super(go_curl, self).__init__()
        self.req=req
    def run(self):
        for item in range(100):
            res_data = urllib2.urlopen(self.req)
            res = res_data.read()

num=100
curl=[]
for item in range(num):
    mesg=item
    curl.append(go_curl(mesg))
for item in range(num):
    curl[item].start()
for item in range(num):
    curl[item].join()

print "ok"

#	print res
