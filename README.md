# alfred-grafana-dashboards-workflow

An Alfred workflow to search your [Grafana](https://grafana.com) dashboards and select one to open in your default browser.

## Installation

Clone this repo and run:

```
make install
```

## Configuration

During installation, set the following required variables:

* `GRAFANA_HOST` - required (e.g. `https://example.com/`)

If your Grafana instance is protected by HTTP basic auth, set the following variables:

* `GRAFANA_BASIC_AUTH_USER` - optional
* `GRAFANA_BASIC_AUTH_PASSWORD` - optional
