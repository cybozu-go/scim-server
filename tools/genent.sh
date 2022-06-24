#!/bin/bash

pushd tools/cmd/genent
go build -o .genent main.go
popd

./tools/cmd/genent/.genent $GENENT_ARGS
rm ./tools/cmd/genent/.genent

