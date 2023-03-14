/usr/local/bin/juicefs mount --background --cache-size 512000 redis://:Txhy2020@39.100.224.84:7000/1 /root/txhyjuicefs

/usr/bin/nohup /root/txhyjuicefs/psutil/linux/psutil -port=8082 -interval=10 >/root/psutil.log 2>&1 &