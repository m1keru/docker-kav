#!/bin/bash
go build
sudo cp docker-kav /usr/local/bin/
sudo cp docker-kav.service /etc/systemd/system/
#sudo cp avcheck /srv/etp/configs/config.d/

