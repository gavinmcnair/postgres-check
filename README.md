# Postgres Check

A simple statically linked go binary for validating metrics functionality. Postgres connectivity.

The check is configured with a `go` style duration (1s, 2m, 4h, etc) which will allow the check to ran continuously.

The pre-built version of the container can be found at

`docker.io/gavinmcnair/prometheuscheck:1.0.1`

| Environment Variable  | Default | Required | Description |
|---|---|---|---|
| LISTEN_PORT | 8080 | - | Listen Port for /health and /metrics endpoint |
| REPEAT_INTERVAL  |  0 | -  | Repeat check on go duration interval (1s,5m,10h,1d,etc)  |
