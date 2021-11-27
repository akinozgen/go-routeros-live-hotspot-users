#!/bin/bash#!/bin/bash

# MySQL Connection Info
export MYSQL_PWD=your_mysql_password
export MYSQL_USER=your_mysql_user
export MYSQL_HOST_PORT=your_mysql_host_port
export MYSQL_DB=your_mysql_db

# Router Connection Info
export ROUTER_IP_PORT=your_router_ip_port
export ROUTER_USER=your_router_user
export ROUTER_PWD=your_router_pwd

# Parameters
export PRINT_PARAMETERS="server,user,address,uptime,bytes-in,bytes-out"
export WATCH_ONLINES=true

go run . --mysql-user=$MYSQL_USER \
         --mysql-host-port=$MYSQL_HOST_PORT \
         --mysql-pwd=$MYSQL_PWD \
         --mysql-db=$MYSQL_DB \
         --router-ip-port=$ROUTER_IP_PORT \
         --router-user=$ROUTER_USER \
         --router-pwd=$ROUTER_PWD \
         --print-parameters=$PRINT_PARAMETERS \
         --watch-onlines=$WATCH_ONLINES