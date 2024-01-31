#!/bin/bash

# Function to pause and wait for user input
function pause() {
    read -n1 -rsp $'Press any key to continue...\n'
}

echo "./wget https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg"
pause
./wget https://pbs.twimg.com/media/EMtmPFLWkAA8CIS.jpg
pause
clear

echo "./wget https://reboot01.com"
pause
./wget https://reboot01.com
pause
clear

echo "./wget https://golang.org/dl/go1.16.3.linux-amd64.tar.gz"
pause
./wget https://golang.org/dl/go1.16.3.linux-amd64.tar.gz
pause
clear

echo "./wget http://ipv4.download.thinkbroadband.com/100MB.zip"
pause
./wget http://ipv4.download.thinkbroadband.com/100MB.zip
pause
clear

echo "./wget -O=test_20MB.zip http://ipv4.download.thinkbroadband.com/20MB.zip"
pause
./wget -O=test_20MB.zip http://ipv4.download.thinkbroadband.com/20MB.zip
pause
clear

echo "./wget -O=test_20MB.zip -P=~/Downloads/ http://ipv4.download.thinkbroadband.com/20MB.zip"
pause
./wget -O=test_20MB.zip -P=~/Downloads/ http://ipv4.download.thinkbroadband.com/20MB.zip
pause
clear

echo "./wget --rate-limit=300k http://ipv4.download.thinkbroadband.com/20MB.zip"
pause
./wget --rate-limit=300k http://ipv4.download.thinkbroadband.com/20MB.zip
pause
clear

echo "./wget --rate-limit=700k http://ipv4.download.thinkbroadband.com/20MB.zip"
pause
./wget --rate-limit=700k http://ipv4.download.thinkbroadband.com/20MB.zip
pause
clear

echo "./wget --rate-limit=2M http://ipv4.download.thinkbroadband.com/20MB.zip"
pause
./wget --rate-limit=2M http://ipv4.download.thinkbroadband.com/20MB.zip
pause
clear

echo "./wget -i=downloads.txt"
pause
./wget -i=downloads.txt
pause
clear

echo "./wget -B http://ipv4.download.thinkbroadband.com/20MB.zip"
pause
./wget -B http://ipv4.download.thinkbroadband.com/20MB.zip
pause
clear

echo "./wget --mirror --convert-links http://corndog.io/"
pause
./wget --mirror --convert-links http://corndog.io/
pause
clear

echo "./wget --mirror https://oct82.com/"
pause
./wget --mirror https://oct82.com/
pause
clear

rm -r oct82.com
echo "./wget --mirror --reject=gif https://oct82.com/"
pause
./wget --mirror --reject=gif https://oct82.com/
pause
clear

echo "./wget --mirror https://trypap.com/"
pause
./wget --mirror https://trypap.com/
pause
rm -r trypap.com
clear


echo "./wget --mirror -X=/img https://trypap.com/"
pause
./wget --mirror -X=/img https://trypap.com/
pause
clear

echo "./wget --mirror https://theuselessweb.com/"
pause
./wget --mirror https://theuselessweb.com/
pause
clear

echo "./wget --mirror <https://link_of_your_choice.com>"
pause
