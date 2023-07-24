#/bin/bash

# install nomad
apt update

apt -y dist-upgrade

curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -

apt-add-repository -y "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"

apt update

apt-get -y install nomad ufw curl

# enable firewall

# nomad http/rpc
ufw allow 4646:4647/tcp

# nomad serf
ufw allow 4648

ufw allow ssh

ufw --force enable

# copy nomad configuration
cp /tmp/nomad.hcl /etc/nomad.d/nomad.hcl

if [ "$1" = "server" ]; then
# run nomad with non privileged user
sed -i 's/\[Service\]/\[Service\]\nUser=nomad\nGroup=nomad/' /lib/systemd/system/nomad.service
fi

# enable and start nomad service

systemctl enable nomad.service

service nomad start

# if this is the first server, bootstrap the acl
if [ "$3" = "$3" ] && [ "$1" = "server" ]; then
# wait for nomad to start
curl --retry 5 --retry-connrefused --retry-delay 5 http://127.0.0.1:4646/v1/status/leader

nomad acl bootstrap /tmp/root.token
fi
