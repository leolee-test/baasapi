#!/bin/bash
rm ./dist/baasapi
cp -rf /media/sf_Downloads/cloud9/api/* ./
cp -rf /media/sf_Downloads/cloud9/api/cmd ./api/
cp -rf /media/sf_Downloads/cloud9/api/exec ./api/
cp -rf /media/sf_Downloads/cloud9/api/http ./api/
cp -rf /media/sf_Downloads/cloud9/api/cron ./api/
cp -rf /media/sf_Downloads/cloud9/api/cli ./api/
cp -rf /media/sf_Downloads/cloud9/api/filesystem ./api/
cp -rf /media/sf_Downloads/cloud9/api/baasapi.go ./api/
cp -rf /media/sf_Downloads/cloud9/api/errors.go ./api/
cp -rf /media/sf_Downloads/cloud9/api ./
cp -rf /media/sf_Downloads/cloud9/api/filesystem ./api/      
cp -rf /media/sf_Downloads/cloud9/api/bolt ./api/      
cp -rp /media/sf_Downloads/baasui/k8s_ansible/k8s/* ./build/linux/k8s/  
tar -zcvf ./build/linux/k8s.tar.gz ./build/linux/k8s
