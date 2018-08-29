#!/bin/bash

VERSION=${VERSION:-1.10.5}

rm -f go.*
go mod init
echo >> go.mod
echo 'require (' >> go.mod
for i in staging/src/k8s.io/*; do
    echo k8s.io/$(basename $i) kubernetes-${VERSION} >> go.mod
done
echo ')' >> go.mod
