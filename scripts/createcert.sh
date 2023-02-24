#!/bin/bash

if [ $# -eq 2 ]; then
  echo $1,$2
  cat >ca-config.json <<EOF
{
  "signing": {
    "default": {
      "expiry": "87600h"
    },
    "profiles": {
      "docker": {
        "usages": [
            "signing",
            "key encipherment",
            "server auth",
            "client auth"
        ],
        "expiry": "87600h"
      }
    }
  }
}
EOF

  cat >ca-csr.json <<EOF
{
  "CN": "LiKe Personal CA",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "SHANXI",
      "L": "TAIYUAN",
      "O": "TXHY",
      "OU": "DEV"
    }
  ],
  "ca": {
    "expiry": "87600h"
 }
}
EOF
  cfssl gencert -initca ca-csr.json | cfssljson -bare ca

  cat >server-csr.json <<EOF
{
  "CN": "Docker Server",
  "hosts": [
    "127.0.0.1",
    "0.0.0.0",
    "localhost",
    "$1",
    "$2"
  ],
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "SHANXI",
      "L": "TAIYUAN",
      "O": "TXHY",
      "OU": "DEV"
    }
  ]
}
EOF

  cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ca-config.json -profile docker server-csr.json | cfssljson -bare server

  cat >client-csr.json <<EOF
{
  "CN": "Docker Client",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "SHANXI",
      "L": "TAIYUAN",
      "O": "TXHY",
      "OU": "DEV"
    }
  ]
}
EOF

  cfssl gencert -ca ca.pem -ca-key ca-key.pem -config ca-config.json -profile docker client-csr.json | cfssljson -bare client
  chmod -v 0400 ca-key.pem server-key.pem client-key.pem
  chmod -v 0444 ca.pem server.pem client.pem

  \cp -rf * /etc/docker
  \cp -rf * ~/.docker

  sed -i "13c ExecStart=/usr/bin/dockerd --tlsverify --tlscacert=/etc/docker/ca.pem --tlscert=/etc/docker/server.pem --tlskey=/etc/docker/server-key.pem -H fd:// --containerd=/run/containerd/containerd.sock -H tcp://0.0.0.0:2376 -H unix:///var/run/docker.sock" /usr/lib/systemd/system/docker.service
  systemctl daemon-reload
  systemctl enable docker
  systemctl restart docker.service
else
  echo 需要参数: 2, 实际参数: $#.
fi
