cd /root/helloweb/myFrame

echo "start save to workspace"
git add .
echo "save to workspace successfully"

echo "start commit to cache"
git commit -m "gitShellTest"
echo "commit to cache successfully"

echo "start pull the latest item"
git pull
echo "pull the latest item successfully"

echo "start push the latest item"
git push
echo "push the latest item successfully"