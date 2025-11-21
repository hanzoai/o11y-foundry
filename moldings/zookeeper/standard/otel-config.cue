package otel

inputs:{
    grpc: bool @tag(OTEL_GRPC)
    http: bool @tag(OTEL_HTTP)
    endpoint: string @tag(OTEL_ENDPOINT)
    tls: bool @tag(OTEL_TLS)
}

#protocols: {
  type: string "grpc | http"

}

#receiver: {
  name: string
  protocols: [#protocols]
  prometheus: [#PrometheusConfig]

}

#PrometheusConfig: {
  global: {
    scrape_interval: string
  }
  scrape_configs: [{
    job_name: string
    static_configs: [{
      targets: [string]
      labels: [string]: string
    }]
  }]
}


receivers: {
  otlp:
    protocols:{
        if inputs.grpc {
          grpc: {
            endpoint: 0.0.0.0:4317
            }
        }
        if inputs.http {
          http: {
            endpoint: 0.0.0.0:4318
            }
         }
      }
  prometheus:
    config:
      global:
        scrape_interval: 60s
      scrape_configs:
        - job_name: otel-collector
          static_configs:
          - targets:
              - localhost:8888
            labels:
              job_name: otel-collector
}


