createcertificatesForOrg1() {
  echo
  echo "Enroll the CA admin"
  echo
  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/
  export FABRIC_CA_CLIENT_HOME=${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/


  fabric-ca-client enroll -u https://admin:adminpw@localhost:7054 --caname ca.org1.secretidentity.com --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem


  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1-secretidentity-com.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1-secretidentity-com.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1-secretidentity-com.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-7054-ca-org1-secretidentity-com.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/msp/config.yaml

  echo
  echo "Register peer0"
  echo
  fabric-ca-client register --caname ca.org1.secretidentity.com --id.name peer0 --id.secret peer0pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  echo
  echo "Register peer1"
  echo
  fabric-ca-client register --caname ca.org1.secretidentity.com --id.name peer1 --id.secret peer1pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  echo
  echo "Register peer2"
  echo
  fabric-ca-client register --caname ca.org1.secretidentity.com --id.name peer2 --id.secret peer2pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  echo
  echo "Register peer3"
  echo
  fabric-ca-client register --caname ca.org1.secretidentity.com --id.name peer3 --id.secret peer3pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  echo
  echo "Register peer4"
  echo
  fabric-ca-client register --caname ca.org1.secretidentity.com --id.name peer4 --id.secret peer4pw --id.type peer --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  echo
  echo "Register user"
  echo
  fabric-ca-client register --caname ca.org1.secretidentity.com --id.name user1 --id.secret user1pw --id.type client --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  echo
  echo "Register the org admin"
  echo
  fabric-ca-client register --caname ca.org1.secretidentity.com --id.name org1admin --id.secret org1adminpw --id.type admin --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers

  # -----------------------------------------------------------------------------------
  #  Peer 0
  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com

  echo
  echo "## Generate the peer0 msp"
  echo
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/msp --csr.hosts peer0.org1.secretidentity.com --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/msp/config.yaml

  echo
  echo "## Generate the peer0-tls certificates"
  echo
  fabric-ca-client enroll -u https://peer0:peer0pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/tls --enrollment.profile tls --csr.hosts peer0.org1.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/tls/signcerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/tls/keystore/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/tls/server.key

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/msp/tlscacerts/ca.crt

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/tlsca
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/tlsca/tlsca.org1.secretidentity.com-cert.pem

  mkdir ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/ca
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer0.org1.secretidentity.com/msp/cacerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/ca/ca.org1.secretidentity.com-cert.pem

  # ------------------------------------------------------------------------------------------------

  # Peer1

  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com

  echo
  echo "## Generate the peer1 msp"
  echo
  fabric-ca-client enroll -u https://peer1:peer1pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com/msp --csr.hosts peer1.org1.secretidentity.com --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com/msp/config.yaml

  echo
  echo "## Generate the peer1-tls certificates"
  echo
  fabric-ca-client enroll -u https://peer1:peer1pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com/tls --enrollment.profile tls --csr.hosts peer1.org1.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com/tls/signcerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com/tls/keystore/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer1.org1.secretidentity.com/tls/server.key

  # --------------------------------------------------------------------------------------------------

  # Peer2

  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com

  echo
  echo "## Generate the peer2 msp"
  echo
  fabric-ca-client enroll -u https://peer2:peer2pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com/msp --csr.hosts peer2.org1.secretidentity.com --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com/msp/config.yaml

  echo
  echo "## Generate the peer2-tls certificates"
  echo
  fabric-ca-client enroll -u https://peer2:peer2pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com/tls --enrollment.profile tls --csr.hosts peer2.org1.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com/tls/signcerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com/tls/keystore/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com/tls/server.key

  # --------------------------------------------------------------------------------------------------

  # Peer3

  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com

  echo
  echo "## Generate the peer3 msp"
  echo
  fabric-ca-client enroll -u https://peer3:peer3pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com/msp --csr.hosts peer3.org1.secretidentity.com --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com/msp/config.yaml

  echo
  echo "## Generate the peer3-tls certificates"
  echo
  fabric-ca-client enroll -u https://peer3:peer3pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com/tls --enrollment.profile tls --csr.hosts peer3.org1.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com/tls/signcerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com/tls/keystore/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer3.org1.secretidentity.com/tls/server.key

  # --------------------------------------------------------------------------------------------------

  # Peer4

  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer2.org1.secretidentity.com

  echo
  echo "## Generate the peer4 msp"
  echo
  fabric-ca-client enroll -u https://peer4:peer4pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer4.org1.secretidentity.com/msp --csr.hosts peer4.org1.secretidentity.com --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer4.org1.secretidentity.com/msp/config.yaml

  echo
  echo "## Generate the peer4-tls certificates"
  echo
  fabric-ca-client enroll -u https://peer4:peer4pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer4.org1.secretidentity.com/tls --enrollment.profile tls --csr.hosts peer4.org1.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer4.org1.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer4.org1.secretidentity.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer4.org1.secretidentity.com/tls/signcerts/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer4.org1.secretidentity.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer4.org1.secretidentity.com/tls/keystore/* ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/peers/peer4.org1.secretidentity.com/tls/server.key

  # --------------------------------------------------------------------------------------------------

  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/users
  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/users/User1@org1.secretidentity.com

  echo
  echo "## Generate the user msp"
  echo
  fabric-ca-client enroll -u https://user1:user1pw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/users/User1@org1.secretidentity.com/msp --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  mkdir -p crypto-config-ca/peerOrganizations/org1.secretidentity.com/users/Admin@org1.secretidentity.com

  echo
  echo "## Generate the org admin msp"
  echo
  fabric-ca-client enroll -u https://org1admin:org1adminpw@localhost:7054 --caname ca.org1.secretidentity.com -M ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/users/Admin@org1.secretidentity.com/msp --tls.certfiles ${PWD}/fabric-ca/org1/tls-cert.pem

  cp ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/peerOrganizations/org1.secretidentity.com/users/Admin@org1.secretidentity.com/msp/config.yaml

}

# createcertificatesForOrg1

createCertificateForOrderer() {
  echo
  echo "Enroll the CA admin"
  echo
  mkdir -p crypto-config-ca/ordererOrganizations/secretidentity.com

  export FABRIC_CA_CLIENT_HOME=${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com


  fabric-ca-client enroll -u https://admin:adminpw@localhost:9054 --caname ca-orderer --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem


  echo 'NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-9054-ca-orderer.pem
    OrganizationalUnitIdentifier: orderer' >${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/config.yaml

  echo
  echo "Register orderer"
  echo

  fabric-ca-client register --caname ca-orderer --id.name orderer --id.secret ordererpw --id.type orderer --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem

  #  orderer admin
  # fabric-ca-client register --caname ca-orderer --id.name ordererAdmin --id.secret ordererAdminpw --id.type admin --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem


  mkdir -p crypto-config-ca/ordererOrganizations/secretidentity.com/orderers
  # mkdir -p crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/secretidentity.com

  # ---------------------------------------------------------------------------
  #  Orderer

  mkdir -p crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com

  echo
  echo "## Generate the orderer msp"
  echo

  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/msp --csr.hosts orderer.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem


  cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/msp/config.yaml

  echo
  echo "## Generate the orderer-tls certificates"
  echo

  fabric-ca-client enroll -u https://orderer:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/tls --enrollment.profile tls --csr.hosts orderer.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem


  cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/tls/ca.crt
  cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/tls/signcerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/tls/server.crt
  cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/tls/keystore/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/tls/server.key

  mkdir ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/msp/tlscacerts/tlsca.secretidentity.com-cert.pem

  mkdir ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/tlscacerts
  cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/tlscacerts/tlsca.secretidentity.com-cert.pem

  # -----------------------------------------------------------------------
  #  Orderer 2

  # mkdir -p crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com

  # echo
  # echo "## Generate the orderer msp"
  # echo

  # fabric-ca-client enroll -u https://orderer2:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/msp --csr.hosts orderer2.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem


  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/msp/config.yaml

  # echo
  # echo "## Generate the orderer-tls certificates"
  # echo

  # fabric-ca-client enroll -u https://orderer2:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/tls --enrollment.profile tls --csr.hosts orderer2.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem


  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/tls/ca.crt
  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/tls/signcerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/tls/server.crt
  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/tls/keystore/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/tls/server.key

  # mkdir ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/msp/tlscacerts
  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/msp/tlscacerts/tlsca.secretidentity.com-cert.pem

  # # mkdir ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/tlscacerts
  # # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer2.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/tlscacerts/tlsca.secretidentity.com-cert.pem

  # # ---------------------------------------------------------------------------
  # #  Orderer 3
  # mkdir -p crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com

  # echo
  # echo "## Generate the orderer msp"
  # echo

  # fabric-ca-client enroll -u https://orderer3:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/msp --csr.hosts orderer3.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem


  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/msp/config.yaml

  # echo
  # echo "## Generate the orderer-tls certificates"
  # echo

  # fabric-ca-client enroll -u https://orderer3:ordererpw@localhost:9054 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/tls --enrollment.profile tls --csr.hosts orderer3.secretidentity.com --csr.hosts localhost --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem


  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/tls/ca.crt
  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/tls/signcerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/tls/server.crt
  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/tls/keystore/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/tls/server.key

  # mkdir ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/msp/tlscacerts
  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/msp/tlscacerts/tlsca.secretidentity.com-cert.pem

  # mkdir ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/tlscacerts
  # cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/orderers/orderer3.secretidentity.com/tls/tlscacerts/* ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/tlscacerts/tlsca.secretidentity.com-cert.pem

  # ---------------------------------------------------------------------------

  mkdir -p crypto-config-ca/ordererOrganizations/secretidentity.com/users
  mkdir -p crypto-config-ca/ordererOrganizations/secretidentity.com/users/Admin@secretidentity.com

  echo
  echo "## Generate the admin msp"
  echo

  fabric-ca-client enroll -u https://ordererAdmin:ordererAdminpw@localhost:9054 --caname ca-orderer -M ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/users/Admin@secretidentity.com/msp --tls.certfiles ${PWD}/fabric-ca/ordererOrg/tls-cert.pem


  cp ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/msp/config.yaml ${PWD}/crypto-config-ca/ordererOrganizations/secretidentity.com/users/Admin@secretidentity.com/msp/config.yaml

}

# createCertificateForOrderer

sudo rm -rf crypto-config-ca/*
# sudo rm -rf fabric-ca/*
createcertificatesForOrg1
createCertificateForOrderer

