#!/bin/bash

rm pirrigo.zip > /dev/null
GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc GOARM=7 go build -o pirrigo
zip -r pirrigo.zip pirrigo templates static init
