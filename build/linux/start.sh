#!/bin/bash
DIR="/data/k8s/"
mkdir $DIR
#if [ -d "$DIR" ]; then
#   sh -c /baasapi
#   exit 0
#fi
tar -zxvf /tmp/k8s.tar.gz -C /data/k8s
sh -c /baasapi
