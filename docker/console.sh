#!/bin/bash

if [ $# -lt 1 ];
then
    echo $0 '[method]'
    exit
fi

if [ $1 == 'up' ];
then

    if [ $# -lt 2 ];
    then
        echo $0 $1 '[container]'
        exit
    fi

    if [ $2 == 'all' ];
    then
        cd /home/ubuntu/mc2020/logclient/
        docker-compose up -d

        cd /home/ubuntu/mc2020/noisemap/
        docker-compose up -d

        cd /home/ubuntu/mc2020/realworld/server/
        docker-compose up -d

        cd /home/ubuntu/mc2020/realworld/
        docker-compose up -d

        cd /home/ubuntu/mc2020/service/
        docker-compose up -d

        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'logclient' ];
    then
        cd /home/ubuntu/mc2020/logclient/
        docker-compose up -d
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'noisemap' ];
    then
        cd /home/ubuntu/mc2020/noisemap/
        docker-compose up -d
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'realworld_server' ];
    then
        cd /home/ubuntu/mc2020/realworld/server/
        docker-compose up -d
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'realworld' ];
    then
        cd /home/ubuntu/mc2020/realworld/
        docker-compose up -d
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'service' ];
    then
        cd /home/ubuntu/mc2020/service/
        docker-compose up -d
        cd /home/ubuntu/
        exit
    fi

fi

if [ $1 == 'down' ];
then

    if [ $# -lt 2 ];
    then
        echo $0 $1 '[container]'
        exit
    fi

    if [ $2 == 'all' ];
    then
        cd /home/ubuntu/mc2020/logclient/
        docker-compose down -t 1

        cd /home/ubuntu/mc2020/noisemap/
        docker-compose down -t 1

        cd /home/ubuntu/mc2020/realworld/server/
        docker-compose down -t 1

        cd /home/ubuntu/mc2020/realworld/
        docker-compose down -t 1

        cd /home/ubuntu/mc2020/service/
        docker-compose down -t 1

        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'logclient' ];
    then
        cd /home/ubuntu/mc2020/logclient/
        docker-compose down -t 1
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'noisemap' ];
    then
        cd /home/ubuntu/mc2020/noisemap/
        docker-compose down -t 1
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'realworld_server' ];
    then
        cd /home/ubuntu/mc2020/realworld/server/
        docker-compose down -t 1
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'realworld' ];
    then
        cd /home/ubuntu/mc2020/realworld/
        docker-compose down -t 1
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'service' ];
    then
        cd /home/ubuntu/mc2020/service/
        docker-compose down -t 1
        cd /home/ubuntu/
        exit
    fi

fi

if [ $1 == 'rm' ];
then

    if [ $# -lt 2 ];
    then
        echo $0 $1 '[container]'
        exit
    fi

    if [ $2 == 'all' ];
    then
        cd /home/ubuntu/mc2020/logclient/
        docker-compose rm -f

        cd /home/ubuntu/mc2020/noisemap/
        docker-compose rm -f

        cd /home/ubuntu/mc2020/realworld/server/
        docker-compose rm -f

        cd /home/ubuntu/mc2020/realworld/
        docker-compose rm -f

        cd /home/ubuntu/mc2020/service/
        docker-compose rm -f

        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'logclient' ];
    then
        cd /home/ubuntu/mc2020/logclient/
        docker-compose rm -f
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'noisemap' ];
    then
        cd /home/ubuntu/mc2020/noisemap/
        docker-compose rm -f
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'realworld_server' ];
    then
        cd /home/ubuntu/mc2020/realworld/server/
        docker-compose rm -f
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'realworld' ];
    then
        cd /home/ubuntu/mc2020/realworld/
        docker-compose rm -f
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'service' ];
    then
        cd /home/ubuntu/mc2020/service/
        docker-compose rm -f
        cd /home/ubuntu/
        exit
    fi

fi

if [ $1 == 'build' ];
then

    if [ $# -lt 2 ];
    then
        echo $0 $1 '[container]'
        exit
    fi

    if [ $2 == 'all' ];
    then
        cd /home/ubuntu/mc2020/logclient/
        docker rmi mc_logclient
        docker-compose build

        cd /home/ubuntu/mc2020/noisemap/
        docker rmi mc_noisemap
        docker-compose build

        cd /home/ubuntu/mc2020/realworld/server/
        docker rmi mc_realworld_server
        docker-compose build

        cd /home/ubuntu/mc2020/realworld/client/
        docker rmi mc_realworld_client
        docker build -t mc_realworld_client .

        cd /home/ubuntu/mc2020/realworld/
        docker rmi mc_realworld
        docker-compose build

        cd /home/ubuntu/mc2020/service/
        docker rmi mc_service
        docker-compose build

        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'logclient' ];
    then
        cd /home/ubuntu/mc2020/logclient/
        docker rmi mc_logclient
        docker-compose build
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'noisemap' ];
    then
        cd /home/ubuntu/mc2020/noisemap/
        docker rmi mc_noisemap
        docker-compose build
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'realworld_server' ];
    then
        cd /home/ubuntu/mc2020/realworld/server/
        docker rmi mc_realworld_server
        docker-compose build
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'realworld_client' ];
    then
        cd /home/ubuntu/mc2020/realworld/client/
        docker rmi mc_realworld_client
        docker build -t mc_realworld_client .
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'realworld' ];
    then
        cd /home/ubuntu/mc2020/realworld/
        docker rmi mc_realworld
        docker-compose build
        cd /home/ubuntu/
        exit
    fi

    if [ $2 == 'service' ];
    then
        cd /home/ubuntu/mc2020/service/
        docker rmi mc_service
        docker-compose build
        cd /home/ubuntu/
        exit
    fi

fi