#!/bin/bash

GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc75-linaro-bin GOARM=7 go build
# arm-linux-gnueabi-gcc GOARM=7 go build
