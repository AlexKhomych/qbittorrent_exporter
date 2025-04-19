# Docker deployment

### Quickstart

```bash
docker run -dit --name qbittorrent_exporter \
    --publish 17171:17171 \
    -e QBE_URL='http://[your_qb_url]:[qb_port]' \
    alexkhomychw/qbittorrent_exporter:v0.2.0-alpha
```

Metrics will be available on `http://[docker_host_ip]:17171/metrics` by default.

## Configuration

Please read [Configuration](../Configuration) for more info.

> Only docker specific(different) options will be explained here.

### Persist state

State file is created under `/app/config/state.json` by default. Can be changed using either config file or envs.

[Read more](../Configuration#state)

### Run with custom config

1. Create config file
```bash
mkdir .config
touch .config/config.yaml
```
2. Edit config file [Read more](../Configuration#config-file)
3. Run docker container with mount flag
```bash
docker run -dit --name qbittorrent_exporter -p 17171:17171 -v ./config:/app/config alexkhomychw/qbittorrent_exporter:v0.2.0-alpha
```

### Run using envs

```bash
docker run -dit --name qbittorrent_exporter \
    --publish 17171:17171 \
    -e QBE_URL='http://[your_qb_url]:[qb_port]' \
    -e QBE_USERNAME='admin' \
    -e QBE_PASSWORD='adminpassword' \
    alexkhomychw/qbittorrent_exporter:v0.2.0-alpha
```

[Available Envs](../Configuration#envs)

## Debug

Check container logs
```bash
docker container logs qbittorrent_exporter
```

Adjust log level
```
# Example is omitted as you may run it in docker-compose, k8s and other environment.
# Therefore, just add/change -log-level flag to debug value
```