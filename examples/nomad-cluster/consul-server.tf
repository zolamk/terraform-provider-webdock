resource "random_uuid" "consul_bootstrap_token" {
}

resource "random_uuid" "consul_agent_token_secret_id" {
}

resource "random_string" "consul_gossip_encryption_key" {
  length = 32
  special = false
}

resource "webdock_server" "consul_server" {
  count = var.consul_server_instance_count
  name = "Consul Server ${count.index + 1}"
  image_slug = "webdock-ubuntu-jammy-cloud"
  profile_slug = "webdockbit-2022"
  location_id = "fi"
}

resource "random_string" "consul_server_user_password" {
  count = var.consul_server_instance_count
  length = 16
  special = false
  min_lower = 4
  min_upper = 4
  min_numeric = 8
}

resource "tls_private_key" "consul_ca_key" {
  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_self_signed_cert" "consul_ca" {
  private_key_pem = tls_private_key.consul_ca_key.private_key_pem
  is_ca_certificate = true
  
  subject {
    country = "ET"
    province = "Addis Ababa"
    locality = "Addis Ababa"
    common_name = "Consul Agent CA"
  }

  validity_period_hours = 43800

  allowed_uses = [
    "digital_signature",
    "any_extended",
    "cert_signing",
    "crl_signing",
    "timestamping",
    "ocsp_signing",
    "email_protection",
    "client_auth",
    "server_auth",
  ]
}

resource "tls_private_key" "consul_server_key" {
  algorithm = "RSA"
  rsa_bits = 2048
}

resource "tls_cert_request" "consul_server_cert_csr" {
  private_key_pem = tls_private_key.consul_server_key.private_key_pem

  dns_names = ["127.0.0.1", "localhost", "server.dc1.consul"]

  subject {
    common_name = "server.dc1.consul"
    country = "ET"
    province = "Addis Ababa"
    locality = "Addis Ababa"
  }
}

resource "tls_locally_signed_cert" "consul_server_cert" {
  cert_request_pem = tls_cert_request.consul_server_cert_csr.cert_request_pem
  ca_private_key_pem = tls_private_key.consul_ca_key.private_key_pem
  ca_cert_pem = tls_self_signed_cert.consul_ca.cert_pem

  validity_period_hours = 43800

  allowed_uses = [
    "digital_signature",
    "any_extended",
    "cert_signing",
    "crl_signing",
    "timestamping",
    "ocsp_signing",
    "email_protection",
    "client_auth",
    "server_auth",
    "key_encipherment",
  ]
}

resource "webdock_shell_user" "consul_server_user" {
  count = var.consul_server_instance_count
  username = "user"
  password = random_string.consul_server_user_password[count.index].result
  server_slug = webdock_server.consul_server[count.index].slug
  public_keys = [ data.webdock_public_keys.public_keys.public_keys[0].id ]

  connection {
    type = "ssh"
    user = "user"
    private_key = "${var.private_key}"
    host = webdock_server.consul_server[count.index].ipv4
  }

  provisioner "file" {
    content = file("./consul/node-policy.hcl")
    destination = "/tmp/node-policy.hcl"
  }

  provisioner "file" {
    content = tls_private_key.consul_ca_key.private_key_pem
    destination = "/tmp/consul-agent-ca-key.pem"
  }

  provisioner "file" {
    content = tls_self_signed_cert.consul_ca.cert_pem
    destination = "/tmp/consul-agent-ca.pem"
  }

  provisioner "file" {
    content = tls_private_key.consul_server_key.private_key_pem
    destination = "/tmp/dc1-server-consul-0-key.pem"
  }

  provisioner "file" {
    content = tls_locally_signed_cert.consul_server_cert.cert_pem
    destination = "/tmp/dc1-server-consul-0.pem"
  }

  provisioner "file" {
    source = "./consul/provision.sh"
    destination = "/tmp/provision.sh"
  }

  provisioner "file" {
    content = templatefile("./consul/server.hcl", {
      ip = webdock_server.consul_server[count.index].ipv4,
      first_server_ip = webdock_server.consul_server[0].ipv4,
      number_of_servers = var.consul_server_instance_count,
      gossip_encryption_key = random_string.consul_gossip_encryption_key.result
    })
    destination = "/tmp/consul.hcl"
  }

  provisioner "remote-exec" {
    inline = [
      "echo ${random_uuid.consul_bootstrap_token.result} > /tmp/root.token",
      "chmod +x /tmp/provision.sh",
      "echo ${random_string.consul_server_user_password[count.index].result} | sudo -k -S /tmp/provision.sh server ${webdock_server.consul_server[0].ipv4} ${webdock_server.consul_server[count.index].ipv4} ${random_uuid.consul_bootstrap_token.result} ${random_uuid.consul_agent_token_secret_id.result}",
      "rm /tmp/provision.sh"
    ]
  }
}
