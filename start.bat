@echo off
cd %~dp0\deployments
docker-compose -p knowledge-hub --env-file ..\.env %*
