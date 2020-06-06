
echo "build from source"
export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
cd /root/gzm/GZMApp/src
go build
echo "build complete"

echo "stop container gzm"
docker stop gzm
echo "delete container gzm"
docker rm gzm
echo "container gzm deleted"

echo "start new container"
# /root/gzm/GZMApp/src/ 这是文件放置的地方 暴露的端口为8080
docker run -d --name=gzm -v /root/gzm/GZMApp/src/:/workspace/ -p 8080:8080 alpine /workspace/gzm

echo "build done"
