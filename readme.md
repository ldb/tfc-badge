# tfc-badge
CI status badges for [Terraform Cloud](https://www.terraform.io/cloud).

## Overview

The system is quite simple:
- On each run, Terraform Cloud informs this server ("_tfcbd_") using a [webhook notification](https://www.terraform.io/docs/cloud/workspaces/notifications.html)
- _tfcbd_ stores the result in memory
- Upon request, _tfcbd_ renders a status badge that can be embedded into READMEs
- For each run, _tfcbd_ can publish annotations to Grafana

## Design decisions

This project is intentionally small and simple and comes with a couple of tradeoffs:
- **No (necessary) persistence.** If the server is restarted, records of previous runs are lost. 
  _tfcbd_ renders a special "unknown" badge in that case. *Note: You can persist state across restarts by providing a 
  file path to the `--file` flag now.*
    
- **No scalability.** Running multiple instances of _tfcbd_ may yield wrong results.
  Honestly though, it should handle even big loads just fine and probably does not need
  to be highly available either.

## Usage

### Publishing Grafana Annotations

On each run, _tfcbd_ can publish annotations to Grafana. These annotations are spans for the entire duration of a run.
They contain the following tags:

- `tfc-badge`: fixed to allow selecting the right annotations easily
- `status`: contains the final status of the run, e.g "applied"
- `workspace`: the name of the workspace of the run

In order to enable the creation of annotation, the `GRAFANA_HOST` and `GRAFANA_API_KEY` environment variables need to be passed (see below).

### Configuration Flags

You can use the following flags to control the behaviour of `tfcbd`:

```shell
  -addr string
        specify the address and port to listen on (default ":3030")
  -debug
        enable debug log output
  -file string
        path to file to store persistent cache in
  -metrics string
        specify the address and port to listen on for metrics (default ":9080")
```

### Environment variables

- `GRAFANA_HOST` can be set to pass the Grafana host (example: `https://grafana.monitoring.local`)
- `GRAFANA_API_KEY` can be set to pass an API key for Grafana
