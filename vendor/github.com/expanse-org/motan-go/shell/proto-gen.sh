#!/bin/bash
#wget https://github.com/google/protobuf/releases/download/v3.5.1/protoc-3.5.1-osx-x86_64.zip
export PATH=$PATH:/software/protoc/bin
protoc -I=./ --go_out=./gengo/ protodef/addressbook.proto
