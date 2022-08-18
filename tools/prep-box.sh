#!/bin/bash

KEY_LOC=~/.ssh/LightsailDefaultKey-us-west-2.pem
REMOTE_KEY_LOC=/opt/gametime/config/remote-aws
REMOTE_WORKDIR=/opt/gametime
SCRIPT_NAME=prep-box.sh
USER=ec2-user
IP_ADDR="52.88.59.140"
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
  sudo systemctl start docker
  sudo chmod 666 /var/run/docker.sock
  make box.dev.down
  make box.docker.dev
  make run &

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
	sudo curl -L https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m) -o /usr/local/bin/docker-compose
	sudo chmod +x /usr/local/bin/docker-compose
    #sudo sed -i "s/SELINUX=.*/SELINUX=disabled/g" /etc/sysconfig/selinux
}

prep ()
{
  echo "Starting prep"
  if [[ $IP_ADDR == "" ]]
  then
    if [[ $1 == "" ]]
    then
      read -ep "What IP do you want to access?" IP_ADDR;
    else IP_ADDR=$1
    fi
  fi
	scp -i $KEY_LOC $0 $USER@$IP_ADDR:~/$SCRIPT_NAME
	scp -r -i $KEY_LOC config $USER@$IP_ADDR:~/config
	scp -i $KEY_LOC .secret $USER@$IP_ADDR:~/.secret
}

cycle ()
{
	ssh -i $KEY_LOC $USER@$IP_ADDR "./$SCRIPT_NAME $1"
}

case $1 in
	start)
		prep $2
		cycle "pickup"
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
		cycle "pickup-restart"
	;;
	pickup-restart)
		server_startup
	;;
	*)
	  echo 'start or restart; start <box ip> or restart <box ip>'
		exit
	;;
esac
exit