resource "webdock_server" "test" {
  name         = "acc-test-server"
  location_id  = "fi"
  profile_slug = "webdock-standard-2vcpu-2gb"
  image_slug   = "ubuntu-2404-noble"
}
