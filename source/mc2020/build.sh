#!/bin/sh

go build -trimpath -ldflags "-w -extldflags '-static'" . && ~/upx-3.96-amd64_linux/upx -9 mc2020