#!/bin/bash
grep -r -a 'https\?:\/\/[a-zA-Z0-9\.]\+\/' -o --include='*.go' . | awk -F '//' '{print $2}' | sort -n | uniq | awk -F '/' '{print $1}' | while read domain
do 
    ip=`curl "http://119.29.29.29/d?dn=$domain" -s $1 $2 $3`; 
    firstip=`echo $ip | awk -F ';' '{print $1}'`;  
    if [ -n "$firstip" ];
    then
        echo $firstip $domain; 
    fi
done
