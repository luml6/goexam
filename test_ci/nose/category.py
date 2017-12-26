from basehttp import NetUtil
import random
import string



domain_baseurl = "http://10.0.52.83:50047/API/V1.0/Category"
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
        
    def test_1categoryAll(self):
        m=NetUtil()  
        
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_get(domain_baseurl + "/All")
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
        

    def test_2categoryAllDeleted(self):
        m=NetUtil()  
        
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                    "AccountID":AccountID,"token":token}
              
        
        resp=m.http_get(domain_baseurl + "/AllDeleted")
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        assert resp["Code"]==0
        
        
    def test_3categorySearch(self):
        m=NetUtil()  
        values =  {
            "Name":"test"
        }

        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Search",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Data"]
        assert resp["Code"]==0
    
    
    def test_4categoryCreate(self):
        m = NetUtil()
        name=string.join(random.sample(['z','y','x','w','v','u','t','s','r','q','p','o','n','m','l','k','j','i','h','g','f','e','d','c','b','a'],5)).replace(' ','')
        values =  {
            "Name":name  
        }

        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Create",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        if resp["Code"] == 0:
            catePath=resp["Data"]["Path"]
            print catePath
        print "aaaaa:",resp["Code"]
        assert resp["Code"]==0
    
    def create(self):
        m = NetUtil()
        name=string.join(random.sample(['z','y','x','w','v','u','t','s','r','q','p','o','n','m','l','k','j','i','h','g','f','e','d','c','b','a'],5)).replace(' ','')
        values =  {
            "Name":name  
        }
        print values
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
        resp=m.http_post(domain_baseurl + "/Create",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg
        print resp
        if resp["Code"] == 0:
            catePath=resp["Data"]["Path"]
            print catePath
        return catePath

    def test_5categoryUpdate(self):
        
        catePath=self.create()
        print catePath 
        m = NetUtil()
        name=string.join(random.sample(['z','y','x','w','v','u','t','s','r','q','p','o','n','m','l','k','j','i','h','g','f','e','d','c','b','a'],5)).replace(' ','')
        values =  {
            "Name":name   
        }
        values["Path"]=catePath
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Update",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Code"]
        assert catePath==resp["Data"]["Path"]
    
    
    
    
    def test_6categoryNotelist(self):
        catePath=self.create()
        m=NetUtil()  
        values =  {
            "Path":catePath    
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/NoteList",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Code"]
        assert resp["Code"]==0
    
    
    def test_7categoryNoteDeletedlist(self):
       
        catePath=self.create()
        print catePath 
        m=NetUtil()  
        values =  {
            "Path":catePath    
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/NoteDeletedList",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Code"]
        assert resp["Code"]==0

    def test_8categoryDelete(self):
        
        catePath=self.create()
        print catePath 
        m=NetUtil()  
        values =  {
            "Path":catePath    
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/Delete",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Code"]
        assert catePath==resp["Data"]["Path"]
    
    
    def test_9categoryDeleteAll(self):
     
        catePath=self.create()
        print catePath 
        m=NetUtil()  
        values =  {
            "Path":catePath    
        }
        headers = {"Content-type": "application/json", "Accept": "text/plain",
                "AccountID":AccountID,"token":token}
              
        
        resp=m.http_post(domain_baseurl + "/DeleteAll",values, headers)
        if resp == None:  
            print m.errCode,m.errmsg  
        else:  
            print resp["Code"] 
        print "aaaaa:",resp["Code"]
        assert catePath==resp["Data"]["Path"]