# Prometheus AMQP proxy

Proxy for feeding prometheus scraper from AMQP RPC sources.

## Building and running

### Build and run using [bazel](http://baze.io/)

    bazel run //proxy:proxy -- -amqp_url=amqp://my_rabbitmq_url/ -amqp_exchange=my_prometheus_exchange

Visiting [http://localhost:8200/proxy?target=prometheus](http://localhost:8200/proxy?target=prometheus) will send an AMQP RPC
request to the "prometheus" queue in my_prometheus_exchange exchange, and return the returned data.

## Prometheus Configuration

The AMQP proxy needs to be passed the AMQP queue as a parameter, this can be done with relabelling.

Example config:
```
scrape_configs:
  - job_name: 'amqp_proxy'
    metrics_path: '/proxy'
    target_groups:
      - targets: ['queue1', 'queue2']
    relabel_configs:
      - source_labels: [__address__]
        regex: (.*)
        target_label: __param_target
        replacement: ${1}
      - source_labels: []
        regex: .*
        target_label: __address__
        replacement: localhost:8200  # AMQP proxy

```
