resource "webdock_public_key" "test" {
  name = "acc-test-key"
  key  = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBb2DvOSJwMmWVQSPm3LfxpdbVG3o5rjRv0a7m/test+key terraform-acc-test"
}

resource "webdock_server" "test" {
  name         = "acc-test-server"
  location_id  = "fi"
  profile_slug = "webdock-standard-2vcpu-2gb"
  image_slug   = "ubuntu-2404-noble"
}

resource "webdock_shell_user" "test" {
  server_slug = webdock_server.test.slug
  username    = "acc-test-user"
  password    = "S3cur3P@ssw0rd!"
  group       = "webdock"
  shell       = "/bin/bash"
  public_keys = [webdock_public_key.test.id]
}
