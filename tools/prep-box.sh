#!/bin/bash
CONTAINER_URL=https://download.docker.com/linux/debian/dists/buster/pool/stable/amd64/containerd.io_1.4.3-1_amd64.deb
DOCKER_CLI_URL=https://download.docker.com/linux/debian/dists/buster/pool/stable/amd64/docker-ce-cli_20.10.2~3-0~debian-buster_amd64.deb
DOCKER_CE_URL=https://download.docker.com/linux/debian/dists/buster/pool/stable/amd64/docker-ce_20.10.2~3-0~debian-buster_amd64.deb

CONTAINER_PATH=containerd.io_1.4.3-1_amd64.deb
DOCKER_CLI_PATH=docker-ce-cli_20.10.2~3-0~debian-buster_amd64.deb
DOCKER_CE_PATH=docker-ce_20.10.2~3-0~debian-buster_amd64.deb
KEY_LOC=/home/holden/.ssh/us-west-2-lightsail-default.pem
REMOTE_KEY_LOC=/opt/gametime/holdongametime/conf/remote-aws

get_docker_debs () 
{
	mkdir docker
	pushd docker
	curl -L -O $CONTAINER_URL
	curl -L -O $DOCKER_CLI_URL
	curl -L -O $DOCKER_CE_URL
	apt install $CONTAINER_PATH
        apt install $DOCKER_CLI_PATH
	apt install $DOCKER_CE_PATH
	popd
}

get_git_stuff ()
{
	eval `ssh-agent -s`
	ssh-add ~/remote-aws
	git clone git@github.com:habruzzo/gametime.git

}	
prep ()
{
	read -ep "What user do you want to access?" user
	read -ep "What IP do you want to access?" ip_addr
	echo "$0"
	echo "$user@$ip_addr"
	echo "sudo $0 pickup"
	sudo chmod u+x $0
	scp -i $KEY_LOC $0 $user@$ip_addr:~/$0
	scp -r -i $KEY_LOC $REMOTE_KEY_LOC $user@$ip_addr:~/remote-aws
	ssh -i $KEY_LOC $user@$ip_addr "sudo chmod u+x $0;sudo $0 pickup"
}

case $1 in
	start)
		echo "Starting prep"
		prep
	;;
	pickup)
		echo "Starting pickup"
		get_docker_debs
		get_git_stuff
	;;
	*)
		exit
	;;
esac
