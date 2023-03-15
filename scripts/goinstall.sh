#!/bin/bash

cd ~
wget https://go.dev/dl/go1.17.5.linux-amd64.tar.gz
tar -zxvf go1.17.5.linux-amd64.tar.gz
rm -rf go1.17.5.linux-amd64.tar.gz
mkdir gopath
cd gopath
mkdir src bin pkg
cd ~
echo "export GOPATH=$HOME/gopath" >> /etc/profile
echo "export PATH=$PATH:$HOME/go/bin:$GOPATH/bin" >> /etc/profile
echo "export GO111MODULE=on" >> /etc/profile
echo "export GOPROXY=https://goproxy.cn" >> /etc/profile
source /etc/profile

go install github.com/cloudflare/cfssl/cmd/cfssl@latest
go install github.com/cloudflare/cfssl/cmd/cfssljson@latest
go install github.com/cloudflare/cfssl/cmd/cfssl-certinfo@latest