resource "webdock_server" "nomad_client" {
  count = var.nomad_client_instance_count
  name = "Nomad Client ${count.index + 1}"
  image_slug = "webdock-ubuntu-jammy-cloud"
  profile_slug = "webdockbit-2022"
  location_id = "fi"
}

resource "random_string" "nomad_client_user_password" {
  count = var.nomad_client_instance_count
  length = 16
  special = false
  min_lower = 4
  min_upper = 4
  min_numeric = 8
}

resource "webdock_shell_user" "nomad_client_user" {
  count = var.nomad_client_instance_count
  username = "user"
  password = random_string.nomad_client_user_password[count.index].result
  server_slug = webdock_server.nomad_client[count.index].slug
  public_keys = [ data.webdock_public_keys.public_keys.public_keys[0].id ]

  connection {
    type = "ssh"
    user = "user"
    private_key = "${var.private_key}"
    host = webdock_server.nomad_client[count.index].ipv4
  }

  provisioner "file" {
    source = "./nomad/provision.sh"
    destination = "/tmp/provision.sh"
  }

  provisioner "file" {
    content = tls_self_signed_cert.consul_ca.cert_pem
    destination = "/tmp/consul-agent-ca.pem"
  }

  provisioner "file" {
    content = templatefile("./consul/client.hcl", {
      ip = webdock_server.nomad_client[count.index].ipv4,
      first_server_ip = webdock_server.consul_server[0].ipv4,
      gossip_encryption_key = random_string.consul_gossip_encryption_key.result
      consul_agent_token = random_uuid.consul_agent_token_secret_id.result
    })
    destination = "/tmp/consul.hcl"
  }

  provisioner "file" {
    content = templatefile("./nomad/client.hcl", {
      ip = webdock_server.nomad_client[count.index].ipv4,
      first_nomad_server_ip = webdock_server.nomad_server[0].ipv4,
    })
    destination = "/tmp/nomad.hcl"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/provision.sh",
      "echo ${random_string.nomad_client_user_password[count.index].result} | sudo -k -S /tmp/provision.sh client ${webdock_server.nomad_server[0].ipv4} ${webdock_server.nomad_client[count.index].ipv4}",
      "rm /tmp/provision.sh"
    ]
  }
}
