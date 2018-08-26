#!/usr/bin/env bash

cd $(dirname $0)
cd ../../

protoc --go_out=plugins=grpc:pb proto/gfalcon.proto
