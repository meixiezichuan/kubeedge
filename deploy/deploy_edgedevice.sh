#!/bin/bash

#cloudip="$1"
KUBEEDGE_ROOT=$(unset CDPATH && cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)
cd $KUBEEDGE_ROOT
#CGO_ENABLED=1  go build -v -o _output/local/bin -ldflags="${GO_LDFLAGS} -w -s -extldflags -static" ./edgedevice/cmd/edgedevice

_output/local/bin/edgedevice --defaultconfig > edgedevice.yaml
#sed -i "s/\(httpServer:\s*\).*/\1${cloudip}:10002/" edgedevice.yaml

#mkdir -p /etc/kubeedge/config
#sudo mv edgedevice.yaml > /etc/kubeedge/config
#sudo cp _output/local/bin/edgedevice /usr/local/bin
#sudo systemctl restart edgedevice
#systemctl status edgedevice
