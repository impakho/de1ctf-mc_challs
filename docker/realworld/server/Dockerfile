FROM python:2.7-alpine

EXPOSE 4080

# Upgrade & Install packages
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk upgrade

# Add user
RUN adduser -D realworld_server

# Copy files
COPY start.sh /
COPY files/ /home/realworld_server/

# Set permission
RUN chmod 770 /start.sh \
    && chown root:root /home/realworld_server/ \
    && chmod 777 /home/realworld_server/ \
    && chmod -R 755 /home/realworld_server/*

# Pip install requirements
RUN pip install -r /home/realworld_server/requirements.txt -i http://mirrors.tuna.tsinghua.edu.cn/pypi/web/simple/ --trusted-host mirrors.tuna.tsinghua.edu.cn

WORKDIR /

ENTRYPOINT ["/start.sh"]