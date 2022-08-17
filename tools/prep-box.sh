#!/bin/bash

KEY_LOC=~/.ssh/LightsailDefaultKey-us-west-2.pem
REMOTE_KEY_LOC=/opt/gametime/config/remote-aws
REMOTE_WORKDIR=/opt/gametime
SCRIPT_NAME=prep-box.sh
USER=ec2-user

# reboot_box ()
# {
# 	ip_addr=$1
# 	ssh -i $KEY_LOC $USER@$ip_addr "sudo reboot"
# }

server_startup ()
{
  pushd /home/$USER
  cp .secret gametime/
  cd gametime
  git pull
  git submodule update

  make box.dev.down
  make box.docker.dev
  make run

  popd
}

get_git_stuff ()
{
	pushd /home/$USER
#	eval `ssh-agent -s`
#	ssh-add config/remote-aws/id_rsa
#	ssh-keyscan -t rsa -H github.com >> /home/$USER/.ssh/known_hosts
	sleep 5
	cd ~

	git clone git@github.com:habruzzo/gametime.git
	cd gametime
	git submodule init
  git submodule update
  make get-dgraph
	sleep 5
	popd
}

setup_deps ()
{
#	sudo yum -y -q install yum-plugin-copr
#	sudo yum -y -q copr enable @caddy/caddy
#	sudo yum -y -q install caddy
	sudo yum -y -q install docker git go make
	sudo yum -y -q update
    #sudo sed -i "s/SELINUX=.*/SELINUX=disabled/g" /etc/sysconfig/selinux
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
	scp -r -i $KEY_LOC config $USER@$ip_addr:~/config
	scp -i $KEY_LOC .secret $USER@$ip_addr:~/.secret
}

cycle ()
{
  	ip_addr=$1
  	if [ $1 == "" ]
  	then
  		read -ep "What IP do you want to access?" ip_addr
  	fi
  	echo "$0"
	ssh -i $KEY_LOC $USER@$ip_addr "chmod u+x $SCRIPT_NAME;./$SCRIPT_NAME $2"
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
	;;
	reboot)
		reboot_box $2
	;;
	restart)
		prep $2
		cycle $2 "pickup-restart $2"
	;;
	pickup-restart)
		server_startup $2
	;;
	*)
	  echo 'start or restart; start <box ip> or restart <box ip>'
		exit
	;;
esac
exit