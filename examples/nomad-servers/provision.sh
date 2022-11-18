#/bin/bash
apt update

apt -y dist-upgrade

curl -fsSL https://apt.releases.hashicorp.com/gpg | apt-key add -

apt-add-repository -y "deb [arch=amd64] https://apt.releases.hashicorp.com $(lsb_release -cs) main"

apt update

apt-get -y install nomad

mkdir /data

echo "
datacenter = \"dc1\"
data_dir  = \"/opt/nomad/data\"
bind_addr = \"$3\"

server {
  enabled          = true
  bootstrap_expect = $1
  server_join {
    retry_join = [\"$2\"]
  }
}

client {
  enabled = true
  host_volume \"data\" {
    path = \"/data/\"
    read_only = false
  }
  server_join {
    retry_join = [\"$3\"]
  }
}

telemetry {
  publish_allocation_metrics = true
  publish_node_metrics = true
}
" > /tmp/nomad.hcl

cp /tmp/nomad.hcl /etc/nomad.d/nomad.hcl

systemctl enable nomad.service

service nomad start

# only bootstrap acl on the first nomad server
# if [ "$2" = "$3" ]; then
# nomad acl bootstrap
# fi