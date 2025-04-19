# Configuration

> Note:
> 
> Naming: `qBittorrent_exporter` (QBE)

- Supports configuration with both file and env.
- Environmental variables take priority over config values.

## Flags

To view supported flags and their description
```bash
./qbittorrent_exporter -h
```
```bash
Version: v0.2.0-alpha

Usage: qbittorrent_exporter [ Options... ]

Available Options:
  -config string
    	Path to yaml config. (default "config.yaml")
  -ff-transient-state
    	[FeatureFlag][transient-state]
  -log-format string
    	Log format (default "default")
  -log-level string
    	Log level (default "info")
  -prefix string
    	Metrics prefix. (default "qb_")
```
Feature flags start with `ff` prefix
```
# Example
-ff-transient-state
```
Feature flags neither overlap with config file nor envs.

## Config file

Config file must be in YAML format.

```yaml
qBittorrent:
  baseUrl: http://127.0.0.1:8080
  username: admin
  password: adminpassword

metrics:
  port: 17171
  urlPath: /metrics

global:
  statePath: state.json
```

## Envs

| Name             | Example                |
| ---------------- | ---------------------- |
| QBE_URL          | https://127.0.0.1:8080 |
| QBE_USERNAME     | admin                  |
| QBE_PASSWORD     | adminpassword          |
| QBE_METRICS_PORT | 17171                  |
| QBE_METRICS_PATH | /metrics               |
| QBE_STATE_PATH   | state.json             |
**Table 1:** supported env and example values

## State

> If following metrics are not important to you, feel free to disable persistent state using ``

QBE creates `state.json` file to store:
- `dl_info_data_total` - Sum of all `dl_info_data` sessions recorded by QBE
- `up_info_data_total` - Sum of all `up_info_data` sessions recorded by QBE

## Log

> Log levels are not case-sensitive.

Possible values are:
- debug
- info
- warn
- error