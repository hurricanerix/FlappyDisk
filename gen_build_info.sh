#!/bin/bash

DATE=`date`
BUILD_HASH=`git log -n 1 | grep commit | cut -d " " -f 2`

CODE="// gen is a generated package, DO NOT EDIT!\n
\n
package gen\n
\n
var BuildDate string = \"$DATE\"\n
var BuildHash string = \"$BUILD_HASH\"\n
"

# go-bindata is required to embed assets into binary
go get -u github.com/jteeuwen/go-bindata/...

# generate all the files we need
mkdir -p gen
go-bindata -pkg gen -ignore="/*.pyxel" -o gen/assets.go assets/
echo -e $CODE | gofmt > gen/build_info.go
