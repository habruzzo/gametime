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


spin ()
{
	for i in {1..12}
	do
		for j in {1..5}
		do
			sleep 1
		done	
		echo "."
	done
}

finish_server_startup ()
{
	copy_conf_files $1
	pushd /opt/holdongametime
	sudo chmod 774 logs/django.log
	sudo chgrp apache logs/django.log
	source django/bin/activate
	#gunicorn -D -b "127.0.0.1:8000" holdongametime.wsgi --log-level DEBUG --log-file error.log

	python manage.py runserver 0.0.0.0:8000 &
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
	pushd /opt/holdongametime
	# grep "ALLOWED_HOSTS" holdongametime/settings.py
	# sed -i "s/%IP_ADDR%/$1/g" holdongametime/settings.py
	# grep "ALLOWED_HOSTS" holdongametime/settings.py
	# popd
	# pushd /etc/httpd/conf
	# sudo chmod 755 *
	# sudo cp /home/$USER/conf/httpd.conf /etc/httpd/conf/httpd.conf


	# sudo cp -v /home/$USER/conf/httpd.conf /etc/httpd/conf
	# #sudo cp -v /home/$USER/conf/Caddyfile /etc/caddy
	# sudo service httpd restart
	popd
}

install_git_deps ()
{	
	pushd /home/$USER/gametime/revel
	export GOPATH="/home/$USER/gametime/revel"
	export PATH=$PATH:$GOPATH/bin
	go get "github.com/revel/revel"
	go get "github.com/revel/cmd/revel"


	#pushd /opt/holdongametime

	#python3 -m venv venv
	#source venv/bin/activate
	#pip install -r /home/$USER/gametime/requirements.txt
	#sudo chgrp -R ec2-user /usr/lib64/httpd/
	#sudo chmod -R g+w /usr/lib64/httpd/modules

	#curl -L -O https://github.com/GrahamDumpleton/mod_wsgi/archive/4.7.1.tar.gz
	#tar -xvzf 4.7.1.tar.gz
	#cd mod_wsgi-4.7.1
	#./configure --with-python=/opt/holdongametime/django/bin/python
	#make
	#sudo make install
	#deactivate
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
	cd gametime
	git clone git@github.com:habruzzo/reviews.git
	cd revel

	#sudo mv gametime/revel/blog_holdongametime /opt
	#sudo chmod -R 775 /opt

	#sudo ln -s /opt/blog_holdongametime gametime/revel/blog_holdongametime
	#sudo chgrp -R apache /opt

	#sudo mkdir /srv/http
	#sudo chmod -R 775 /srv
	#sudo ln -s /opt/holdongametime/static /srv/http/static
	#sudo ln -s /opt/holdongametime/templates /srv/http/templates
	popd
}

setup_deps ()
{
	sudo yum -y -q install yum-plugin-copr
	sudo yum -y -q copr enable @caddy/caddy
	sudo yum -y -q install caddy
	#sudo yum -y -q install docker git npm gcc python3 python3-devel python3-pip httpd httpd-devel mod_ssl openssl 
	sudo yum -y -q install docker git gcc golang mercurial postgres psql
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
	ssh -i $KEY_LOC $USER@$1 "chmod u+x $SCRIPT_NAME;./$SCRIPT_NAME $2"
}

case $1 in
	start)
		echo "Starting prep"
		prep $2
		cycle $2 "pickup"
	;;
	pickup)
		echo "Starting pickup"
		setup_deps
		get_git_stuff
		#fix_python
		install_git_deps
		copy_conf_files
		#sudo reboot
		exit
	;;
	reload)
		reload $2
		prep $2
		cycle $2 "copy $2"
	;;
	copy)
		copy_conf_files $2
	;;
	reboot)
		reboot_box $2
	;;
	finish)
		#spin $2
		reload $2
		prep $2
		cycle $2 "pickup-finish $2" 
	;;
	# TODO: (26 Jan 20201) Fix terrible names
	pickup-finish)
		
		finish_server_startup $2
	;;
	*)
		exit
	;;
esac
exit