#!/bin/bash
CONTAINER_URL=https://download.docker.com/linux/debian/dists/buster/pool/stable/amd64/containerd.io_1.4.3-1_amd64.deb
DOCKER_CLI_URL=https://download.docker.com/linux/debian/dists/buster/pool/stable/amd64/docker-ce-cli_20.10.2~3-0~debian-buster_amd64.deb
DOCKER_CE_URL=https://download.docker.com/linux/debian/dists/buster/pool/stable/amd64/docker-ce_20.10.2~3-0~debian-buster_amd64.deb

CONTAINER_PATH=containerd.io_1.4.3-1_amd64.deb
DOCKER_CLI_PATH=docker-ce-cli_20.10.2~3-0~debian-buster_amd64.deb
DOCKER_CE_PATH=docker-ce_20.10.2~3-0~debian-buster_amd64.deb
KEY_LOC=/home/holden/.ssh/us-west-2-lightsail-default.pem
REMOTE_KEY_LOC=/opt/gametime/conf/remote-aws
SCRIPT_NAME=prep-box.sh


reboot_box ()
{
	user=centos
	#read -ep "What user do you want to access?" user
	ip_addr=$1
	ssh -i $KEY_LOC $user@$ip_addr "sudo reboot"
}

reload ()
{
	user=centos
	#read -ep "What user do you want to access?" user
	ip_addr=$1
	ssh -i $KEY_LOC $user@$ip_addr "./$SCRIPT_NAME copy"
}

copy_conf_files ()
{
	cd /etc/httpd/conf
	sudo chmod 755 *
	sudo cp /home/centos/conf/httpd.conf /etc/httpd/conf/httpd.conf

	sudo service httpd restart
}

install_git_deps ()
{
	pip3 install --user -r gametime/requirements.txt
	
	npm install lessc
}

get_git_stuff ()
{
	cd /home/centos
	eval `ssh-agent -s`
	ssh-add remote-aws/id_ed25519
	ssh-keyscan -t rsa -H github.com >> /home/centos/.ssh/known_hosts	
	sleep 5
	
	git clone git@github.com:habruzzo/gametime.git
	
	sudo mv gametime/holdongametime /opt
	sudo ln -s /opt/holdongametime gametime/holdongametime
	sudo chmod -R 775 /opt
	sudo chgrp -R apache /opt

	sudo mkdir /srv/http
	sudo ln -s /home/centos/gametime/holdongametime/static /srv/http/static
	sudo ln -s /home/centos/gametime/holdongametime/templates /srv/http/templates
}

setup_deps ()
{
	sudo yum -y install docker
	sudo yum -y install httpd
	sudo yum -y install git
	sudo yum -y install python3
	sudo yum -y install epel-release
	sudo yum -y install python-pip
	sudo yum -y install npm
	sudo yum -y install gcc
	sudo yum -y install mod_wsgi
    sudo yum -y install python3-devel
    sudo sed -i "s/SELINUX=.*/SELINUX=disabled/g" /etc/sysconfig/selinux
}

fix_python ()
{
	sudo unlink /bin/python
	sudo unlink /bin/pip
	sudo ln -s /bin/python3 /bin/python
	sudo ln -s /bin/pip3 /bin/pip
}	

prep ()
{
	user=centos
	#read -ep "What user do you want to access?" user
	ip_addr=$1
	if [ $1 == "" ]
	then
		read -ep "What IP do you want to access?" ip_addr
	fi
	echo "$0"
	scp -i $KEY_LOC $0 $user@$ip_addr:/home/centos/$SCRIPT_NAME
	scp -r -i $KEY_LOC $REMOTE_KEY_LOC $user@$ip_addr:/home/centos/remote-aws
	scp -r -i $KEY_LOC conf $user@$ip_addr:/home/centos/conf

}

cycle ()
{
	prep $1
	ssh -i $KEY_LOC $user@$ip_addr "chmod u+x $SCRIPT_NAME;./$SCRIPT_NAME pickup"
}

case $1 in
	start)
		echo "Starting prep"
		cycle $2
	;;
	pickup)
		echo "Starting pickup"
		setup_deps
		get_git_stuff
		##fix_python
		install_git_deps
		copy_conf_files
		sudo reboot
	;;
	reload)
		prep $2
		reload $2
	;;
	copy)
		copy_conf_files
	;;
	reboot)
		reboot_box $2
	;;
	*)
		exit
	;;
esac
