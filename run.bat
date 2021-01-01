@echo off
set config_path=conf/config.yml
set enable_log=true
set log_path=./log/proxy.log
set log_level=INFO
proxy.exe t -c cert.pem -k key.pem