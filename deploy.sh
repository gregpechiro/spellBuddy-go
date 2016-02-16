#!/bin/bash

clear
NODE="$1"
DIR=`echo ${PWD##*/}`
EXCLUDE=(".git" ".idea" "db" ".gitignore" "spellBuddy.tar" "deploy.sh")

if [ "$NODE" == "" ]; then
	echo "No node specified!"
	exit 1
fi

if [ -f "${DIR}.tar" ]; then
	echo "Removing old tar ${DIR}.tar..."
	rm $DIR.tar
fi

echo "Removing old binary ${DIR}..."
go clean
echo "Building ${DIR}..."
go build

if [ ! -f $DIR ]; then
	echo "Build $DIR failed."
	exit 1
fi

echo "Creating tar ${DIR}.tar..."

for item in ${EXCLUDE[*]}; do
	TOGETHER="$TOGETHER --exclude $item"
done

tar cf $DIR.tar * $TOGETHER
if [ ! -f "$DIR.tar" ]; then
	echo "Create $DIR.tar failed."
	exit 1
fi

echo "SCP to spellbuddy.xiphoid24.com..."
scp $DIR.tar greg@xiphoid24.com:/home/greg
