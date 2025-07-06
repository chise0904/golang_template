#!/bin/bash

# Install required dependencies
apt-get update && apt-get install -y \
    curl \
    git \
    make \
    binutils \
    bison \
    gcc \
    build-essential

# Install Go directly first
curl -OL https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
export PATH=$PATH:/usr/local/go/bin

# Install GVM
bash < <(curl -s -S -L https://raw.githubusercontent.com/moovweb/gvm/master/binscripts/gvm-installer)
source /root/.gvm/scripts/gvm

# Install Go 1.21 using GVM
gvm install go1.23.10
gvm use go1.23.10 --default

# Set up GOPATH
mkdir -p /workspaces/go
echo 'export GOPATH=/workspaces/go' >> /root/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> /root/.bashrc

# Install common Go tools
go install golang.org/x/tools/gopls@latest
go install github.com/go-delve/delve/cmd/dlv@latest
go install github.com/cosmtrek/air@latest

# Set permissions
chown -R root:root /workspaces/go 