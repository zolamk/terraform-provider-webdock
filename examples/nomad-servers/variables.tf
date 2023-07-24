variable "token" {
  type = string
  description = "Webdock API Token"
}

variable "nomad_server_instance_count" {
  type = number
  default = 3
  description = "The number of nomad servers to deploy"
}

variable "private_key" {
  type = string
  description = "the private key to use for ssh connections. used by the provisioner to install and setup the nomad server"
}