# Omada Exporter

[![Release](https://github.com/filipegalo/omada-exporter-go/actions/workflows/release.yml/badge.svg)](https://github.com/filipegalo/omada-exporter-go/actions/workflows/release.yml)
[![Latest release](https://img.shields.io/github/v/release/filipegalo/omada-exporter-go?sort=semver)](https://github.com/filipegalo/omada-exporter-go/releases/latest)
[![Go version](https://img.shields.io/github/go-mod/go-version/filipegalo/omada-exporter-go)](go.mod)

Prometheus Exporter written in Go for TP-Link Omada SDN.

This exporter is built with native Go HTTP and [prometheus/client_golang](https://github.com/prometheus/client_golang)
libraries to expose metrics from Omada SDN to Prometheus.
Metrics are queried live from the controller when the `/metrics` endpoint is accessed,
avoiding unnecessary polling when not in use.

Supports both Omada OpenAPI and Web API to accommodate current platform limitations.

## 📊 Grafana Dashboards

Ready-to-use dashboards live in [`dashboards/`](/dashboards/). To import: in Grafana go to
**Dashboards → New → Import**, upload the JSON file, and select your Prometheus data source.

| Dashboard | Description | JSON |
| --------- | ----------- | ---- |
| **Site Overview** | High-level health and totals across the whole site | [site_overview.json](/dashboards/site_overview.json) |
| **Router** | Gateway throughput, WAN status, and port metrics | [router.json](/dashboards/router.json) |
| **Switch** | Per-port status, PoE, link speed, and traffic | [switch.json](/dashboards/switch.json) |
| **Access Point** | Radios, clients, and wireless uplink metrics | [access_point.json](/dashboards/access_point.json) |

<details>
<summary><b>Site Overview</b> — preview</summary>

![Site Overview](/pictures/site_overview.png)

</details>

<details>
<summary><b>Router</b> — preview</summary>

![Router 1](/pictures/router_1.png)
![Router 2](/pictures/router_2.png)
![Router 3](/pictures/router_3.png)

</details>

<details>
<summary><b>Switch</b> — preview</summary>

![Switch 1](/pictures/switch_1.png)
![Switch 2](/pictures/switch_2.png)
![Switch 3](/pictures/switch_3.png)

</details>

<details>
<summary><b>Access Point</b> — preview</summary>

![Access Point 1](/pictures/access_point_1.png)
![Access Point 2](/pictures/access_point_2.png)
![Access Point 3](/pictures/access_point_3.png)
![Access Point 4](/pictures/access_point_4.png)

</details>

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
> After starting the exporter, navigate to `http://<docker_host_ip>:8080/metrics`
> and check the `omada_http_client_requests_total` metric.
> There should be **no** entries with a status code other than 200.

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
| `OMADA_USERNAME`      | Username for Web API access                                            | -          |
| `OMADA_PASSWORD`      | Password for Web API access                                            | -          |
| `METRICS_PATH`        | HTTP path the metrics are served on                                    | `/metrics` |
| `METRICS_PORT`        | Port the exporter listens on                                           | `8080`     |

## Tested Devices

- **Switch**: [SG2218 v1.20](https://www.tp-link.com/en/business-networking/omada-switch-smart/sg2218/)
- **Router**: [ER707-M2 v1.0](https://www.tp-link.com/en/business-networking/omada-router-wired-router/er707-m2/v1/)
- **Access Point**: [EAP650(EU) v1.0](https://www.tp-link.com/en/business-networking/omada-wifi-ceiling-mount/eap650/v1/)
- **Omada controller**: [Software Controller v5.15.8.s hosted in Docker container](https://github.com/mbentley/docker-omada-controller)
