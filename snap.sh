#!/bin/sh

cd /go/bin
mkdir -p /var/log/snap
snapteld --plugin-trust 0 --log-level 1 --log-path /var/log/snap &
snaptel plugin load snap-plugin-publisher-file
snaptel plugin load snap-plugin-collector-psutil
snaptel plugin load snap-plugin-publisher-elasticsearch
snaptel plugin list
echo snaptel task create -t /tmp/task.yml
