#!/bin/bash


# rm pirri.db
CGO_ENABLED=1 go build
cd static
bower install
cd ../
./pirrigo init/appconfig.json