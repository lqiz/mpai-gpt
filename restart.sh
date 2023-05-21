git pull
rm -f mpai
go build
killall -w mpai
nohup ./mpai -c config.toml >> /data/logs/mpai/mpai.log 2>&1 &
