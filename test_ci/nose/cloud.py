#coding=utf-8
from basehttp import NetUtil
import note
import time
from qiniu import Auth, put_file, etag, urlsafe_base64_encode
import qiniu.config
import uuid
import random
import string





domain_baseurl = "http://10.0.52.83:50047/API/V1.0/Cloud"
AccountID = "76618461"
token="00308835d1bbe351b58418ca1e3bfbc4"
notePath="test"
def setup_module(module):
    print('setup_module every func exec ')

def setup_deco():
    print('setup_deco use with_setup ')

def teardown_deco():
    print('teardown_deco all use  with_setup')


class TestUM():
    
    def setup(self):
        print('setup each fuc this class ')

    @classmethod
    def setup_class(cls):
        print('setup_class use for this class, just one time')

    def test_1Qiniu(self):
        filekey=str(uuid.uuid4())

        n=note.TestUM()
        global notePath
        notePath=n.notecreate()
        fileName=string.join(random.sample(['z','y','x','w','v','u','t','s','r','q','p','o','n','m','l','k','j','i','h','g','f','e','d','c','b','a'],5)).replace(' ','')
        m=NetUtil()  
        values={
            "FileKey": filekey,
            "FileType":"test",
            "FileName":fileName,
            "NotePath":notePath,
            "Expires":3600
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/GetUploadTokenWithCB",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 

        Token = resp["Data"]["UploadToken"]
        print Token
        localfile = './bbb.jpg'

        ret, info = put_file(Token, filekey, localfile)
        print(info)
        print ret
        assert ret['key'] == filekey
        assert ret['hash'] == etag(localfile)
        m=NetUtil()  
        values={
            "FileKey": filekey,
            "Expires":3600
        }
        '''
        七牛回调地址为外网地址，获取回调接口调用外网地址
        '''
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post("http://devnotepad.zeusis.com:50047/API/V1.0/Cloud/GetDownloadURL",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print resp
        assert resp["Code"]==0
    
    
    
    