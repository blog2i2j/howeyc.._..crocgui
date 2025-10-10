#!/usr/bin/env bash

curl -Lso go.tar.gz https://go.dev/dl/go1.25.1.linux-amd64.tar.gz
echo "7716a0d940a0f6ae8e1f3b3f4f36299dc53e31b16840dbd171254312c41ca12e go.tar.gz" | sha256sum -c -
mkdir -p gobuild/go{lang,path,cache}
tar -C gobuild/golang -xzf go.tar.gz
rm go.tar.gz
export GOPATH="$PWD/gobuild/gopath"
export GOCACHE="$PWD/gobuild/gocache"
export GO_LANG="$PWD/gobuild/golang/go/bin"
export GO_COMPILED="$GOPATH/bin"
export PATH="$GO_LANG:$GO_COMPILED:$PATH"
go version
go install fyne.io/fyne/v2/cmd/fyne\@v2.6.3
fyne version
if [[ $# -eq 0 ]]; then
	fyne package -os android -release
	zip -d crocgui.apk "META-INF/*"
else
	fyne package -os android
fi
chmod -R u+w gobuild
rm -rf gobuild
