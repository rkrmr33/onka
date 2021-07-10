// +build tools

// This package is used to lock the versions of certain tools being used for code
// generation and testing
package tools

import (
	_ "github.com/gogo/protobuf/protoc-gen-gofast"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
)
