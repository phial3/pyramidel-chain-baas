#!/bin/bash
#参数1组织uscc, 参数2 peer名,名字和密码使用同一参数, 参数3 peer域名, 参数4 ca端口
# 设置fabric-ca-client home环境变量
source /etc/profile

ORG_NAME=$1
PEER_NAME=$2
PEER_DOMAIN=$3
PORT=$4
CA_CLIENT_HOME=/txhyjuicefs/organizations/fabric-ca/$ORG_NAME
PEERS_DIR=/txhyjuicefs/organizations/$ORG_NAME/peers
PEER_DIR=$PEERS_DIR/$PEER_DOMAIN
ORG_DIR=/txhyjuicefs/organizations/$ORG_NAME

# Set environment variables
export FABRIC_CA_CLIENT_HOME=$CA_CLIENT_HOME

# 创建目录函数
create_dir() {
  mkdir -p $1
}

# 复制文件函数
copy_file() {
  cp $1 $2
}

# 创建组织目录
create_dir $PEERS_DIR
create_dir $PEER_DIR

fabric-ca-client enroll -d -u https://admin:$ORG_NAME@localhost:$PORT --caname ca-$ORG_NAME --tls.certfiles $CA_CLIENT_HOME/ca-cert.pem
# 注册peer
fabric-ca-client register -d --caname ca-$ORG_NAME --id.name $PEER_NAME --id.secret $PEER_NAME --id.type peer --tls.certfiles $CA_CLIENT_HOME/ca-cert.pem

# Enroll peer身份
fabric-ca-client enroll -u https://$PEER_NAME:$PEER_NAME@localhost:$PORT --caname ca-$ORG_NAME -M $PEER_DIR/msp --csr.hosts $PEER_DOMAIN --tls.certfiles $CA_CLIENT_HOME/ca-cert.pem

# Enroll peer TLS身份
fabric-ca-client enroll -u https://$PEER_NAME:$PEER_NAME@localhost:$PORT --caname ca-$ORG_NAME -M $PEER_DIR/tls --enrollment.profile tls --csr.hosts $PEER_DOMAIN --csr.hosts localhost --tls.certfiles $CA_CLIENT_HOME/ca-cert.pem

# 复制配置文件和证书
copy_file $ORG_DIR/msp/config.yaml $PEER_DIR/msp/config.yaml
copy_file $PEER_DIR/tls/tlscacerts/* $PEER_DIR/tls/ca.crt
copy_file $PEER_DIR/tls/signcerts/* $PEER_DIR/tls/server.crt
copy_file $PEER_DIR/tls/keystore/* $PEER_DIR/tls/server.key
