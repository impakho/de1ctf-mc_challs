#/bin/sh

while true
do
    if [ `ls /proc | wc -l` -gt 168 ]
    then
        killall -STOP -u logclient
        killall -9 -u logclient
        pkill -9 python
        pkill -9 su
        sleep 0.5
    fi
    sleep 0.5
done