#!/bin/sh
set -e

if [[ ! -z "${GO_FLAGS}" ]]; then
    echo Building \"${OUT_FILE}\" with flags: \"${GO_FLAGS}\" starting at: \"${MAIN}\"
    for d in ${GO_FLAGS}; do
        export $d
    done
fi

go build -ldflags=" \
    -extldflags '-static'
    -X 'github.com/rkrmr33/onka/common.BinaryName=${BINARY_NAME}'
    -X 'github.com/rkrmr33/onka/common.Version=${VERSION}'
    " \
    -v -o ${OUT_FILE} ${MAIN}
