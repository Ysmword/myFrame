cd /root/helloweb/myFrame

echo "set git config"
git config --global user.email "1843121593@qq.com"
git config --global user.name "Ysmword"
echo "set git config successfully"

echo "start save to workspace"
git add .
echo "save to workspace successfully"

echo "start commit to cache"
git commit -m "[modify]:simple gitshellTest"
echo "commit to cache successfully"

echo "start pull the latest item"
git pull
echo "pull the latest item successfully"

echo "start push the latest item"
git push origin master
echo "push the latest item successfully"