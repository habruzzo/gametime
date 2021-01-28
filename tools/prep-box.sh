#!/bin/bash
CONTAINER_URL=https://download.docker.com/linux/debian/dists/buster/pool/stable/amd64/containerd.io_1.4.3-1_amd64.deb
DOCKER_CLI_URL=https://download.docker.com/linux/debian/dists/buster/pool/stable/amd64/docker-ce-cli_20.10.2~3-0~debian-buster_amd64.deb
DOCKER_CE_URL=https://download.docker.com/linux/debian/dists/buster/pool/stable/amd64/docker-ce_20.10.2~3-0~debian-buster_amd64.deb

CONTAINER_PATH=containerd.io_1.4.3-1_amd64.deb
DOCKER_CLI_PATH=docker-ce-cli_20.10.2~3-0~debian-buster_amd64.deb
DOCKER_CE_PATH=docker-ce_20.10.2~3-0~debian-buster_amd64.deb
KEY_LOC=~/.ssh/LightsailDefaultKey-us-west-2.pem
REMOTE_KEY_LOC=/opt/gametime/conf/remote-aws
SCRIPT_NAME=prep-box.sh
USER=ec2-user


finish_server_startup ()
{
	pushd /opt/holdongametime
	source django/bin/activate
	#python manage.py runserver 0.0.0.0:8000 &
	sudo service httpd start
	deactivate
}

reboot_box ()
{
	ip_addr=$1
	ssh -i $KEY_LOC $USER@$ip_addr "sudo reboot"
}

reload ()
{
	pushd /opt/gametime/holdongametime
	ip_addr=$1
	scp -i $KEY_LOC holdongametime/* $USER@$ip_addr:/opt/holdongametime/holdongametime
	popd
}

copy_conf_files ()
{
	cd /etc/httpd/conf
	sudo chmod 755 *
	sudo cp /home/$USER/conf/httpd.conf /etc/httpd/conf/httpd.conf

	sudo service httpd restart
}

install_git_deps ()
{	
	pushd /opt/holdongametime
	python3 -m venv django
	source django/bin/activate
	pip install -r /home/$USER/gametime/requirements.txt
	sudo chgrp -R ec2-user /usr/lib64/httpd/
	sudo chmod -R g+w /usr/lib64/httpd/modules


	curl -L -O https://github.com/GrahamDumpleton/mod_wsgi/archive/4.7.1.tar.gz
	tar -xvzf 4.7.1.tar.gz
	cd mod_wsgi-4.7.1
	./configure --with-python=/opt/holdongametime/django/bin/python
	make
	sudo make install
	deactivate
	#npm install lessc
	popd
}

get_git_stuff ()
{
	pushd /home/$USER
	eval `ssh-agent -s`
	ssh-add conf/remote-aws/id_ed25519
	ssh-keyscan -t rsa -H github.com >> /home/$USER/.ssh/known_hosts	
	sleep 5
	
	git clone git@github.com:habruzzo/gametime.git
	
	sudo mv gametime/holdongametime /opt
	sudo chmod -R 775 /opt

	sudo ln -s /opt/holdongametime gametime/holdongametime
	#sudo chgrp -R apache /opt

	sudo mkdir /srv/http
	sudo chmod -R 775 /srv
	sudo ln -s /opt/holdongametime/static /srv/http/static
	sudo ln -s /opt/holdongametime/templates /srv/http/templates
	popd
}

setup_deps ()
{
	sudo yum -y -q install docker httpd httpd-devel mod_ssl openssl git npm gcc python3 python3-devel python3-pip
	sudo yum -y -q update
	#sudo yum -y install epel-release
    #sudo sed -i "s/SELINUX=.*/SELINUX=disabled/g" /etc/sysconfig/selinux
}

fix_python ()
{
	sudo update-alternatives --install /usr/bin/python python /usr/bin/python3 1
	sudo update-alternatives --install /usr/bin/pip pip /usr/bin/pip3 1

}	

prep ()
{
	ip_addr=$1
	if [ $1 == "" ]
	then
		read -ep "What IP do you want to access?" ip_addr
	fi
	echo "$0"
	scp -i $KEY_LOC $0 $USER@$ip_addr:~/$SCRIPT_NAME
	scp -r -i $KEY_LOC conf $USER@$ip_addr:~/conf

}

cycle ()
{
	prep $1
	ssh -i $KEY_LOC $USER@$ip_addr "chmod u+x $SCRIPT_NAME;./$SCRIPT_NAME $2"
}

case $1 in
	start)
		echo "Starting prep"
		cycle $2 "pickup"
	;;
	pickup)
		echo "Starting pickup"
		setup_deps
		get_git_stuff
		#fix_python
		install_git_deps
		copy_conf_files
		sudo reboot
		exit
	;;
	reload)
		reload $2
		cycle $1 "copy"
	;;
	copy)
		copy_conf_files
	;;
	reboot)
		reboot_box $2
	;;
	finish)
		reload $2
		cycle $2 "pickup-finish"
	;;
	# TODO: (26 Jan 20201) Fix terrible names
	pickup-finish)
		finish_server_startup
	;;
	*)
		exit
	;;
esac
