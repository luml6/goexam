import os
import time
import json
from contextlib import contextmanager
from exec_cmd import *

class BaseProces():
    def __init__(self, image_domain,port_list,docker_path):
        self._image_domain = image_domain
        self._port_list = port_list
        #Dockerfile must in this path
        self._docker_path = docker_path
        self._image_domain = image_domain
        self._container_id = None
        pass
 
    def __enter__(self):
        print "In  __enter__()"
        return "Foo"
        
    def __exit__(self, type,value, trace):
        print "In__exit__()"
    
    def env_prepare(self):
        print "docker build -t "+self._image_domain+" " + self._docker_path
        process = exec_cmd("docker build -t "+self._image_domain+" " + self._docker_path, timeout=1800,  stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        if process.returncode != 0 :
            print "-------------------create image : std  out--------------------------"
            print process.stdout.read()
            print "-------------------create image : err  out--------------------------"
            print process.stderr.read()
            exit(1)
        print process.stdout.read()
        port_list_str = ""
        for one_port  in self._port_list:
            port_list_str += " -p " + one_port
            
        process = exec_cmd("docker run -d  "+ port_list_str+"  "+self._image_domain, timeout=1800, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        if process.returncode != 0 :
            print "run docker ret : ", process.returncode
            print "-------------------run docker : std  out--------------------------"
            print process.stdout.read()
            print "-------------------run docker : err  out--------------------------"
            print process.stderr.read()
            exit(1)
        docker_container_id = process.stdout.read().strip()
        
        print "container id is : ", docker_container_id
        if len(docker_container_id) != 64:
            print "container id error : ", docker_container_id
            print "-------------------run docker : std  out--------------------------"
            print process.stdout.read()
            print "-------------------run docker : err  out--------------------------"
            print process.stderr.read()
            exit(1)
        self._container_id = docker_container_id
        
        time.sleep(5)
        process = exec_cmd("docker ps | grep "+docker_container_id[0:12], timeout=1800, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
        if process.returncode != 0 or len(process.stdout.read()) == 0:
            print "-------------------docker check : std  out--------------------------"
            print process.stdout.read()
            print "-------------------docker check : err  out--------------------------"
            print process.stderr.read()
            exit(1)
        print process.stdout.read()
    
    @contextmanager
    def acquire_lock(self):
        print "acquire_lock  start"
        try:
            print "acquire_lock  doing"
            yield
            print "acquire_lock  end  doing"
        except Exception,e:
            print "error"
            print e
        finally:
            if self._container_id != None:
                print "----start stop container----"
                process = exec_cmd("docker stop "+self._container_id[0:12], timeout=1800, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                print process.stdout.read()
                print "----start remover container----"
                process = exec_cmd("docker rm "+self._container_id[0:12], timeout=1800, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                print process.stdout.read()
                print "----start remove image "
                process = exec_cmd("docker rmi "+self._image_domain, timeout=1800, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
                print process.stdout.read()
            print "acquire_lock  end"


def CheckCfg(cfgmap):
    if isinstance(cfgmap, dict):
        if cfgmap.has_key("Docker") and cfgmap.has_key("SysTest"):
            return True
        else:
            return False
    return False

def GetHostIP():
    import socket
    hostname = socket.gethostname()
    #print hostname

    ip = socket.gethostbyname(hostname)
    return ip 
    #ipList = socket.gethostbyname_ex(socket.gethostname())
    #print(ipList)

if __name__ == '__main__':

    try:
        cfgfile = file("cfg.json");
        cfgMap = json.load(cfgfile)
        cfgfile.close
        
        
    except Exception, e:
        print "-----Read cfg.json Faild-------"
        print cfgMap
        print "error : ", e
        exit(1)
    if not CheckCfg(cfgMap):
        print "-----Cfg File Error----"
        print cfgMap
        exit(1)
     
    testobj = BaseProces(cfgMap["Docker"]["ImageDomain"], cfgMap["Docker"]["PortList"], cfgMap["Docker"]["DockerFilePath"])
    testobj.env_prepare()
    test_result = True
    
    with testobj.acquire_lock():
        print ("-----system test start-----")
        for onetask in cfgMap["SysTest"]:
            cmd_str = "cd " + onetask["ScriptPath"] + " && " + onetask["RunCmd"] + " && cd - "
            process = exec_cmd(cmd_str,timeout=1800, stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            if process.returncode != 0:
                test_result = False
                print "------execute task "+onetask["TestName"] + " Error-------"
                print process.stderr.read()
                print process.stdout.read()
                break
        print ("-----system test end-----")
        
    if not test_result:
        print "++++++system test error++++++++"
        exit(1)
