#!/bin/sh

cd /home/service

touch log/log.txt
touch log/weblog.txt
chmod 772 log/log.txt
chmod 772 log/weblog.txt
chmod 757 logs/

su service -c './mc2020 &'
su service -c './webserver 172.20.154.79:4080'