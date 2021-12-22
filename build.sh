#!/bin/bash
GOARCH=arm CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc GOARM=7 go build -o pirrigo
# JAVA_HOME="/home/joe/java/jdk1.8.0_301/" PATH=$JAVA_HOME:$PATH ANDROID_SDK_ROOT=/home/joe/Android/Sdk cordova build android
JAVA_HOME="/home/joe/java/jdk1.8.0_301/" PATH=/home/joe/.android/flutter/flutter/bin:$JAVA_HOME:$PATH ANDROID_SDK_ROOT=/home/joe/Android/Sdk cordova build android
JAVA_HOME="/home/joe/java/jdk1.8.0_301/" PATH=/home/joe/.android/flutter/flutter/bin:$JAVA_HOME:$PATH ANDROID_SDK_ROOT=/home/joe/Android/Sdk cordova emulate android