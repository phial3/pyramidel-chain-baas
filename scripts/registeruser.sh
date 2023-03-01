#!/usr/bin/env bash
# $1 uscc $2 用户名 $3 密码 $4 类型 $5 ca端口
export FABRIC_CA_CLIENT_HOME=/root/txhyjuicefs/organizations/fabric-ca/$1/

infoln "Register user"
set -x
fabric-ca-client register -d --caname ca-$1 --id.name $2 --id.secret $3 --id.type $4 --tls.certfiles /root/txhyjuicefs/organizations/fabric-ca/$1/ca-cert.pem
{ set +x; } 2>/dev/null

mkdir -p /root/txhyjuicefs/organizations/$1/users
mkdir -p /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com

infoln "Generate the user msp"
set -x
fabric-ca-client enroll -u https://$2:$3@localhost:$5 --caname ca-$1 -M /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com/msp --tls.certfiles /root/txhyjuicefs/organizations/fabric-ca/$1/ca-cert.pem
{ set +x; } 2>/dev/null

cp /root/txhyjuicefs/organizations/$1/msp/config.yaml /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com/msp/config.yaml

infoln "Generate the user tls"
set -x
fabric-ca-client enroll -u https://$2:$3@localhost:$5 --caname ca-$1 -M /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com/tls --enrollment.profile tls --tls.certfiles /root/txhyjuicefs/organizations/fabric-ca/$1/ca-cert.pem
{ set +x; } 2>/dev/null

cp /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com/tls/tlscacerts/* /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com/tls/ca.crt
cp /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com/tls/signcerts/* /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com/tls/client.crt
cp /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com/tls/keystore/* /root/txhyjuicefs/organizations/$1/users/$2@$1.pcb.com/tls/client.key
