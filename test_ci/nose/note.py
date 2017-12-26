
from basehttp import NetUtil
import category



domain_baseurl = "http://10.0.52.83:50047/API/V1.0/Note"
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
        
    def test_1noCategory(self):
        m=NetUtil()  
        
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_get(domain_baseurl + "/NoCategory")
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
        

    def test_2noteAll(self):
        m=NetUtil()  
        
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_get(domain_baseurl + "/All")
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
        
    def test_3noteGetID(self):
        m=NetUtil()  
        
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_get(domain_baseurl + "/GetNoteID")
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
            
        assert resp["Code"]==0   
    
    def test_4noteCreate(self):
        c=category.TestUM()
        catePath=c.create()
        m=NetUtil()
        values = {
          "Author": "string",
          "Title": "string",
          "Summary": "string",
          "Content": "string",
          "CheckSum": "string",
          "IsPinned": 0
        }
        values["CategoryPath"]=catePath
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
        resp=m.http_post(domain_baseurl + "/Create",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp 
        print "aaaaa:",resp["Data"]
        if resp["Code"]==0:
            global notePath
            notePath=resp["Data"]["Path"]
        assert resp["Code"]==0
        
    def notecreate(self):
        c=category.TestUM()
        catePath=c.create()
        print catePath
        m=NetUtil()
        values = {
          "Author": "string",
          "Title": "string",
          "Summary": "string",
          "Content": "string",
          "CheckSum": "string",
          "IsPinned": 0
        }
        values["CategoryPath"]=catePath
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
        resp=m.http_post(domain_baseurl + "/Create",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        if resp["Code"]==0:
            notePath=resp["Data"]["Path"]
        return notePath
    def test_5noteGet(self):
        m=NetUtil()  
        values =  {
            
        }
        values["Path"]=self.notecreate()

        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Get",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Code"]
        assert resp["Code"]==0
        
    def test_6NoteUpdate(self):
        m = NetUtil()  
        values =  {
          "Author": "test",
          "Title": "string",
          "Summary": "test",
          "CheckSum": "string",
          "Content": "string",
          "IsPinned": 0,
          "UpdateTime": 0
        }
        values["Path"]=self.notecreate()

        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Update",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp 
        print "aaaaa:",resp["Code"]
        assert resp["Code"]==0
    
    
    def test_7noteDelete(self):
        m=NetUtil()  
        values =  {
                
        }
        notePath=self.notecreate()
        values["Path"]=notePath
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Delete",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Code"]
        assert notePath==resp["Data"]["Path"]
    
    def test_8noteGetDeleted(self):
       
        m=NetUtil()  
        values =  {
           
        }
        notePath=self.notecreate()
        values["Path"]=notePath 
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Delete",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Code"]
        assert notePath==resp["Data"]["Path"]
    
    
    def test_9noteRecover(self):
        m=NetUtil()  
        values =  {
           
        }
        notePath=self.notecreate()
        values["Path"]=notePath 
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Delete",values, headers)
        if resp["Code"]==0:
            m=NetUtil()  
            values =  {
            
            }
            values["Path"]=notePath
            headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
            resp=m.http_post(domain_baseurl + "/Recover",values, headers)
            if resp == None:  
                print m.errCode,m.errmsg  
            else:  
                print resp["Code"] 
            print "aaaaa:",resp["Code"]
            assert notePath==resp["Data"]["Path"]
    
    
    def test_10noteMove(self):
       
        c=category.TestUM()
        catePath=c.create()
        print catePath
        m=NetUtil()  
        values =  {
        }
        notePath=self.notecreate()
        values["Path"]=notePath
        values["CategoryPath"]=catePath
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Move",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Code"]
        assert notePath==resp["Data"]["Path"]
    
    
    