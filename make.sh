#!/usr/bin/env bash

if [[ ! -f make.sh ]]; then
	echo 'make.sh must be run from src' 1>&2
	exit 1
fi

# Setup
METAPOD_BK_GOPATH="${GOPATH}"
METAPOD_BK_GOARCH="${GOARCH}"
export GOPATH=`pwd`
export CGO_ENABLED=1

if [[ "${GOOS}" == "" ]]; then
    if [[ "${OSTYPE}" == "darwin" ]]; then
        export GOOS="darwin"
    else
        export GOOS="linux"
    fi
fi

mkdir -p bin/${GOOS}/x64
mkdir -p bin/${GOOS}/x32

# 64-bit
echo "Building 64-bit"
export GOARCH=amd64
pushd bin/${GOOS}/x64
go build -buildmode=c-shared metapod
popd

# 32-bit
echo "Building 32-bit"
export GOARCH=386
pushd bin/${GOOS}/x32
go build -buildmode=c-shared metapod
popd

# Cleanup
export GOPATH="${METAPOD_BK_GOPATH}"
export GOARCH="${METAPOD_BK_GOARCH}"
unset METAPOD_BK_GOPATH
unset METAPOD_BK_GOARCH
