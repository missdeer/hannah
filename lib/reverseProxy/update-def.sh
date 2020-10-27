#!/bin/bash
echo "EXPORTS" > rp.def
grep -r -a '//export' main.go | while read line
do
    fn=`echo $line | awk -F ' ' '{ print $2 }'`
    echo "    $fn" >> rp.def
done

sed -i.bak '/#line/d' librp.h
sed -i.bak '/Complex/d' librp.h
sed -i.bak '/SIZE_TYPE/d' librp.h
rm -f librp.h.bak
