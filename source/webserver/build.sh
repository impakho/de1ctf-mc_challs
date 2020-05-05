#!/bin/sh

go build -trimpath -ldflags "-s -w -extldflags '-static'" . && ~/upx-3.96-amd64_linux/upx -9 webserver