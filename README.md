# distributed_database

一、安装配置
docker pull debian

docker run -it --name my_debian_0 --network mynet --ip 172.18.0.2 debian
docker start my_debian_0
docker exec -it my_debian_0 /bin/bash

docker run -it --name my_debian_1 --network mynet --ip 172.18.0.3 debian
docker start my_debian_1
docker exec -it my_debian_1 /bin/bash

docker run -it --name my_debian_2 --network mynet --ip 172.18.0.4 debian
docker start my_debian_2
docker exec -it my_debian_2 /bin/bash

docker run -it --name my_debian_3 --network mynet --ip 172.18.0.5 debian
docker start my_debian_3
docker exec -it my_debian_3 /bin/bash

docker run -it --name my_debian_4 --network mynet --ip 172.18.0.6 debian
docker start my_debian_4
docker exec -it my_debian_4 /bin/bash





apt-get update  
apt-get install git  
git


apt-get update
apt-get install vsftpd


apt-get update
apt-get install wget



wget https://mirrors.ustc.edu.cn/golang/go1.17.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version




adduser lsw
lsw


编辑/etc/vsftpd.conf，把下面这行的注释去掉。
write_enable=YES


mkdir /var/run/vsftpd
mkdir /var/run/vsftpd/empty







export GOPROXY=https://goproxy.io,direct


cd ~
git clone -b v3.4.16 https://kgithub.com/etcd-io/etcd.git
cd etcd
./build
ls
export PATH="$PATH:`pwd`/bin"


cd /home/lsw
git clone https://kgithub.com/YunzeTong/Distributed-MiniSQL

二、运行

0.	显示ip

hostname -I


1.	开启ftp

vsftpd


2.	开启etcd

export THIS_NAME=host0
export THIS_IP=172.18.0.2
cd /home/lsw/Distributed-MiniSQL
./scripts/etcd.sh

export THIS_NAME=host1
export THIS_IP=172.18.0.3
cd /home/lsw/Distributed-MiniSQL
./scripts/etcd.sh

export THIS_NAME=host2
export THIS_IP=172.18.0.4
cd /home/lsw/Distributed-MiniSQL
./scripts/etcd.sh

export THIS_NAME=host3
export THIS_IP=172.18.0.5
cd /home/lsw/Distributed-MiniSQL
./scripts/etcd.sh

export THIS_NAME=host4
export THIS_IP=172.18.0.6
cd /home/lsw/Distributed-MiniSQL
./scripts/etcd.sh


3.	开启master和region

export THIS_NAME=host0
export THIS_IP=172.18.0.2
cd /home/lsw/Distributed-MiniSQL
./scripts/master.sh

export THIS_NAME=host1
export THIS_IP=172.18.0.3
cd /home/lsw/Distributed-MiniSQL
./scripts/region.sh

export THIS_NAME=host2
export THIS_IP=172.18.0.4
cd /home/lsw/Distributed-MiniSQL
./scripts/region.sh

export THIS_NAME=host3
export THIS_IP=172.18.0.5
cd /home/lsw/Distributed-MiniSQL
./scripts/region.sh

export THIS_NAME=host4
export THIS_IP=172.18.0.6
cd /home/lsw/Distributed-MiniSQL
./scripts/region.sh


4.	开启client







go get kgithub.com/bgentry/speakeasy@v0.0.0-20210830174823-7b0c5cfc6624
注意：
--initial-advertise-peer-urls=http://172.17.0.2:2380



