# qrgen

qrgen is a QR generator service written in Go.

## Env variables:

The service requires the following environment variables:

- APP_PORT (service port)
- GRPC_SERVER_PORT (grpc server port)
- GRPC_OTEL_COLLECTOR (otel collector receiver grpc port)

## Run it

To start all the services run:

```bash
docker compose -f docker-compose.yaml -f docker-compose.prod.yaml up -d
```

And stop them all with:

```bash
docker compose down
```

## Endpoints
- GET / (main page)
- POST /qr (generate QR)
- GET /metrics (prometheus metrics)

## Instrumentation

The service is instrumented using the [OpenTelemtry Go SDK](https://opentelemetry.io/docs/instrumentation/go/).

### Tracing

Traces are implemented using the Go SDK. When produced they're sent to the open telemetry collector instance through a gRPC endpoint using OTLP (open telemetry protocol), batch processed, exported and loaded to a Jaeger instance (port `14250`). The Jaeger UI runs at http://localhost:16686.

### Metrics

Metrics (_just the simple counter example for now_) are exposed through a HTTP endpoint at http://localhost/metrics and consumed by the Prometheus instance.

Metrics could be instrumented manually with the open telemetry SDK and sent to the otel collector to be exported to a external prometheus instance, but [the Go SDK does not support this yet](https://opentelemetry.io/docs/instrumentation/go/manual/#creating-metrics).

## Logs

Logs are generated using the default gorilla/mux logging handler, but they're not collected in any way yet.