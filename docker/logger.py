import time

# define
logs = {
    b'logclient': '/var/lib/docker/volumes/de1ctf-logclient-log/_data/log.txt',
    b'realworld_deploy': '/var/lib/docker/volumes/de1ctf-realworld-pipe/_data/pipe.txt',
    b'realworld_server': '/var/lib/docker/volumes/de1ctf-realworld-server-log/_data/log.txt',
    b'realworld_talk': '/var/lib/docker/volumes/de1ctf-realworld-server-log/_data/talklog.txt',
    b'service': '/var/lib/docker/volumes/de1ctf-service-log/_data/log.txt',
    b'ticktock': '/var/lib/docker/volumes/de1ctf-service-log/_data/weblog.txt',
}

while True:
    alls = b''
    for i in logs:
        try:
            read = open(logs[i], 'rb').read()
            open(logs[i], 'wb').write(b'')
            read = read.split(b'\n')
            for j in read:
                if len(j) <= 0:
                    continue
                alls += time.strftime("%Y-%m-%d %H:%M:%S", time.localtime()).encode() + b'[' + i + b'] ' + j + b'\n'
        except:
            pass
    open('log.txt', 'ab').write(alls)
    time.sleep(1)