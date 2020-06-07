# kube-event-logger
Logging Kubernetes Events for Observabulity


Running kube-event-logger inside a cluster that you're allows you to export kubernetes events as stdout logs (allowing you to use fluent-bit/fluentd) as part of your normal kubernetes log handling.

# Build

`docker build -t jamesglennan/kube-event-logger .` 