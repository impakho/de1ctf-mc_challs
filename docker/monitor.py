import socket
import os
import time

CONSOLE_CMD = '/bin/bash /home/ubuntu/mc2020/console.sh'

def CheckConn(ip, port):
    try:
        client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client.settimeout(1)
        client.connect((ip, port))
        client.close()
        return True
    except:
        pass
    return False

while True:
    if not CheckConn('172.20.20.20', 25565):
        os.system(CONSOLE_CMD + ' down service')
        os.system(CONSOLE_CMD + ' up service')
    time.sleep(3)