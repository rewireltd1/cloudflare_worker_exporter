# Cloudflare Worker Exporter for Prometheus

This is a simple server that scrapes Cloudflare graphql analytics API for worker stats and exports them via HTTP for
Prometheus consumption.

## Getting Started

### Prerequisites
In order to run the exporter you need to prepare the following items:
1. Cloudlfare API Bearer token with `Account.Account Analytics` permissions
2. Account id of your account

### running the code

```bash
./cloudflare_worker_exporter [flags]
```

Help on flags:

```bash
./cloudflare_worker_exporter --help
```

### Dot env file
The implementation supports a .env file. You can override the path using the DOTENV_FILE environment variable

```bash
DOTENV_FILE=/env/variables ./cloudflare_worker_exporter [flags]
```

### Docker

* Build
```bash
docker build -t cloudflare_worker_exporter:local -f Dockerfile .
```

* Run
```bash
docker run -it\
 -e CLOUDFLARE_ANALYTICS_TOKEN=[token]\
 -e CLOUDFLARE_ACCOUNT_ID=[account id]\
 -p 9184:9184\
  cloudflare_worker_exporter:local
```

## Metrics
The exporter exposes the following worker metrics:
-----------------------------------------------------------------------------------------------------------------------------------
| Name                                        | Description                                                             | Type    |
|:--------------------------------------------|:------------------------------------------------------------------------|:-------:|
| `cloudflare_worker_requests_up`             | Did the last scrape request to get requests stats finished successfully | gauge   |
| `cloudflare_worker_cpu_time_up`             | Did the last scrape request to get cpu usage finished successfully      | gauge   |
| `cloudflare_worker_cpu_time_percentile`     | worker cpu time per worker, status and percentile                       | gauge   |
| `cloudflare_worker_requests_received_total` | total number of requests per worker and status                          | counter |
| `cloudflare_worker_errors_total`            | total number of errors per worker and status                            | counter |
| `cloudflare_worker_subrequests_total`       | total number of subrequests per worker and status                       | counter |
-----------------------------------------------------------------------------------------------------------------------------------

**The labels for the metrics are:**
* worker - the worker script name
* state - worker final status - success, clientDisconnected, scriptThrewException
* percentile - for the cpu time time series, possible values `25,50,75,90,99,999`


## License

Copyright 2021 Rewire (O.S.G) Research and Development Ltd. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License"), see [LICENSE](https://github.com/rewireltd1/cloudflare_worker_exporter/blob/master/LICENSE).

For more information, see [the coming blog post](https://coming-soon).
