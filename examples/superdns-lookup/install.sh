#!/bin/bash

go install ../superdns-generate
go install

superdns-generate
sudo cp -r ./superdns /var/
sudo cp ./superdns.conf /etc/superdns.conf
rm -rf ./superdns
rm ./superdns.conf

