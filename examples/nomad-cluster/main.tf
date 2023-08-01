terraform {
  cloud {
    organization = "zolamk"
    workspaces {
      name = "terraform-provider-webdock"
    }
  }

  required_providers {
    webdock = {
      version = "0.1.1"
      source = "zolamk/webdock"
    }
    random = {
      source = "hashicorp/random"
      version = "3.4.3"
    }

    tls = {
      source = "hashicorp/tls"
      version = "4.0.4"
    }

    consul = {
      source = "hashicorp/consul"
      version = "2.18.0"
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