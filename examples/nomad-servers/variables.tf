variable "token" {
  type = string
  description = "Webdock API Token"
}

variable "nomad_server_instance_count" {
  type = number
  default = 3
  description = "The number of nomad servers to deploy"
}