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
set23='wxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuv'
set24='xyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvw'
set25='yz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwx'
set26='z0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxy'
set27='0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
set28='123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0'
set29='23456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01'
set30='3456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz012'
set31='456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123'
set32='56789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234'
set33='6789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz012345'
set34='789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456'
set35='89ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz01234567'
set36='9ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz012345678'
set37='ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'
set38='BCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789A'
set39='CDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789AB'
set40='DEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABC'
set41='EFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCD'
set42='FGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCDE'
set43='GHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCDEF'
set44='HIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCDEFG'
set45='IJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCDEFGH'
set46='JKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHI'
set47='KLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJ'
set48='LMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJK'

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
    local pass1="$1"
    local pass2=""

    if [ "$pass1" = "" ]; then
        read -s -p "Password: " pass1
        echo ""
        read -s -p "Confirm: " pass2
        echo ""
        if [ "$pass1" != "$pass2" ]; then
            echo "Passwords do not match!!!"
            exit -1
        fi
    fi

    tar cvzO data |
        base64 -w 0 | tr $set1 $set2 | gzip -9 | base64 -w 0 | tr $set3 $set4 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set5 $set6 | gzip -9 | base64 -w 0 | tr $set7 $set8 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set9 $set10 | gzip -9 | base64 -w 0 | tr $set11 $set12 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set13 $set14 | gzip -9 | base64 -w 0 | tr $set15 $set16 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set17 $set18 | gzip -9 | base64 -w 0 | tr $set19 $set20 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set21 $set22 | gzip -9 | base64 -w 0 | tr $set23 $set24 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set25 $set26 | gzip -9 | base64 -w 0 | tr $set27 $set28 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set29 $set30 | gzip -9 | base64 -w 0 | tr $set31 $set32 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set33 $set34 | gzip -9 | base64 -w 0 | tr $set35 $set36 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set37 $set38 | gzip -9 | base64 -w 0 | tr $set39 $set40 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set41 $set42 | gzip -9 | base64 -w 0 | tr $set43 $set44 | \
        openssl aes-256-cbc -salt -pbkdf2 -pass "pass:$pass1" | \
        base64 -w 0 | tr $set45 $set46 | gzip -9 | base64 -w 0 | tr $set47 $set48 | \
        tr '0' '\n' > data.enc

    unset pass1
    unset pass2
}

function perform_decode() {
    local pass1="$1"

    if [ "$pass1" = "" ]; then
        read -s -p "Password: " pass1
        echo ""
    fi

    cat data.enc | \
        tr '\n' '0' | \
        tr $set48 $set47 | base64 -w 0 -d | gunzip | tr $set46 $set45 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set44 $set43 | base64 -w 0 -d | gunzip | tr $set42 $set41 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set40 $set39 | base64 -w 0 -d | gunzip | tr $set38 $set37 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set36 $set35 | base64 -w 0 -d | gunzip | tr $set34 $set33 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set32 $set31 | base64 -w 0 -d | gunzip | tr $set30 $set29 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set28 $set27 | base64 -w 0 -d | gunzip | tr $set26 $set25 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set24 $set23 | base64 -w 0 -d | gunzip | tr $set22 $set21 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set20 $set19 | base64 -w 0 -d | gunzip | tr $set18 $set17 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set16 $set15 | base64 -w 0 -d | gunzip | tr $set14 $set13 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set12 $set11 | base64 -w 0 -d | gunzip | tr $set10 $set9 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set8 $set7 | base64 -w 0 -d | gunzip | tr $set6 $set5 | base64 -w 0 -d | \
        openssl aes-256-cbc -d -pbkdf2 -pass "pass:$pass1" | \
        tr $set4 $set3 | base64 -w 0 -d | gunzip | tr $set2 $set1 | base64 -w 0 -d | \
        tar xzv

    unset pass1
}

if [ "$cmd" = "encode" ]; then
    echo '>> encoding...'
    touch data.enc
    rm -rf data.bak data.enc.bak
    mv data.enc data.enc.bak

    pass1=""
    pass2=""
    read -s -p "Password: " pass1
    echo ""
    read -s -p "Confirm: " pass2
    echo ""
    if [ "$pass1" != "$pass2" ]; then
        echo "Passwords do not match!!!"
        exit -1
    fi

    perform_encode "$pass1"
    echo '>> testing...'
    mv data data.bak
    perform_decode "$pass1"
    diff data/highlights.txt data.bak/highlights.txt
    if [ "$?" != "0" ]; then
        rm -rf data
        mv data.bak data
    fi

    unset pass1
    unset pass2
else
    echo '>> decoding...'
    rm -rf data.bak
    mkdir -p data
    mv data data.bak
    perform_decode
    head data/highlights.txt
fi
