#!/bin/bash

# maintainer liyao.miao@yeepay.com

REG_DOMAIN="registry.test.com"
REG_PORT="5000"
CURL_OPT="-skXGET"

image=""
tag=""

#echo $1 > .arg
#image=`echo $1 | grep ^[a-zA-Z/]*`
#cat .arg
#image=`grep ^[a-zA-Z/]* .arg`
if [ $# -eq 0 ] 
then
    echo "Usage:"
    echo "./regsearch IMAGE[:TAG]"
fi

image=`echo $1 | cut -d ":" -f1`
tag=`echo $1 | cut -d ":" -f2`
#echo $image
#echo $tag

echo $1 | grep [:] > /dev/null 
if [ $? -ne 0 ]
then
    curl $CURL_OPT https://$REG_DOMAIN:$REG_PORT/v2/$image/tags/list | jq .   
else 
    #curl $CURL_OPT https://$REG_DOMAIN:$REG_PORT/v2/$1
    curl $CURL_OPT https://$REG_DOMAIN:$REG_PORT/v2/$image/tags/list | jq . | grep $tag > /dev/null

    if [ $? -eq 0 ]
    then
        echo $image:$tag "found!"
    else
        echo $image:$tag "NOT found!"
    fi
fi 
#rm -f .arg

