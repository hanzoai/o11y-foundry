package example

directory:       string @tag(a,var=cwd)
operatingSystem: string @tag(b,var=os)
cpuArchitecture: string @tag(c,var=arch)
currentUsername: string @tag(d,var=username)
currentHostname: string @tag(e,var=hostname)
randomnessA:     int    @tag(f,var=rand,type=int)
randomnessB:     int    @tag(g,var=rand,type=int)
currentTimeA:    string @tag(h,var=now)
currentTimeB:    string @tag(i,var=now)

processors: {
  batch:
    send_batch_size: 10000
    send_batch_max_size: 11000
    timeout: 10s
  resourcedetection:
    detectors: [env, system]
    timeout: 2s
  signozspanmetrics/delta:
    metrics_exporter: signozclickhousemetrics
    metrics_flush_interval: 60s
    latency_histogram_buckets: [100us, 1ms, 2ms, 6ms, 10ms, 50ms, 100ms, 250ms, 500ms, 1000ms, 1400ms, 2000ms, 5s, 10s, 20s, 40s, 60s ]
    dimensions_cache_size: 100000
    aggregation_temporality: aggregation_temporality_delta
    enable_exp_histogram: true
    dimensions:
      - name: service.namespace
        default: default
      - name: deployment.environment
        default: default
      - name: signoz.collector.id
      - name: service.version
      - name: browser.platform
      - name: browser.mobile
      - name: k8s.cluster.name
      - name: k8s.node.name
      - name: k8s.namespace.name
      - name: host.name
      - name: host.type
      - name: container.name
}

extensions: {
  health_check:
    endpoint: 0.0.0.0:13133
  pprof:
    endpoint: 0.0.0.0:1777
}

exporters:{ 
  otlp/signoz:
    endpoint: inputs.endpoint
    tls:
      if inputs.tls {
        insecure: false
      } else {
        insecure: true
      }
  service:
    telemetry:
      logs:
        encoding: json
    extensions:
      - health_check
      - pprof
    pipelines:
      traces:
        receivers: [otlp]
        processors: [signozspanmetrics/delta, batch]
        exporters: [otlp/signoz]
      metrics:
        receivers: [otlp]
        processors: [batch]
        exporters: [otlp/signoz]
      metrics/prometheus:
        receivers: [prometheus]
        processors: [batch]
        exporters: [otlp/signoz]
      logs:
        receivers: [otlp]
        processors: [batch]
        exporters: [otlp/signoz]
}
