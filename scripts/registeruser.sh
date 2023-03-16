#!/bin/bash
# $1 uscc $2 用户名 $3 密码 $4 类型 $5 ca端口
# Check if required argument(s) is missing
if [ "$#" -ne 5 ]; then
  echo "Usage: $0 <uscc> <username> <password> <type> <ca-port>"
  exit 1
fi
source /etc/profile
export FABRIC_CA_CLIENT_HOME=/root/txhyjuicefs/organizations/fabric-ca/$1


# Define variables
uscc=$1
username=$2
password=$3
type=$4
ca_port=$5
ca_name=ca-$uscc
ca_client_home=/root/txhyjuicefs/organizations/fabric-ca/$uscc
org_dir=/root/txhyjuicefs/organizations/$uscc
user_home=/root/txhyjuicefs/organizations/$uscc/users/$username@$uscc.pcb.com

fabric-ca-client enroll -d -u https://admin:$uscc@localhost:$ca_port --caname ca-$uscc --tls.certfiles $ca_client_home/ca-cert.pem
# Register user with CA
fabric-ca-client register -d --caname $ca_name --id.name $username --id.secret $password --id.type $type --tls.certfiles $ca_client_home/ca-cert.pem

# Enroll user with CA
fabric-ca-client enroll -u https://$username:$password@localhost:$ca_port --caname $ca_name -M $user_home/msp --tls.certfiles $ca_client_home/ca-cert.pem
fabric-ca-client enroll -u https://$username:$password@localhost:$ca_port --caname $ca_name -M $user_home/tls --enrollment.profile tls --tls.certfiles $ca_client_home/ca-cert.pem

# Copy user's certs to appropriate locations
mkdir -p  $user_home/msp/tlscacerts
cp $ca_client_home/ca-cert.pem $user_home/msp/tlscacerts/ca.crt
cp $user_home/tls/tlscacerts/* $user_home/tls/ca.crt
cp $user_home/tls/signcerts/* $user_home/tls/client.crt
cp $user_home/tls/keystore/* $user_home/tls/client.key

# Copy config.yaml to user's msp directory
cp $org_dir/msp/config.yaml $user_home/msp/config.yaml
