#!/bin/bash

#启动docker

docker run -d -v /tmp/multichain-explorer-mock-`date +'%Y%m%d%H%M%S'`:/home/jovyan/logs -p 28080:28080 -t lansheng228/multichain-explorer-mock:latest