# Postgres Check

A simple statically linked go binary for validating Postgres connectivity.

In its default mode it will connect to a database and validate connectivity using a `ping`

The check can also be configured with a `go` style duration (1s, 2m, 4h, etc) which will allow the check to ran continuously.

The pre-built version of the container can be found at

`docker.io/gavinmcnair/postgrescheck:1.1.15`

| Environment Variable  | Default | Required | Description |
|---|---|---|---|
| DATABASE  |  - | No | Validates presence of database  |
| DB_HOST  |  - | Yes  | Hostname of target database or file location containing hostname  |
| DB_PORT  | 5432  |  - |  Port number of listening database |
| DB_USER  |  - |  Yes |  Username for target database or file location containing username |
| DB_PASS  |  - |  Yes |  Password for target database or file location containing password |
| LISTEN_PORT | 8080 | - | Listen Port for /health and /metrics endpoint |
| REPEAT_INTERVAL  |  0 | -  | Repeat check on go duration interval (1s,5m,10h,1d,etc)  |
| SSLMODE | verify-ca | -| Valid Options are [ 'disable']  |
