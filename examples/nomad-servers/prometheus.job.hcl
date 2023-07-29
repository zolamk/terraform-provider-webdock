job "prometheus" {
  type        = "service"

  group "monitoring" {
    count = 1

    network {
      port "prometheus_ui" {
        static = 9090
      }
    }

    task "prometheus" {
      template {
        change_mode = "noop"
        data = <<EOF
---
global:
  scrape_interval:     15s
  evaluation_interval: 15s

remote_write:
  - url: ${ prometheus_url }
    basic_auth:
      username: ${ prometheus_username }
      password: ${ prometheus_password }

scrape_configs:
  - job_name: 'nomad_metrics'
    metrics_path: /v1/metrics
    params:
      format: [prometheus]

EOF
        destination = "local/prometheus.yml"
      }

      driver = "docker"

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
        provider = "nomad"

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
