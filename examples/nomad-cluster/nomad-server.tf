resource "random_string" "gossip_encryption_key" {
  length = 32
  special = false
}

resource "random_uuid" "nomad_bootstrap_token" {
}

resource "webdock_server" "nomad_server" {
  count = var.nomad_server_instance_count
  name = "Nomad Server ${count.index + 1}"
  image_slug = "webdock-ubuntu-jammy-cloud"
  profile_slug = "webdockbit-2022"
  location_id = "fi"
}

resource "random_string" "nomad_server_user_password" {
  count = var.nomad_server_instance_count
  length = 16
  special = false
  min_lower = 4
  min_upper = 4
  min_numeric = 8
}

resource "webdock_shell_user" "nomad_server_user" {
  count = var.nomad_server_instance_count
  username = "user"
  password = random_string.nomad_server_user_password[count.index].result
  server_slug = webdock_server.nomad_server[count.index].slug
  public_keys = [ data.webdock_public_keys.public_keys.public_keys[0].id ]

  connection {
    type = "ssh"
    user = "user"
    private_key = "${var.private_key}"
    host = webdock_server.nomad_server[count.index].ipv4
  }

  provisioner "file" {
    content = tls_self_signed_cert.consul_ca.cert_pem
    destination = "/tmp/consul-agent-ca.pem"
  }

  provisioner "file" {
    source = "./nomad/provision.sh"
    destination = "/tmp/provision.sh"
  }

  provisioner "file" {
    content = templatefile("./consul/client.hcl", {
      ip = webdock_server.nomad_server[count.index].ipv4,
      first_server_ip = webdock_server.consul_server[0].ipv4,
      gossip_encryption_key = random_string.consul_gossip_encryption_key.result
      consul_agent_token = random_uuid.consul_agent_token_secret_id.result
    })
    destination = "/tmp/consul.hcl"
  }

  provisioner "file" {
    content = templatefile("./nomad/server.hcl", {
      ip = webdock_server.nomad_server[count.index].ipv4,
      first_nomad_server_ip = webdock_server.nomad_server[0].ipv4,
      number_of_servers = var.nomad_server_instance_count,
      gossip_encryption_key = random_string.gossip_encryption_key.result
    })
    destination = "/tmp/nomad.hcl"
  }

  provisioner "remote-exec" {
    inline = [
      "echo ${random_uuid.nomad_bootstrap_token.result} > /tmp/root.token",
      "chmod +x /tmp/provision.sh",
      "echo ${random_string.nomad_server_user_password[count.index].result} | sudo -k -S /tmp/provision.sh server ${webdock_server.nomad_server[0].ipv4} ${webdock_server.nomad_server[count.index].ipv4}",
      "rm /tmp/root.token /tmp/provision.sh"
    ]
  }
}
