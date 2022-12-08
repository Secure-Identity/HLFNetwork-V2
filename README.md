# HLFNetwork-V2
The second implementation of the fabric network which involves API calls and separate couchdb instances for each peer.

After cloning the repository and switching to the HLFNetwork-V2 directory, follow the steps:

## Install Requirements
1. `chmod +x installRequirements.sh && chmod +x installGo.sh`
2. `sudo ./installRequirements.sh`
3. `sudo ./installGo.sh`

## Updating your shell environment
1. Open your bashrc file in an editor: `vi ~/.bashrc`
2. Add the following lines: <br>
 `export GOROOT=/usr/local/go` <br>
`export PATH=$PATH:$GOROOT/bin`<br>
`export PATH=$PATH:$HOME/secureID/fabric-samples/bin`
3. Run `source ~/.bashrc` to apply the new changes

## Creating crypto artifacts and starting the Hyperledger Fabric network
1. `cd $HOME/HLFNetwork-V2/artifacts/channel && ./create-artifacts.sh`
2. `cd $HOME/HLFNetwork-V2/artifacts/channel && docker-compose up -d`

### Wait for atleast 1-2 minutes before executing the following steps

## Creating a channel
`cd $HOME/HLFNetwork-V2 && ./createChannel.sh`

## Deploying the SecureID chaincode on the channel
`cd $HOME/HLFNetwork-V2 && ./deploySecureID.sh`
