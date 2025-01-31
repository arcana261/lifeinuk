#!/bin/bash

set1='abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ'
set2='0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
set3='IJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCDEFGH'
set4='efghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcd'
set5='wxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuv'
set6='JKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHI'
cmd='encode'

while [[ $# > 0 ]]; do
    case $1 in
    -e | --encode)
        cmd='encode'
        ;;
    -d | --decode)
        cmd='decode'
        ;;
    *)
        echo "unknown cmd: $1" > /dev/stderr
        exit -1
        ;;
    esac
    shift
done

function perform_encode() {
    cat highlights.txt | tr $set1 $set2 | base64 -w 0 | tr $set3 $set4 | gzip -9 | base64 -w 0 | tr $set5 $set6 > highlights.enc
}

function perform_decode() {
    cat highlights.enc | tr $set6 $set5 | base64 -w 0 -d | gunzip | tr $set4 $set3 | base64 -w 0 -d | tr $set2 $set1 > highlights.txt
}

if [ "$cmd" = "encode" ]; then
    echo '>> encoding...'
    cp highlights.enc highlights.enc.bak
    cp highlights.txt highlights.txt.bak
    perform_encode
    echo '>> testing...'
    perform_decode
    diff highlights.txt highlights.txt.bak
else
    echo '>> decoding...'
    cp highlights.txt highlights.txt.bak
    perform_decode
    head highlights.txt
fi