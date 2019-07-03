# ps -aux | grep sensible

echo "Make sure you have stopped the server otherwise deploy will fail"
env GOOS=linux GOARCH=amd64 go build
echo "Binary built ..."
scp sensible kgthegreat@46.101.51.200:/home/kgthegreat/sensible
echo "Binary copied to server ..."
rsync -av static kgthegreat@46.101.51.200:/home/kgthegreat/sensible
echo "Statics copied to server ..."
rsync -av templates kgthegreat@46.101.51.200:/home/kgthegreat/sensible
echo "templates copied to server ..."
rsync -av token.json kgthegreat@46.101.51.200:/home/kgthegreat/sensible
echo "Token copied to server ..."
#rsync -av keyword_template.json kgthegreat@46.101.51.200:/home/kgthegreat/sensible
#echo "Keyword template copied to server ..."
#rsync -av keyword.json kgthegreat@46.101.51.200:/home/kgthegreat/sensible
#echo "Keyword seed copied to server ..."
echo "Deployment complete."
