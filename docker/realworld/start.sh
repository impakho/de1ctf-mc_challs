#!/bin/sh

cd /home/realworld

touch pipe/pipe.txt
chmod 772 pipe/pipe.txt

su realworld -c 'python app.py'