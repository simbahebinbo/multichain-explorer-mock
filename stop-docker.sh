#!/bin/bash

#停掉正在运行的docker容器


docker stop $(docker ps -qa --filter ancestor=lansheng228/multichain-explorer-mock)
docker rm $(docker ps -qa --filter ancestor=lansheng228/multichain-explorer-mock)
