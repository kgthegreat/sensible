# ps -aux | grep sensible

ip=46.101.51.200
user=kgthegreat
project_root=/home/kgthegreat/sensible

echo "Make sure you have stopped the server otherwise deploy will fail"
echo "Stopping server..."
ssh -t $user@$ip "sudo service sensible stop"
env GOOS=linux GOARCH=amd64 go build
echo "Binary built ..."
scp sensible $user@$ip:$project_root
echo "Binary copied to server ..."
rsync -av static $user@$ip:$project_root
echo "Statics copied to server ..."
rsync -av templates $user@$ip:$project_root
echo "templates copied to server ..."
rsync -av token.json $user@$ip:$project_root
echo "Token copied to server ..."
rsync -av keyword_template.json $user@$ip:$project_root
echo "Keyword template copied to server ..."
rsync -av keyword.json $user@$ip:$project_root
echo "Keyword seed copied to server ..."
echo "Restarting app..."
ssh -t $user@$ip "sudo service sensible start"
echo "App restarted"
echo "Deployment complete. Please go to http://trysensible.com"
