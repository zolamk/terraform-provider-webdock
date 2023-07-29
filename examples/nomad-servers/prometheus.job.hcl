job "prometheus" {
  type        = "service"

  variable "prometheus" {
    type = object({
      url = string
      username = string
      password = string
    })
  }

  group "monitoring" {
    count = 1

    network {
      port "prometheus_ui" {
        static = 9090
      }
    }

    restart {
      attempts = 2
      interval = "30m"
      delay    = "15s"
      mode     = "fail"
    }

    task "prometheus" {
      template {
        change_mode = "noop"
        source = "prometheus.yml.tpl"
        destination = "local/prometheus.yml"
      }

      driver = "podman"

      config {
        image = "prom/prometheus:latest"

        volumes = [
          "local/prometheus.yml:/etc/prometheus/prometheus.yml",
        ]

        ports = ["prometheus_ui"]
      }

      service {
        name = "prometheus"
        tags = ["urlprefix-/"]
        port = "prometheus_ui"

        check {
          name     = "prometheus_ui port alive"
          type     = "http"
          path     = "/-/healthy"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}
