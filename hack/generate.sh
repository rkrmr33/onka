#!/bin/sh

protoc --version

PROJECT_ROOT=$(cd $(dirname ${BASH_SOURCE})/..; pwd)

# Building protoc-gen-gogofast from vendor directory because it is locked to
# a specific vesion in go.sum thanks to tools.go package
export PATH=$PATH:./dist
go mod vendor
GRPCPROTOBINARY=go-grpc

PROTO_FILES=$(find $PROJECT_ROOT -path $PROJECT_ROOT/vendor -prune -false -o -name '*.proto' | sort)

protoc \
    -I${PROJECT_ROOT} \
    -I/usr/local/include \
    -I./vendor \
    -I$GOPATH/src \
    -I${PROJECT_ROOT}/vendor/github.com/gogo/protobuf/types \
    --${PROTO_GEN}_out=\
Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,\
Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types,\
plugins=grpc:$GOPATH/src \
    $PROTO_FILES

rm -r vendor