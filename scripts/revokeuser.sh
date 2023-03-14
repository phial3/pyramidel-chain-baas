#!/bin/bash
# $1 uscc $2 用户名 $3 密码 $4 类型 $5 ca端口
# Check if required argument(s) is missing
if [ "$#" -ne 3 ]; then
  echo "Usage: $0 <uscc> <username> <ca-port>"
  exit 1
fi

source /etc/profile
export FABRIC_CA_CLIENT_HOME=/txhyjuicefs/organizations/fabric-ca/$1

# Define variables
uscc=$1
username=$2
port=$3

fabric-ca-client enroll -d -u https://admin:$uscc@localhost:$port --caname ca-$uscc --tls.certfiles $FABRIC_CA_CLIENT_HOME/ca-cert.pem

fabric-ca-client revoke -e $username -r unspecified --gencrl --tls.certfiles $FABRIC_CA_CLIENT_HOME/ca-cert.pem