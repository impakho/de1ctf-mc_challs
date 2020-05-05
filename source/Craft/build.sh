#!/bin/sh

cmake .
make
~/upx-3.96-amd64_linux/upx -9 craft