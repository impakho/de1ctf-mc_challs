#!/bin/sh

cd /home/logclient

touch log/log.txt
chmod 772 log/log.txt

/bin/sh /cleaner.sh &

while true
do
    su logclient -c 'python app.py'
    sleep 2
done