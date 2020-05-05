#!/bin/sh

cd /home/realworld_server

touch log/log.txt
touch log/talklog.txt
chmod 772 log/log.txt
chmod 772 log/talklog.txt

su realworld_server -c 'python server.pyc'