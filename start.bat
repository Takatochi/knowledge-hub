@echo off
cd %~dp0\deployments
docker-compose --env-file ..\.env %*
