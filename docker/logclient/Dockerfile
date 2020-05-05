FROM python:3.8-slim

EXPOSE 80

# Upgrade & Install packages
RUN echo 'deb http://mirrors.tuna.tsinghua.edu.cn/debian/ buster main contrib non-free' > /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/debian/ buster-updates main contrib non-free' >> /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/debian/ buster-backports main contrib non-free' >> /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/debian-security buster/updates main contrib non-free' >> /etc/apt/sources.list \
    && apt-get update -y \
    && apt-get dist-upgrade -y
RUN apt-get install libcap2-bin iputils-ping procps -y

# Add user
RUN adduser --disabled-password --gecos "" logclient

# Copy files
COPY flag readflag start.sh cleaner.sh /
COPY files/ /home/logclient/

# Set permission
RUN chmod 740 /flag \
    && chmod 755 /readflag \
    && chmod +s /readflag \
    && chmod 770 /start.sh \
    && chmod 770 /cleaner.sh \
    && chown root:root /home/logclient/ \
    && chmod -R 755 /home/logclient/*

# Pip install requirements
RUN pip install -r /home/logclient/requirements.txt -i http://mirrors.tuna.tsinghua.edu.cn/pypi/web/simple/ --trusted-host mirrors.tuna.tsinghua.edu.cn

# Security Strategy
RUN chmod o-x /bin/* /sbin/* /usr/bin/* /usr/local/bin/* \
    && chmod o+x /bin/bash /bin/cat /bin/cp /bin/dash /bin/date /bin/dir /bin/echo /bin/false /bin/hostname /bin/ls /bin/mkdir /bin/mv /bin/ping /bin/ping4 /bin/ping6 /bin/pwd /bin/rm /bin/rmdir /bin/sed /bin/sh /bin/tar /bin/touch /bin/true /bin/uname \
    && chmod o+x /usr/bin/awk /usr/bin/base32 /usr/bin/base64 /usr/bin/cut /usr/bin/perl /usr/bin/tr /usr/bin/wc /usr/bin/whereis /usr/bin/whoami \
    && chmod o+x /usr/local/bin/python /usr/local/bin/python3 /usr/local/bin/python3.8
RUN echo "logclient soft nproc 128"  >> /etc/security/limits.conf \
    && echo "logclient hard nproc 128"  >> /etc/security/limits.conf

# Allow programs to listen ports below 1024
RUN setcap 'cap_net_bind_service=+ep' /usr/local/bin/python3.8

WORKDIR /

ENTRYPOINT ["/start.sh"]