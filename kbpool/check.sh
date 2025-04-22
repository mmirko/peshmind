#!/bin/bash

curl -s \
  -X POST http://localhost:3030/pengine/create \
  -d "format=json" \
  -d "ask=direct(X,Y)." \
  -d "chunk=10"
