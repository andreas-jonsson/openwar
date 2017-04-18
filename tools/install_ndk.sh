#!/bin/bash

# curl -L http://dl.google.com/android/ndk/android-ndk-${NDK_VERSION}-linux-x86_64.bin -O
# chmod u+x android-ndk-${NDK_VERSION}-linux-x86_64.bin
# ./android-ndk-${NDK_VERSION}-linux-x86_64.bin > /dev/null

curl -L http://dl.google.com/android/repository/android-ndk-${NDK_VERSION}-linux-x86_64.zip -O
unzip android-ndk-${NDK_VERSION}-linux-x86_64.zip

# rm android-ndk-${NDK_VERSION}-linux-x86_64.bin
rm android-ndk-${NDK_VERSION}-linux-x86_64.zip
export ANDROID_NDK_HOME=`pwd`/android-ndk-${NDK_VERSION}
export PATH=${ANDROID_NDK_HOME}:${PATH}
