set -o errexit
set -o nounset
set -o pipefail


KUBEEDGE_ROOT=$(unset CDPATH && cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd -P)
echo "KUBEEDGE_ROOT: $KUBEEDGE_ROOT"
cd $KUBEEDGE_ROOT
make all WHAT=cloudcore BUILD_WITH_CONTAINER=false
sudo cp _output/local/bin/cloudcore /usr/local/bin
sudo cp build/tools/cloudcore.service /etc/systemd/system

_output/local/bin/cloudcore --defaultconfig > cloudcore.yaml
sed -i 's/\(kubeConfig\s*:\s*\)"[^"]*"/\1"\/etc\/rancher\/k3s\/k3s.yaml"/g' cloudcore.yaml

# get ip from cloudcore.yaml
IP=$(awk '/advertiseAddress:/{getline; print $2}' cloudcore.yaml)
echo "ip: $IP"

# generate cert
sudo ./build/tools/certgen.sh genCertAndKey server $IP

# copy config file
sudo mkdir -p /etc/kubeedge/config
sudo mv cloudcore.yaml /etc/kubeedge/config

sudo systemctl start cloudcore
sudo systemctl enable cloudcore
systemctl status cloudcore
