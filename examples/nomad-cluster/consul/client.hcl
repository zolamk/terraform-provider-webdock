datacenter = "dc1"
data_dir  = "/opt/consul/data"
bind_addr = "${ip}"
retry_join = ["${first_server_ip}"]
encrypt = "${gossip_encryption_key}"

performance {
    raft_multiplier = 1
}

acl {
    enabled = true
    default_policy = "deny"
    enable_token_persistence = true
    tokens = {
      default = "${consul_agent_token}"
    }
}

tls {
   defaults {
      ca_file = "/etc/consul.d/consul-agent-ca.pem"

      verify_incoming = true
      verify_outgoing = true
   }
   internal_rpc {
      verify_server_hostname = true
   }
}

auto_encrypt {
  tls = true
}
