#!/bin/bash

# Install paho.mqtt
go mod init paho.mqtt
go get github.com/eclipse/paho.mqtt.golang

# Install serial
go mod init serial
go get github.com/jacobsa/go-serial/serial
