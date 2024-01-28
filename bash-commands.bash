#!/bin/bash

# Function to pause and wait for user input
function pause() {
    read -n1 -rsp $'Press any key to continue...\n'
}

./wget https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
pause

./wget https://reboot01.com
pause

./wget https://golang.org/dl/go1.16.3.linux-amd64.tar.gz
pause

./wget http://ipv4.download.thinkbroadband.com/100MB.zip
pause

./wget -O=test_20MB.zip http://ipv4.download.thinkbroadband.com/20MB.zip
pause

./wget -O=test_20MB.zip -P=~/Downloads/ http://ipv4.download.thinkbroadband.com/20MB.zip
pause

./wget --rate-limit=300k http://ipv4.download.thinkbroadband.com/20MB.zip
pause

./wget --rate-limit=700k http://ipv4.download.thinkbroadband.com/20MB.zip
pause

./wget --rate-limit=2M http://ipv4.download.thinkbroadband.com/20MB.zip
pause

./wget -i=downloads.txt
pause

./wget -B http://ipv4.download.thinkbroadband.com/20MB.zip
pause
