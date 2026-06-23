# Omada Exporter

Prometheus Exporter written in Go for TP-Link Omada SDN.

This exporter is built with native Go HTTP and [prometheus/client_golang](https://github.com/prometheus/client_golang)
libraries to expose metrics from Omada SDN to Prometheus.
Metrics are queried live from the controller when the `/metrics` endpoint is accessed,
avoiding unnecessary polling when not in use.

Supports both Omada OpenAPI and Web API to accommodate current platform limitations.

## Grafana Dashboards

### [Site Overview](/dashboards/site_overview.json)

![image](/pictures/site_overview.png)

### [Router](/dashboards/router.json)

![image](/pictures/router_1.png)

### [Switch](/dashboards/switch.json)

![image](/pictures/switch_1.png)

### [Access Point](/dashboards/access_point.json)

![image](/pictures/access_point_1.png)

## 🚀 Getting Started

### Prerequisites

- **Omada Controller** with configured:
  - [**OpenAPI Client**](#omada-authentication-setup)
  - [**Service user account**](#omada-authentication-setup)
- **Prometheus** (or another metrics scraper) to collect and store the metrics
- **Grafana** to visualize the data using built-in dashboards

### Setup

1. Clone the Repository

```shell
git clone https://github.com/filipegalo/omada-exporter-go.git
cd omada-exporter-go
```

2. Create the `.env` file from the template and fill in your values

```shell
cp .env.example .env
```

   Then edit `.env` — see [Configuration Parameters](#configuration-parameters) for each variable.

3. Build the Docker image

```shell
docker build -t omada_exporter .
```

4. Run the container

```shell
docker run -d \
  --env-file .env \
  -p 8080:8080 \
  --name omada_exporter \
  omada_exporter
```

> [!TIP]
> A `Makefile` is included with shortcuts: `make docker-build`, `make docker-run`,
> or `make run` to run locally (loads `.env` automatically). Run `make help` for all targets.

5. Configure Prometheus to scrape the data:

```YAML
scrape_configs:
  - job_name: "Omada"
    static_configs:
      - targets: ["<docker_host_ip>:8080"]
        labels:
          device_name: "<omada_controller_friendly_name>"
    metrics_path: /metrics
```

6. Import dashboards from [`dashboards`](/dashboards/) directory to your Grafana instance.

> [!TIP]
> After starting up the exporter navigate to `http://"<docker_host_ip>:8080/metrics`
> and check `omada_http_client_requests_total` metrics.
> There should be **no** metrics with code different different than 200

### Omada Authentication Setup

- OpenAPI Client – Created via: `Settings -> Platform Integration`.
  Assign admin role for full API access.
- Service User – Create under: `Account section` at `Global level`.
  Assign viewer role for read-only access.

### Configuration Parameters

All configuration values are read from environment variables. These can be provided via a .env file, environment variables in your system, or directly in Docker Compose.

| Variable              | Description                                                            | Default |
| --------------------- | ---------------------------------------------------------------------- | ------- |
| `LOG_LEVEL`           | Logging verbosity level (_debug_, _info_, _warn_, _error_)             | `error` |
| `OMADA_URL`           | Full base URL to the Omada controller (e.g., https://192.168.1.1:8043) | -       |
| `OMADA_SITE_NAME`     | Name of the site you wish to monitor                                   | -       |
| `OMADA_CLIENT_ID`     | OpenAPI Client ID created in Omada Platform Integration                | -       |
| `OMADA_CLIENT_SECRET` | OpenAPI Client Secret                                                  | -       |
| `OMADA_USERNAME`      | Username for Web API access                                            | -       |
| `OMADA_PASSWORD`      | Password for Web API access                                            | -       |

## Tested on devices:

- **Switch**: [SG2218 v1.20](https://www.tp-link.com/en/business-networking/omada-switch-smart/sg2218/)
- **Router**: [ER707-M2 v1.0](https://www.tp-link.com/en/business-networking/omada-router-wired-router/er707-m2/v1/)
- **Access Point**: [EAP650(EU) v1.0](https://www.tp-link.com/en/business-networking/omada-wifi-ceiling-mount/eap650/v1/)
- **Omada controller**: [Software Controller v5.15.8.s hosted in Docker container](https://github.com/mbentley/docker-omada-controller)
