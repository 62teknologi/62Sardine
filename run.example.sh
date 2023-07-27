#run sardine
echo "killing port 10083"
sudo kill -9 $(sudo lsof -t -i:10083)  || true
echo "build sardine"
go build main.go
echo "running go sardine"
nohup ./main &
