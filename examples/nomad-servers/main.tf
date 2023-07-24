terraform {
  required_providers {
    webdock = {
      version = "0.1.1"
      source = "zolamk/webdock"
    }
    random = {
      source = "hashicorp/random"
      version = "3.4.3"
    }
  }
}

provider "webdock" {
  token = "${var.token}"
}

data "webdock_images" "images" {}

data "webdock_locations" "locations" {}

data "webdock_profiles" "profiles" {
  location_id = "${data.webdock_locations.locations.locations[0].id}"
}

data "webdock_public_keys" "public_keys" {}

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
  profile_slug = "webdockbitmore-2022"
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

output "nomad_server_user_passwords" {
  value = random_string.nomad_server_user_password[*].result
  description = "nomad server user password"
  sensitive = true
}

output "nomad_server_bootstrap_token" {
  value = random_uuid.nomad_bootstrap_token.result
  description = "the nomad server bootstrap token"
  sensitive = true
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
    source = "./provision.sh"
    destination = "/tmp/provision.sh"
  }

  provisioner "file" {
    content = templatefile("./nomad.hcl", {
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
      "rm /tmp/root.token /tmp/provision.sh /tmp/nomad.hcl"
    ]
  }
}

# Nomad Clients

resource "webdock_server" "nomad_client" {
  count = var.nomad_client_instance_count
  name = "Nomad Client ${count.index + 1}"
  image_slug = "webdock-ubuntu-jammy-cloud"
  profile_slug = "webdockbitmore-2022"
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

output "nomad_client_user_passwords" {
  value = random_string.nomad_server_user_password[*].result
  description = "nomad client user password"
  sensitive = true
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
    source = "./provision.sh"
    destination = "/tmp/provision.sh"
  }

  provisioner "file" {
    content = templatefile("./nomad-client.hcl", {
      ip = webdock_server.nomad_server[count.index].ipv4,
      first_nomad_server_ip = webdock_server.nomad_server[0].ipv4,
      number_of_servers = var.nomad_server_instance_count,
      gossip_encryption_key = random_string.gossip_encryption_key.result
    })
    destination = "/tmp/nomad.hcl"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/provision.sh",
      "echo ${random_string.nomad_client_user_password[count.index].result} | sudo -k -S /tmp/provision.sh client",
      "rm /tmp/provision.sh /tmp/nomad.hcl"
    ]
  }
}