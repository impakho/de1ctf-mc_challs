#!/bin/sh

cd /home/ubuntu/mc2020/logclient/
docker-compose down -t 1
docker-compose rm -f
docker-compose up -d

cd /home/ubuntu/mc2020/noisemap/
docker-compose down -t 1
docker-compose rm -f
docker-compose up -d

cd /home/ubuntu/mc2020/realworld/server/
docker-compose down -t 1
docker-compose rm -f
docker-compose up -d

cd /home/ubuntu/mc2020/realworld/
docker-compose down -t 1
docker-compose rm -f
docker-compose up -d

cd /home/ubuntu/mc2020/service/
docker-compose down -t 1
docker-compose rm -f
docker-compose up -d

cd /

#*/10 * * * * root /bin/sh /home/ubuntu/mc2020/crontab.sh