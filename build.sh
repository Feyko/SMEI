#!/bin/bash

go build -o smei.exe
go build -o smei-companion.exe cmd/auth-companion/proxy/auth-proxy.go