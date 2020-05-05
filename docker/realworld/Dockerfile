FROM python:3.8-alpine

EXPOSE 80

# Upgrade & Install packages
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk upgrade \
    && apk add libcap

# Add user
RUN adduser -D realworld

# Copy files
COPY start.sh /
COPY files/ /home/realworld/

# Set permission
RUN chmod 770 /start.sh \
    && chown root:root /home/realworld/ \
    && chmod -R 755 /home/realworld/*

# Pip install requirements
RUN pip install -r /home/realworld/requirements.txt -i http://mirrors.tuna.tsinghua.edu.cn/pypi/web/simple/ --trusted-host mirrors.tuna.tsinghua.edu.cn

# Allow programs to listen ports below 1024
RUN setcap 'cap_net_bind_service=+ep' /usr/local/bin/python3.8

WORKDIR /

ENTRYPOINT ["/start.sh"]