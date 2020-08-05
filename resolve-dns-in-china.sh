#!/bin/bash
grep -r -a 'http:\/\/[a-zA-Z0-9\.]\+\/' -o --include='*.go' . | awk -F '//' '{print $2}' | sort -n | uniq | awk -F '/' '{print $1}' | while read domain
do 
    ip=`curl "http://119.29.29.29/d?dn=$domain" -s --socks5 127.0.0.1:23333`; 
    firstip=`echo $ip | awk -F ';' '{print $1}'`;  
    echo $firstip $domain; 
done
