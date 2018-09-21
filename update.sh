#!/bin/sh

pm2 stop feedback
rm feedback.exe
go build -o feedback.exe
pm2 start feedback.exe

while true; do
    
    if [[ $(git pull origin master) == *up-to-date* ]]; 
    then
        echo "no change"
    else
        echo "detect changes"
        sleep 2s
        echo "stop app"
		pm2 stop feedback
		rm feedback.exe
        go build -o feedback.exe
        pm2 restart feedback
    fi

    echo "sleep 30s"

    sleep 30s        

done
