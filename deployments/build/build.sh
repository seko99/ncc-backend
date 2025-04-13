#!/bin/bash

go test ./...

if [ $? -eq 0 ]
then
  echo Test success, building binary
  mkdir out
  go build -o out/ncc-backend
  if [ $? -eq 0 ]
  then
    echo Build success
  fi
else
  echo Test failed
fi