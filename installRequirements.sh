#!/bin/sh

cd $HOME
sudo apt-get update
sudo apt-get install curl
sudo apt-get install nodejs
sudo apt-get install npm
sudo apt-get install python-is-python3

sudo npm install -g nodemon

sudo apt-get install \
    ca-certificates \
    curl \
    gnupg \
    lsb-release

sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  bionic stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get install docker.io
apt-get update
apt-cache policy docker-ce
sudo apt-get install -y docker-ce
sudo apt-get install docker-compose
sudo apt-get upgrade

curl -sL https://deb.nodesource.com/setup_14.x | sudo bash -
sudo apt-get install -y nodejs

sudo service docker start
sudo groupadd docker
sudo usermod -aG docker $USER
sudo chmod 666 /var/run/docker.sock

wget https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh
chmod +x install-fabric.sh
./install-fabric.sh docker
mkdir $HOME/secureID && cd $HOME/secureID
curl -sSL https://bit.ly/2ysbOFE | bash -s - 2.4.7 1.5 0.4.18
