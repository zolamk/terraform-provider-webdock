#/bin/bash

# install nomad
apt update

apt -y dist-upgrade

curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -

apt-add-repository -y "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"

apt update

apt-get -y install consul ufw

# enable firewall

# consul dns
ufw allow 8600

# consul http/s grpc/grpc tls
ufw allow 8500:8503/tcp

# consul serf
ufw allow 8301

ufw allow 8302

# consul rpc
ufw allow 8300

ufw allow ssh

ufw --force enable

# copy consul configuration
cp /tmp/consul.hcl /etc/consul.d/consul.hcl

# enable and start consul service

systemctl enable consul.service

mv /tmp/consul-agent-ca.pem /etc/consul.d/consul-agent-ca.pem

if [ "$1" = "server" ]; then
mv /tmp/consul-agent-ca-key.pem /etc/consul.d/consul-agent-ca-key.pem

mv /tmp/dc1-server-consul-0.pem /etc/consul.d/dc1-server-consul-0.pem

mv /tmp/dc1-server-consul-0-key.pem /etc/consul.d/dc1-server-consul-0-key.pem
fi

chown consul:consul -R /etc/consul.d

service consul start

# if this is the first server, bootstrap the acl
# and the tls ca and server certificates
if [ "$2" = "$3" ] && [ "$1" = "server" ]; then
# wait for conul to start
curl --retry 5 --retry-connrefused --retry-delay 5 http://127.0.0.1:8500/v1/status/leader

consul acl bootstrap /tmp/root.token

consul acl policy create -token="$4" -name node-policy -rules @/tmp/node-policy.hcl

consul acl token create -token="$4" -secret="$5" --description "consul agent token" -policy-name node-policy
fi
