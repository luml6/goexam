
from basehttp import NetUtil
import note
import time


domain_baseurl = "http://10.0.52.83:50047/API/V1.0/Share"
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

    def test_1Publish(self):
        n=note.TestUM()
        global notePath
        notePath=n.notecreate()
        m=NetUtil()  
        values={
            "Path": notePath,
            "Type": 4,
            "Data": {
                "Question": "string",
                "Answer": "string"
            }
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Publish",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
    def test_2GetShareMessage(self):
        m=NetUtil()
        values={
            "Path": notePath
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/GetShareMessage",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
        
    def test_6GetQuestion(self):
        m=NetUtil()
        values={
            "Path": notePath
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/GetQuestion",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
		
    def test_3Publish(self):
        n=note.TestUM()
        global notePath
        notePath=n.notecreate()
      
     
        m=NetUtil()  
        values={
            "Path": notePath,
            "Type": 2,
            "Data": {
                "Password": "string"
            }
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Publish",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
    
    def test_4Publish(self):
        n=note.TestUM()
        global notePath
        notePath=n.notecreate()
      
     
        m=NetUtil()  
        values={
            "Path": notePath,
            "Type": 1,
            "Data": {
                
            }
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Publish",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
    
    
    