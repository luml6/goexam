
from basehttp import NetUtil
import category
import random


domain_baseurl = "http://10.0.52.83:50047/API/V1.0/Sync"
AccountID = "76618461"
token="00308835d1bbe351b58418ca1e3bfbc4"

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

    def test_1GetUserUsn(self):
        m=NetUtil()  
        
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_get(domain_baseurl + "/GetUserUsn",headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
       
    def test_2GetSyncNotes(self):
        m=NetUtil()  
        Usn=random.randint(1,50)
        values={
            "Usn":Usn
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/GetSyncNotes",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
        
    def test_3GetSyncCategorys(self):
        m=NetUtil()  
        Usn=random.randint(1,50)
        values={
            "Usn":Usn
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/GetSyncCategorys",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
    
    
    
    
    
    