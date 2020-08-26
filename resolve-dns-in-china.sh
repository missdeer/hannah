#!/bin/bash
cat /etc/hosts | grep "\(.migu.cn\|.163.com\|.kugou.com\|kuwo.cn\|.xiami.com\|.qq.com\|.kgimg.com\|.bilibili.com\)$" | while read line; do echo $line | awk -F ' ' '{print $2}'; done > domainlist.txt
grep -r -a 'https\?:\/\/[a-zA-Z0-9\.]\+\/' -o --include='*.go' . | awk -F '//' '{print $2}' | awk -F '/' '{print $1}' | grep "\(.migu.cn\|.163.com\|.kugou.com\|kuwo.cn\|.xiami.com\|.qq.com\|.kgimg.com\|.bilibili.com\)$"  >>domainlist.txt 
cat domainlist.txt | sort -n | uniq | while read domain
do 
    ip=`curl "http://119.29.29.29/d?dn=$domain" -s $1 $2 $3`; 
    firstip=`echo $ip | awk -F ';' '{print $1}'`;  
    if [ -n "$firstip" ];
    then
        printf "%-15s %s\n" $firstip $domain; 
    fi
done > result.txt
cat result.txt | sort -n
rm -f domainlist.txt result.txt 
