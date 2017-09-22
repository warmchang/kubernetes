#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

KUBE_ROOT=$(dirname "${BASH_SOURCE}")/..
cd $KUBE_ROOT

setup_links()
{
    for i in staging/src/k8s.io/*; do
        ln -s ../../staging/src/k8s.io/$(basename $i) vendor/k8s.io/$(basename $i)
    done
}

cat > trash.conf << EOF
# package
k8s.io/kubernetes

########################################################################################
#### Don't edit this file, it is generated, edit trash.in and run hack/trash-gen.sh ####
########################################################################################

EOF

cat trash.in >> trash.conf

cat ./Godeps/Godeps.json | jq -r '(.Deps | .[] | "\(.ImportPath) \(.Rev)\n")' | sed '/^$/d' >> trash.conf

sed -i -e '/zz_generated.openapi.go/d' .gitignore

trash -k
setup_links
./hack/update-codegen.sh

trash
setup_links
