if [[ $DOWNLOAD_LINK == "" ]]
then
	DOWNLOAD_LINK=$1
	if [[ $DOWNLOAD_LINK == "" ]]
	then
		echo "dont forget to supply a download link with export \"DOWNLOAD_LINK=<link>\" or ./get-gecko.sh <link>"
		exit
	fi
fi
#DOWNLOAD_LINK=curl -L -O https://github.com/mozilla/geckodriver/releases/download/v0.28.0/geckodriver-v0.28.0-linux64.tar.gz
cd /home/holden/Downloads/zips
GECKO_FILE=$(find . -name "geckodriver-*-linux64.tar.gz")
echo $GECKO_FILE
read -r -p "Thats the geckodriver you currently have. Keep going? (y/n)" answer1
case $answer1 in 
	y|Y|yes)
	curl -L -O $DOWNLOAD_LINK
	if [ $? -ne 0 ]
		then echo "oh no a problem downloading"
		exit
	fi
	QUERY_RET=($(find . -name "geckodriver-*-linux64.tar.gz"))
	GECKO_NEW=${QUERY_RET[0]}
	echo $GECKO_NEW
	read -r -p "We downloaded a new driver. Keep going?" answer2
	case $answer2 in
		y|Y|yes)
			tar -xvzf $GECKO_NEW
			cd -
			read -r -p "Keep going?" answer3
			case $answer3 in
				y|Y|yes)
					GECKO_FULL="/home/holden/Downloads/zips/geckodriver"
					ln -s $GECKO_FULL geckodriver
					exit
					;;
				*)
				echo "okay bye"
				exit
				;;
			esac
		;;
		*)
		echo "okay bye"
		exit
		;;
	esac
	;;
	*)
	echo "okay bye"
	exit
	;;
esac

