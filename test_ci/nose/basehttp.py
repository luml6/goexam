import httplib2
import json
import time
import datetime
import hashlib

class NetUtil:
    def http_get(self,url, headers={}, timeout=5):
        resp_data = None
        try:
            http = httplib2.Http(timeout = timeout)   
            innerheaders = {"Content-type": "application/x-www-form-urlencoded"
                        , "Accept": "text/plain"}
            newheaders = dict( innerheaders.items() + headers.items() )            
            response, content = http.request(url, 'GET', headers=headers)

            print response
            status=response.status
            if 200==status:
                resp_data = json.loads(content)
            else:
                self.errcode=''
                self.errmsg='http response code(%s):%s' % (status,response.reason)
        except Exception, e:
            print e
            resp_data = None

        return resp_data
        
    def http_post(self,url,ps={}, headers={}, timeout=5):
        resp_data = None
        try:
            #print 1
            http = httplib2.Http(timeout = timeout)     
            
            #print 2
            innerheaders = {"Content-type": "application/x-www-form-urlencoded"
                        , "Accept": "text/plain"}
            
            #print 3
            newheaders = dict( innerheaders.items() + headers.items() ) 
            print newheaders

            #print 4
            jdata = json.dumps(ps)

            #print 5
            data = jdata.encode(encoding='UTF8') 

            #print 6
            response, content = http.request(url, 'POST', headers=newheaders, body=data) 

            #print 7
            print response
            status=response.status
            if 200==status:
                resp_data = json.loads(content)
            else:
                self.errcode=''
                self.errmsg='http response code(%s):%s' % (status,response.reason)
        except Exception, e:
            print "http_post CRASH: ",e
            resp_data = None
        return resp_data
    
    def sign(self, values, appkey="bce17962876747028879fc0fe47bbd5a"):
        values["AppKey"] = appkey

        keys = list(values.keys())
        keys.sort()
        string = ""

        for key in keys:
            string += key + str(values[key])
        
        return hashlib.sha1(string.encode(encoding="utf-8")).hexdigest()



if __name__ == '__main__':
    
    url='http://10.0.12.200:10005/demanpoint/demanlist'
    
    m=NetUtil()  
    resp=m.http_get(url)  
    if resp == None:  
        print m.errcode,m.errmsg  
    else:  
        print resp["data"] 
    
    
    m=NetUtil()  
    timestamp = str(int(round(time.time() * 1000)))    
    values = {
        "appId":"2000004944",
        "uid":"75753123",
        "tkt":"1.75753123.2000004944.78ac62778408edad1639755204468a82.1478143668.7760000"
    }
    
    #timestamp = str(int(time.mktime(datetime.datetime.now().timetuple())))
    values["Timestamp"] = timestamp
    headers = {"Content-type": "application/json", "Accept": "text/plain",
                "Timestamp":timestamp,"Sign":m.sign(values)}
          
    
    resp=m.http_post("http://passport.coolyun.com/uac/m/check_tkt",values, headers)
    if resp == None:  
        print m.errcode,m.errmsg  
    else:  
        print resp["rtncode"] 
    
