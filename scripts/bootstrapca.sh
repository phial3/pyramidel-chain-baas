mkdir -p /root/txhyjuicefs/organizations/$1/

export FABRIC_CA_CLIENT_HOME=/root/txhyjuicefs/organizations/fabric-ca/$1/

fabric-ca-client enroll -d -u https://admin:$1@localhost:$2 --caname ca-$1 --tls.certfiles /root/txhyjuicefs/organizations/fabric-ca/$1/ca-cert.pem

mkdir -p /root/txhyjuicefs/organizations/$1/msp
touch /root/txhyjuicefs/organizations/$1/msp/config.yaml
cat >/root/txhyjuicefs/organizations/$1/msp/config.yaml <<EOF
NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/localhost-$2-ca-$1.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/localhost-$2-ca-$1.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/localhost-$2-ca-$1.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/localhost-$2-ca-$1.pem
    OrganizationalUnitIdentifier: orderer
EOF

mkdir -p /root/txhyjuicefs/organizations/$1/msp/tlscacerts
cp /root/txhyjuicefs/organizations/fabric-ca/$1/ca-cert.pem /root/txhyjuicefs/organizations/$1/msp/tlscacerts/ca.crt
cp /root/txhyjuicefs/organizations/fabric-ca/$1/ca-cert.pem /root/txhyjuicefs/organizations/$1/msp/tlscacerts/tlsca.example.com-cert.pem

# Copy org's CA cert to org's /tlsca directory (for use by clients)
mkdir -p /root/txhyjuicefs/organizations/$1/tlsca
cp /root/txhyjuicefs/organizations/fabric-ca/$1/ca-cert.pem /root/txhyjuicefs/organizations/$1/tlsca/tlsca.$1.pcb.com-cert.pem

# Copy org's CA cert to org's /ca directory (for use by clients)
mkdir -p /root/txhyjuicefs/organizations/$1/ca
cp /root/txhyjuicefs/organizations/fabric-ca/$1/ca-cert.pem /root/txhyjuicefs/organizations/$1/ca/ca.$1.pcb.com-cert.pem
