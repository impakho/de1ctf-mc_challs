FROM ubuntu:16.04

# Upgrade & Install packages
RUN echo 'deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial main restricted universe multiverse' > /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-updates main restricted universe multiverse' >> /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-backports main restricted universe multiverse' >> /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/ubuntu/ xenial-security main restricted universe multiverse' >> /etc/apt/sources.list \
    && apt-get update -y \
    && apt-get dist-upgrade -y
RUN apt-get install libx11-6 libxrandr2 libxinerama1 libxi6 libxxf86vm1 libxcursor1 libgl1 libcurl3 xorg -y

# Add user
RUN adduser --disabled-password --gecos "" realworld_client

# Copy files
COPY flag start.sh /
COPY files/ /home/realworld_client/

# Set permission
RUN chmod 775 /flag \
    && chmod 770 /start.sh \
    && chmod -R 755 /home/realworld_client/*

WORKDIR /

ENTRYPOINT ["/start.sh"]