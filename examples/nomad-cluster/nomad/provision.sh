#/bin/bash

# install nomad
apt update

apt -y dist-upgrade

curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -

apt-add-repository -y "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"

apt update

apt-get -y install nomad ufw curl consul

# enable firewall

# nomad http/rpc
ufw allow 4646:4647/tcp

# nomad serf
ufw allow 4648

ufw allow ssh

# consul dns
ufw allow 8600

# consul http/s grpc/grpc tls
ufw allow 8500:8503/tcp

# consul serf
ufw allow 8301

ufw allow 8302

# consul rpc
ufw allow 8300

ufw --force enable

systemctl enable consul.service

# copy consul configuration
mv /tmp/consul.hcl /etc/consul.d/consul.hcl

mv /tmp/consul-agent-ca.pem /etc/consul.d/consul-agent-ca.pem

sudo chown consul:consul -R /etc/consul.d

service consul start

# copy nomad configuration
mv /tmp/nomad.hcl /etc/nomad.d/nomad.hcl

if [ "$1" = "server" ]; then
# run nomad with non privileged user
sed -i 's/\[Service\]/\[Service\]\nUser=nomad\nGroup=nomad/' /lib/systemd/system/nomad.service
fi

# enable and start nomad service

systemctl enable nomad.service

service nomad start

# if this is the first server, bootstrap the acl
if [ "$2" = "$3" ] && [ "$1" = "server" ]; then
# wait for nomad to start
curl --retry 5 --retry-connrefused --retry-delay 5 http://127.0.0.1:4646/v1/status/leader

nomad acl bootstrap /tmp/root.token
fi

if [ "$1" = "client" ]; then
apt-get -y install podman unzip

curl https://releases.hashicorp.com/nomad-driver-podman/0.5.0/nomad-driver-podman_0.5.0_linux_amd64.zip --output /tmp/nomad-driver-podman_0.5.0_linux_amd64.zip

mkdir -p /opt/nomad/data/plugins

unzip /tmp/nomad-driver-podman_0.5.0_linux_amd64.zip -d /opt/nomad/data/plugins/

service nomad restart
fi