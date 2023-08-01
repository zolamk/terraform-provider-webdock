datacenter = "dc1"
data_dir  = "/opt/consul/data"
bind_addr = "${ip}"
bootstrap_expect = ${number_of_servers}
retry_join = ["${first_server_ip}"]
encrypt = "${gossip_encryption_key}"
server = true

ui_config {
    enabled = true
}

performance {
    raft_multiplier = 1
}

acl {
    enabled = true
    default_policy = "deny"
    enable_token_persistence = true
}

tls {
   defaults {
      ca_file = "/etc/consul.d/consul-agent-ca.pem"
      cert_file = "/etc/consul.d/dc1-server-consul-0.pem"
      key_file = "/etc/consul.d/dc1-server-consul-0-key.pem"

      verify_incoming = true
      verify_outgoing = true
   }
   internal_rpc {
      verify_server_hostname = true
   }
}

auto_encrypt {
  allow_tls = true
}
