#!/bin/bash
# 设置fabric-ca-client home环境变量
source /etc/profile
# Set variables
ORG_NAME=$1
ORDERER_NAME=$2
ORDERER_DOMAIN=$3
PORT=$4
CA_CLIENT_HOME=/root/txhyjuicefs/organizations/fabric-ca/$ORG_NAME
ORDERERS_DIR=/root/txhyjuicefs/organizations/$ORG_NAME/orderers
ORDERER_DIR=/root/txhyjuicefs/organizations/$ORG_NAME/orderers/$ORDERER_DOMAIN
ORG_DIR=/root/txhyjuicefs/organizations/$ORG_NAME

# Set environment variables
export FABRIC_CA_CLIENT_HOME=$CA_CLIENT_HOME

# Define functions
create_dir() {
  mkdir -p $1
}

copy_file() {
  cp $1 $2
}

# Create organization directory
create_dir $ORDERERS_DIR
create_dir $ORDERER_DIR

fabric-ca-client enroll -d -u https://admin:$ORG_NAME@localhost:$PORT --caname ca-$ORG_NAME --tls.certfiles $CA_CLIENT_HOME/ca-cert.pem
# Register orderer with CA
fabric-ca-client register -d --caname ca-$ORG_NAME --id.name $ORDERER_NAME --id.secret $ORDERER_NAME --id.type orderer --tls.certfiles $CA_CLIENT_HOME/ca-cert.pem

# Enroll orderer with CA
fabric-ca-client enroll -u https://$ORDERER_NAME:$ORDERER_NAME@localhost:$PORT --caname ca-$ORG_NAME -M $ORDERER_DIR/msp --csr.hosts $ORDERER_DOMAIN --tls.certfiles $CA_CLIENT_HOME/ca-cert.pem

fabric-ca-client enroll -u https://$ORDERER_NAME:$ORDERER_NAME@localhost:$PORT --caname ca-$ORG_NAME -M $ORDERER_DIR/tls --enrollment.profile tls --csr.hosts $ORDERER_DOMAIN --csr.hosts localhost --tls.certfiles $CA_CLIENT_HOME/ca-cert.pem

# Copy TLS and MSP files
copy_file $ORG_DIR/msp/config.yaml $ORDERER_DIR/msp/config.yaml
copy_file $ORDERER_DIR/tls/tlscacerts/* $ORDERER_DIR/tls/ca.crt
copy_file $ORDERER_DIR/tls/signcerts/* $ORDERER_DIR/tls/server.crt
copy_file $ORDERER_DIR/tls/keystore/* $ORDERER_DIR/tls/server.key
create_dir $ORDERER_DIR/msp/tlscacerts
copy_file $ORDERER_DIR/tls/tlscacerts/* $ORDERER_DIR/msp/tlscacerts/tlsca.example.com-cert.pem
