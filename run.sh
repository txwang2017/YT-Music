#!/bin/bash
id=$1
filename=$2
go build -o run . && ./run $id $filename