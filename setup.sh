#!/bin/bash 

sudo apt-get update
sudo apt-get upgrade 

# install mariadb
sudo apt-get install mariadb-server

sudo systemctl start mariadb
sudo systemctl status mariadb

# create user & password, and grant priviliges 
# TODO 

# /etc/mysql/mariadb.conf.d/50-server.cnf file has "bind-address = 127.0.0.1" line uncommented 

# create database and tables inside

# setup service using 8080 port

# setup nginx to redirect 80 port request to 8080

# config GoDaddy DNS config if have not done it

# test if server is okay using https://duftee.com/moma-api/rate?fromCode=CNY&toCode=USD