terraform {
  required_providers {
    webdock = {
      version = "0.1.0"
      source = "github.com/zolamk/webdock"
    }
    random = {
      source = "hashicorp/random"
      version = "3.4.3"
    }
    nomad = {
      source = "hashicorp/nomad"
      version = "1.4.19"
    }
  }
}

provider "webdock" {
  token = "${var.token}"
}

data "webdock_account" "account" {}

data "webdock_images" "images" {}

data "webdock_locations" "locations" {}

data "webdock_profiles" "profiles" {
  location_id = "${data.webdock_locations.locations.locations[0].id}"
}

data "webdock_public_keys" "public_keys" {}

resource "random_string" "nomad_server_user_password" {
  length = 16
  special = false
  min_lower = 4
  min_upper = 4
  min_numeric = 8
}

resource "webdock_server" "nomad_server" {
  count = var.nomad_server_instance_count
  name = "Nomad Server ${count.index + 1}"
  image_slug = "webdock-ubuntu-jammy-cloud"
  profile_slug = "webdockbit-2022"
  location_id = "fi"
}

resource "webdock_shell_user" "nomad_server_user" {
  count = var.nomad_server_instance_count
  username = "user"
  password = random_string.nomad_server_user_password.result
  server_slug = webdock_server.nomad_server[count.index].slug
  public_keys = [ data.webdock_public_keys.public_keys.public_keys[0].id ]

  connection {
    type = "ssh"
    user = "user"
    password = random_string.nomad_server_user_password.result
    host = webdock_server.nomad_server[count.index].ipv4
  }

  provisioner "file" {
    source = "./provision.sh"
    destination = "/tmp/provision.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/provision.sh",
      "echo ${random_string.nomad_server_user_password.result} | sudo -k -S /tmp/provision.sh ${var.nomad_server_instance_count} ${webdock_server.nomad_server[0].ipv4} ${webdock_server.nomad_server[count.index].ipv4}"
    ]
  }
}