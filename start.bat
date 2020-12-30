@echo off
go build -tags=jsoniter
set config_path=conf/config.yml
proxy.exe