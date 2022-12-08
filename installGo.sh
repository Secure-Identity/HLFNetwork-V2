#!/bin/sh
cd $HOME
wget https://go.dev/dl/go1.19.3.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.19.3.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
echo "---------------------------"
echo "go installed successfully!"
go version
echo "---------------------------"
