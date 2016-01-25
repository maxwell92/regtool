#! /bin/bash
# maintainer liyao.yao@yeepay.com

REG_DOMAIN="registry.test.com"
REG_PORT="5000"

CURL_OPTS="-skXGET"


curl $CURL_OPTS https://$REG_DOMAIN:$REG_DOMAIN/v2/_catalog | jq .







