#!/bin/sh

#保持容器不退出，用于调试

while true;do
   sleep 5
   echo "sleep 5" >> $HOME/logs/idle.log 2>&1
done
