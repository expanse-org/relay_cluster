#!/bin/bash
#AfterInstall

echo "start build.sh......."
WORK_DIR=/opt/loopring/relay
SVC_DIR=/etc/service/relay
GOROOT=/usr/lib/go-1.9
export PATH=$PATH:$GOROOT/bin
export GOPATH=/opt/loopring/go-src

sudo cp -rf $WORK_DIR/src/bin/svc/* $SVC_DIR
sudo chmod -R 755 $SVC_DIR

SRC_DIR=$GOPATH/src/github.com/expanse-org/relay_cluster
if [ ! -d $SRC_DIR ]; then
      sudo mkdir -p $SRC_DIR
	  sudo chown -R ubuntu:ubuntu $GOPATH
fi

cd $SRC_DIR
rm -rf ./*
cp -r $WORK_DIR/src/* ./
go get -u github.com/shopify/sarama
go get -u github.com/expanse-org/relay-lib
go build -ldflags -s -v  -o build/bin/relay cmd/main.go
echo "go build finished......."
cp build/bin/relay $WORK_DIR/bin
