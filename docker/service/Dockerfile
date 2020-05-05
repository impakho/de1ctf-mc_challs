FROM ubuntu:16.04

EXPOSE 80 8080 25565

# Upgrade & Install packages
RUN echo 'deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial main restricted universe multiverse' > /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-updates main restricted universe multiverse' >> /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-backports main restricted universe multiverse' >> /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-security main restricted universe multiverse' >> /etc/apt/sources.list \
    && apt-get update -y \
    && apt-get dist-upgrade -y
RUN apt-get install libcap2-bin -y

# Add user
RUN adduser --disabled-password --gecos "" service

# Copy files
COPY start.sh /
COPY files/ /home/service/

# Set permission
RUN chmod 770 /start.sh \
    && chown root:root /home/service/ \
    && chmod -R 755 /home/service/*

# Allow programs to listen ports below 1024
RUN setcap 'cap_net_bind_service=+ep' /home/service/webserver

WORKDIR /

ENTRYPOINT ["/start.sh"]