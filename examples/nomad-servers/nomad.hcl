datacenter = "dc1"
data_dir  = "/opt/nomad/data"
bind_addr = "${ip}"

addresses {
  http = "${ip} 127.0.0.1"
}

server {
  enabled          = true
  bootstrap_expect = ${number_of_servers}
  encrypt = "${gossip_encryption_key}"
  server_join {
    retry_join = ["${first_nomad_server_ip}"]
  }
}

telemetry {
  collection_interval = "15s"
  prometheus_metrics = true
  publish_allocation_metrics = true
  publish_node_metrics = true
}

acl {
  enabled = true
}