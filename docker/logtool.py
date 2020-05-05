from uuid import uuid4
import atexit
import os
import time

PATH = "/var/lib/docker/volumes/de1ctf-service-logs/_data"

uuid = str(uuid4())

filepath = PATH + "/" + uuid

def exit_handler():
    try:
        os.remove(filepath)
    except:
        pass

atexit.register(exit_handler)

print("Please input your player name (max: 32 bytes): ")
name = input()
if len(name) > 32:
    print("name too long")
    exit()

print("Hello " + name + ".")
print("Your UUID: " + uuid)
print("")

while True:
    message = input("> ")
    if len(message) > 128:
        print("message too long")
        continue
    try:
        message = "%d <%s> %s\n" % (time.time(), name, message)
        open(filepath, 'ab').write(message.encode())
    except:
        pass