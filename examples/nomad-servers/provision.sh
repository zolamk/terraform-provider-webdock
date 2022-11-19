#/bin/bash
apt update

apt -y dist-upgrade

curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -

apt-add-repository -y "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"

apt update

apt-get -y install nomad ufw

# nomad http/rpc
ufw allow 4646:4647/tcp

# nomad serf
ufw allow 4648

ufw --force enable

cp /tmp/nomad.hcl /etc/nomad.d/nomad.hcl

# run nomad with non privileged user
sed -i 's/\[Service\]/\[Service\]\nUser=nomad\nGroup=nomad/' /lib/systemd/system/nomad.service

systemctl enable nomad.service

service nomad start
