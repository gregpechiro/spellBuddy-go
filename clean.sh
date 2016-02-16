#!/bin/bash

# a script to remove old files and folders
# this should only be run on the server where the project lives right before
# redeployment. It will remove all files that are contained in the new .tar

DIR=`echo ${PWD##*/}`

rm -rf static/ templates/ *.go ${DIR} ${DIR}.tar
