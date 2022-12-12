# HLFNetwork-V2
The V2 implementation of the fabric network which involves API calls and separate CouchDB instances for each peer.

## AWS EC2 Setup
1. Create a new EC2 Ubuntu instance with more than 50gb of storage.
2. SLAC employees can use the launch-wizard-18 security group, otherwise ensure that port 4000 (App) and 5984 (CouchDB) are open
3. SSH into the instance and clone this repository using `git clone https://github.com/Secure-Identity/HLFNetwork-V2`

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

## Start the application
1. `cd api-1.4`
2. `nodemon app.js`

## API requests
1. User Registration: `http://<ec2-instance-ip>:4000/users` Method: POST, Body: {
    "username":<enter-username>,
    "orgName":"Org1"
}
Copy the token from the response and use that as the bearer token in the following requests.
- IP address of a docker container can be found by using - `docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' <container id>`
- We'll use IPs of dev-peer.... docker containers only.
1. ProvisionID: `http://<ec2-instance-ip>:4000/channels/mychannel/chaincodes/secureID?args=["DEVICE"]&peer=peer0.org1.secretidentity.com&fcn=provisionID`
2. Sharding: `http://<ec2-instance-ip>:4000/channels/mychannel/chaincodes/secureID?args=["DEVICE0"]&peer=peer0.org1.secretidentity.com&fcn=keymaker`
3. Listen Shards: `http://<ec2-instance-ip>:4000/channels/mychannel/chaincodes/secureidfinal?args=["<listener-docker-container-IP>"]&peer=peer0.org1.secretidentity.com&fcn=listenShard`
4. Send Shards: `http://<ec2-instance-ip>:4000/channels/mychannel/chaincodes/secureidfinal?args=["<sender-docker-container-IP>"]&peer=peer0.org1.secretidentity.com&fcn=sendShard`
 
## References
 1. https://hyperledger-fabric.readthedocs.io/en/latest/test_network.html
 2. https://www.youtube.com/watch?v=SJTdJt6N6Ow&list=PLSBNVhWU6KjW4qo1RlmR7cvvV8XIILub6&index=1
 3. https://github.com/adhavpavan/BasicNetwork-2.0
