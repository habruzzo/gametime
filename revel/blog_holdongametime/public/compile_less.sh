LESS_END=".less"
CSS_END=".css"

for i in $(find . -type f | grep "\.less" | cut -d "." -f2 | cut -b 2-)
do
	REALNAME=$(echo $i)
	JUSTFILE=$(echo $REALNAME | rev | cut -d "/" -f 1 | rev)
	echo $REALNAME$LESS_END $JUSTFILE$CSS_END
    lessc $REALNAME$LESS_END $JUSTFILE$CSS_END
done	
