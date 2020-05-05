#!/bin/sh

docker network create --driver=bridge --subnet=172.20.0.0/16 de1ctf-mc-net
docker volume create de1ctf-service-log
docker volume create de1ctf-service-logs
docker volume create de1ctf-logclient-log
docker volume create de1ctf-realworld-server-log
docker volume create de1ctf-realworld-pipe
docker pull ubuntu:16.04
docker pull python:3.8-alpine
docker pull python:2.7-alpine
docker pull node:12-alpine

#iptables -t filter -I DOCKER-USER -i br-683519c2e141 ! -o br-683519c2e141 -m state --state NEW -j DROP
#iptables -t filter -I DOCKER-USER -i br-683519c2e141 ! -o br-683519c2e141 -p icmp -m state --state NEW -j ACCEPT