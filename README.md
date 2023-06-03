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

#安装git
apt-get update  
apt-get install git  
git

#安装vsftpd
apt-get update
apt-get install vsftpd

#安装wget
apt-get update
apt-get install wget

#安装go1.17
wget https://mirrors.ustc.edu.cn/golang/go1.17.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.17.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version

#增加用户，后面还需要输入密码等信息
adduser lsw
lsw

#进入docker文件，删除对应行
编辑/etc/vsftpd.conf，把下面这行的注释去掉。
write_enable=YES

#建立文件夹
mkdir /var/run/vsftpd
mkdir /var/run/vsftpd/empty

#引入网络代理，否则内网外网不能联通
export GOPROXY=https://goproxy.io,direct

#编译ectd,产生bin文件夹即成功
cd ~
git clone -b v3.4.16 https://github.com/etcd-io/etcd.git
cd etcd
./build
ls
export PATH="$PATH:`pwd`/bin"

#拷贝git库
cd /home/lsw
git clone https://github.com/lsw7777/DDB

二、运行

0.	显示ip
#让每个docker的终端显示ip，模拟不同主机
hostname -I


1.	开启ftp
#在每个region中开启ftp
vsftpd


2.	开启etcd
#在每个master和region中都开启etcd
#master
export THIS_NAME=host0
export THIS_IP=172.18.0.2
cd /home/lsw/DDB
./scripts/etcd.sh

#region_1
export THIS_NAME=host1
export THIS_IP=172.18.0.3
cd /home/lsw/DDB
./scripts/etcd.sh

#region_2
export THIS_NAME=host2
export THIS_IP=172.18.0.4
cd /home/lsw/DDB
./scripts/etcd.sh

#region_3
export THIS_NAME=host3
export THIS_IP=172.18.0.5
cd /home/lsw/DDB
./scripts/etcd.sh

#region_4
export THIS_NAME=host4
export THIS_IP=172.18.0.6
cd /home/lsw/DDB
./scripts/etcd.sh


3.	开启master和region

#master
export THIS_NAME=host0
export THIS_IP=172.18.0.2
cd /home/lsw/DDB
./scripts/master.sh

#region_1
export THIS_NAME=host1
export THIS_IP=172.18.0.3
cd /home/lsw/DDB
./scripts/region.sh

#region_2
export THIS_NAME=host2
export THIS_IP=172.18.0.4
cd /home/lsw/DDB
./scripts/region.sh

#region_3
export THIS_NAME=host3
export THIS_IP=172.18.0.5
cd /home/lsw/DDB
./scripts/region.sh

#region_4
export THIS_NAME=host4
export THIS_IP=172.18.0.6
cd /home/lsw/DDB
./scripts/region.sh


4.	开启client
#在一个region上开启
cd /home/lsw/DDB
./scripts/client.sh










