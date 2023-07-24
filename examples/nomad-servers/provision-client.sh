#/bin/bash

# install nomad
apt update

apt -y dist-upgrade

curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -

apt-add-repository -y "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"

apt update

apt-get -y install nomad ufw

# enable firewall

# nomad http/rpc
ufw allow 4646:4647/tcp

# nomad serf
ufw allow 4648

ufw allow ssh

ufw --force enable

# copy nomad configuration

cp /tmp/nomad.hcl /etc/nomad.d/nomad.hcl

# enable and start nomad service

systemctl enable nomad.service

service nomad start