#!/bin/sh
cd `dirname $0`
rm -rf ./config/data/output 
mkdir ./config/data/output 
# cp defs/errno.go ./config/data/
excel_export $1 strict ./config/data/
cp ./config/data/output/*.sql ./config
cp -rf ./config/data/output/db_defs_design_select.go ./design/
sed  -i '' 's|package db|package design|g' ./design/db_defs_design_select.go


