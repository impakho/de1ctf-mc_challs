import os, time, threading

pipe_file = '/var/lib/docker/volumes/de1ctf-realworld-pipe/_data/pipe.txt'

class killThread(threading.Thread):
    def __init__(self, name):
        threading.Thread.__init__(self)
        self.name = name


    def run(self):
        count = 0
        while True:
            if count >= 30:
                for i in range(3):
                    try:
                        os.system("docker kill mc_realworld_client_" + self.name)
                        time.sleep(1)
                        os.system("docker rm -f mc_realworld_client_" + self.name)
                        time.sleep(1)
                    except:
                        continue
                    time.sleep(2)
                break
            count += 1
            time.sleep(1)


def deploy(name):
    os.system('./x11docker --env NAME=' + name + ' --env ADDR=172.20.154.79 --env PORT=4080 --name=mc_realworld_client_' + name + ' --user=root -- --network=de1ctf-mc-net -- mc_realworld_client &')
    time.sleep(5)
    thread = killThread(name)
    thread.start()


while True:
    pipe = open(pipe_file, 'r').read()
    open(pipe_file, 'w').write('')
    if '\n' in pipe:
        for i in pipe.split('\n'):
            if len(i) == 12:
                print(time.time(), i)
                open('deploylog.txt', 'a').write(str(time.time()) + ' ' + i + '\n')
                deploy(i)
    time.sleep(1)