#!/bin/bash
set -x

export THISPROJ=snap-plugin-publisher-elasticsearch
export GOPATH=$WORKSPACE/$THISPROJ/gopath
export IMPORTPATH=github.com/intelsdi-x/snap-plugin-lib-go
rm -rf $GOPATH && mkdir -p $GOPATH/bin $GOPATH/pkg $GOPATH/src/$IMPORTPATH && mv $THISPROJ/src/* $GOPATH/src/$IMPORTPATH
go get $IMPORTPATH/...
cd $GOPATH/snap-plugin-lib-go
glide up
# go test $IMPORTPATH || exit 1
rm -rf dist && mkdir dist && cp -a $GOPATH/bin/* $GOPATH/src/$IMPORTPATH/{js,font,css,html,img} dist/
