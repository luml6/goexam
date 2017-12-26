
from basehttp import NetUtil
import note
import time


domain_baseurl = "http://10.0.52.83:50047/API/V1.0/Vote"
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

    def test_1Vote(self):
        n=note.TestUM()
        global notePath
        notePath=n.notecreate()
        invalidTime=int(round(time.time() * 1000))
        print invalidTime
        m=NetUtil()  
        values={
            "Path": notePath,
            "VoteType": 0,
            "InvalidTime": invalidTime,
            "VoteQuestion": "string",
            "VoteSelect": [
                "string",
                "test1"
            ]
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/VoteAdd",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
    def test_2VoteGetByNote(self):
        m=NetUtil()
        print notePath		
        values={
            "Path": notePath
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/VoteGetByNote",values,headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
        if resp["Code"]==0:
            VoteID=resp["Data"]["VoteID"]
            selectID = resp["Data"]["VoteSelect"][0]["ID"]
            m=NetUtil()  
            values={
                "VoteID": VoteID    
            }
            selects=[]
            selects.append(selectID)
            print selects			
            values["SelectID"]=selects
            print values
            headers = {"Content-type": "application/json", "Accept": "text/plain",
                        "AccountID":AccountID,"token":token}
            resp=m.http_post(domain_baseurl + "/ClickVote",values,headers)
            print resp
            assert resp["Code"]==0
    
    
    
    
    
    
    