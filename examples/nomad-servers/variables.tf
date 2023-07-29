variable "token" {
  type = string
  description = "Webdock API Token"
}

variable "nomad_server_instance_count" {
  type = number
  description = "The number of nomad servers to deploy"
}

variable "private_key" {
  type = string
  description = "the private key to use for ssh connections. used by the provisioner to install and setup the nomad server"
}

variable "nomad_client_instance_count" {
  type = number
  description = "The number of nomad clients to deploy"
}

variable "prometheus" {
  type = object({
    url = "string"
    username = "string"
    password = "string"
  })
  description = "The prometheus server configuration"
}