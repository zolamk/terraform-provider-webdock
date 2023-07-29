---
global:
  scrape_interval:     15s
  evaluation_interval: 15s

{{ with nomadVar "prometheus" }}
remote_write:
  - url: {{.url}}
    basic_auth:
      username: {{.username}}
      password: {{.password}}
{{ end}}

scrape_configs:
  - job_name: 'nomad_metrics'
    metrics_path: /v1/metrics
    params:
      format: [prometheus]
