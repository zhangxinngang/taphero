#!/bin/sh
rm -rf ~/tmp
mkdir ~/tmp
wget -q ftp://www.chaoshenxy.com/th/taphero.tar.bz2 -O ~/tmp/taphero.tar.bz2
mkdir -p ~/bin
mkdir -p /data/taphero/config/
tar -xf ~/tmp/taphero.tar.bz2 -C ~/tmp/
cp -rf ~/tmp/dist/taphero ~/bin/
cp -rf ~/tmp/dist/config/* /data/taphero/config/
chmod 755 ~/bin/taphero
if [ -d ~/go/bin/ ]; then
    rm -rf  ~/go/bin/taphero
    ln -f -s ~/bin/taphero ~/go/bin/taphero
fi