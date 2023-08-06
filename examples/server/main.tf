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

resource "webdock_server" "server" {
  count = var.server_instance_count
  name = "Server ${count.index + 1}"
  image_slug = "webdock-ubuntu-jammy-cloud"
  profile_slug = "webdockbit-2022"
  location_id = "fi"
}

resource "random_string" "server_user_password" {
  count = var.server_instance_count
  length = 16
  special = false
  min_lower = 4
  min_upper = 4
  min_numeric = 8
}

resource "webdock_shell_user" "nomad_server_user" {
  count = var.server_instance_count
  username = "user"
  password = random_string.server_user_password[count.index].result
  server_slug = webdock_server.server[count.index].slug
  public_keys = [ data.webdock_public_keys.public_keys.public_keys[0].id ]
}
