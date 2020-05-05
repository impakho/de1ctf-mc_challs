FROM node:12-stretch-slim

EXPOSE 80

# Upgrade & Install packages
RUN echo 'deb http://mirrors.tuna.tsinghua.edu.cn/debian/ stretch main contrib non-free' > /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/debian/ stretch-updates main contrib non-free' >> /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/debian/ stretch-backports main contrib non-free' >> /etc/apt/sources.list \
    && echo 'deb http://mirrors.tuna.tsinghua.edu.cn/debian-security stretch/updates main contrib non-free' >> /etc/apt/sources.list \
    && apt-get update -y \
    && apt-get dist-upgrade -y
RUN apt-get install libx11-6 libx11-xcb1 libxcb1 libxcomposite1 libxcursor1 libxdamage1 libxext6 libxfixes3 libxi6 libxrender1 libxtst6 libglib2.0-0 libnss3 libcups2 libdbus-1-3 libxss1 libexpat1 libxrandr2 libasound2 libatk1.0-0 libatk-bridge2.0-0 libpangocairo-1.0-0 libpango-1.0-0 libcairo2 libatspi2.0-0 libgtk-3-0 procps libcap2-bin -y

# Add user
RUN adduser --disabled-password --gecos "" noisemap

# Copy files
COPY start.sh /
COPY files/ /home/noisemap/

# Set permission
RUN chmod 770 /start.sh \
    && chown root:root /home/noisemap/ \
    && chmod -R 755 /home/noisemap/* \
    && chmod -R 757 /home/noisemap/www/maps

# Install cnpm & cnpm install requirements
RUN npm install -g cnpm --registry=http://registry.npm.taobao.org \
    && cd /home/noisemap \
    && cnpm install \
    && cd /

# Allow programs to listen ports below 1024
RUN setcap 'cap_net_bind_service=+ep' /usr/local/bin/node

WORKDIR /

ENTRYPOINT ["/start.sh"]