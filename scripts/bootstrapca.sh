#!/bin/bash
if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <org_name> <port>" >&2
  exit 1
fi
# 设置fabric-ca-client home环境变量
source /etc/profile

# 设置变量
org_name=$1
port=$2
ca_client_home=/root/txhyjuicefs/organizations/fabric-ca/$org_name
org_dir=/root/txhyjuicefs/organizations/$org_name

export FABRIC_CA_CLIENT_HOME=ca_client_home

# 创建目录函数
create_dir() {
  mkdir -p $1
}

# 复制文件函数
copy_file() {
  cp $1 $2
}

# 创建组织目录
create_org_dir() {
  create_dir $org_dir
  create_dir $org_dir/msp
  create_dir $org_dir/msp/tlscacerts
  create_dir $org_dir/tlsca
  create_dir $org_dir/ca

  # 写入 config.yaml 文件
  cat >$org_dir/msp/config.yaml <<EOF
NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-$port-ca-$org_name.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-$port-ca-$org_name.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-$port-ca-$org_name.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-$port-ca-$org_name.pem
    OrganizationalUnitIdentifier: orderer
EOF

  # 使用 fabric-ca-client 命令进行注册和 enroll
  fabric-ca-client enroll -d -u https://admin:$org_name@localhost:$port \
    --caname ca-$org_name --tls.certfiles $ca_client_home/ca-cert.pem

  # 复制证书文件到指定目录
  copy_file $ca_client_home/ca-cert.pem $org_dir/msp/tlscacerts/ca.crt
  copy_file $ca_client_home/ca-cert.pem $org_dir/msp/tlscacerts/tlsca.example.com-cert.pem
  copy_file $ca_client_home/ca-cert.pem $org_dir/tlsca/tlsca.$org_name.pcb.com-cert.pem
  copy_file $ca_client_home/ca-cert.pem $org_dir/ca/ca.$org_name.pcb.com-cert.pem
}

# 调用函数
create_org_dir

# 输出提示信息
echo "Organization directory created successfully."
