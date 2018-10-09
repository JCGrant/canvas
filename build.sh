#!/bin/bash

if ! [ -x "$(command -v statik)" ]; then
  go get -u github.com/rakyll/statik
fi

statik -src client/build
go build .
