#!/bin/bash

go build -ldflags -s -v -o build/bin/motan-server main/serverdemo.go
go build -ldflags -s -v -o build/bin/motan-client main/clientdemo.go
