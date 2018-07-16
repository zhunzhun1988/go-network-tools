#!/bin/bash

#build ssh
echo "start build ssh for linux"
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64   go build -o release/linux/sshhack   ./ssh/main.go
echo "start build ssh for windows"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64   go build -o release/windows/sshhack.exe ./ssh/main.go
echo "start build ssh for mac"
CGO_ENABLED=0 GOOS=darwin   GOARCH=amd64  go build -o release/mac/sshhack   ./ssh/main.go

#build ipscan
echo "start build ipscan for linux"
CGO_ENABLED=0 GOOS=linux   GOARCH=amd64   go build -o release/linux/ipscan   ./ipscan/main.go 
echo "start build ipscan for windows"
CGO_ENABLED=0 GOOS=windows GOARCH=amd64   go build -o release/windows/ipscan.exe ./ipscan/main.go
echo "start build ipscan for mac"
CGO_ENABLED=0 GOOS=darwin   GOARCH=amd64  go build -o release/mac/ipscan   ./ipscan/main.go 
