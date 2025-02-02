#!/bin/bash

set1='abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ'
set2='bcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZa'
set3='cdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZab'
set4='defghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabc'
set5='efghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcd'
set6='fghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcde'
set7='ghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef'
set8='hijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefg'
set9='ijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefgh'
set10='jklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghi'
set11='klmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghij'
set12='lmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijk'
set13='mnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijkl'
set14='nopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklm'
set15='opqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmn'
set16='pqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmno'
set17='qrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnop'
set18='rstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopq'
set19='stuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqr'
set20='tuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrs'
set21='uvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrst'
set22='vwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstu'


# tar cvzO data > data.tar.gz

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
    tar cvzO data |
        base64 -w 0 | tr $set1 $set2 | base64 -w 0 | tr $set3 $set4 | gzip -9 | base64 -w 0 | tr $set5 $set6 | \
        base64 -w 0 | tr $set7 $set8 | gzip -8 | base64 -w 0 | tr $set9 $set10 | \
        base64 -w 0 | tr $set11 $set12 | gzip -7 | base64 -w 0 | tr $set13 $set14 | \
        base64 -w 0 | tr $set15 $set16 | gzip -6 | base64 -w 0 | tr $set17 $set18 | \
        base64 -w 0 | tr $set19 $set20 | gzip -5 | base64 -w 0 | tr $set21 $set22 > data.enc
}

function perform_decode() {
    cat data.enc | \
        tr $set22 $set21 | base64 -w 0 -d | gunzip | tr $set20 $set19 | base64 -w 0 -d | \
        tr $set18 $set17 | base64 -w 0 -d | gunzip | tr $set16 $set15 | base64 -w 0 -d | \
        tr $set14 $set13 | base64 -w 0 -d | gunzip | tr $set12 $set11 | base64 -w 0 -d | \
        tr $set10 $set9 | base64 -w 0 -d | gunzip | tr $set8 $set7 | base64 -w 0 -d | \
        tr $set6 $set5 | base64 -w 0 -d | gunzip | tr $set4 $set3 | base64 -w 0 -d | tr $set2 $set1 | base64 -w 0 -d | \
        tar xzv
}

if [ "$cmd" = "encode" ]; then
    echo '>> encoding...'
    touch data.enc
    cp data.enc data.enc.bak
    cp data/highlights.txt highlights.txt.bak
    perform_encode
    echo '>> testing...'
    perform_decode
    diff data/highlights.txt highlights.txt.bak
else
    echo '>> decoding...'
    cp data/highlights.txt highlights.txt.bak
    perform_decode
    head data/highlights.txt
fi