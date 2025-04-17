# qBittorrent Exporter

### Quickstart

```bash
curl --proto '=https' --tlsv1.3 -sSfL  https://raw.githubusercontent.com/AlexKhomych/qbittorrent_exporter/refs/tags/v0.1.0-alpha/deploy/quickstart.sh | bash
```

1) Edit `config.yaml` under `$HOME/.config/qbittorrent_exporter/config.yaml` path.
2) And restart `systemctl restart qbittorrent_exporter.service --user`.

### Configuration

- Supports configuration with both file and env.
- Environmental variables take priority over config values.

Example ENVs:
```
QBT_USERNAME="admin"
QBT_PASSWORD="adminpassword"
QBT_BASE_URL="http://127.0.0.1:8080"

METRICS_PORT="17171"
METRICS_URL_PATH="/metrics"
```

#### Optional(Enable Lingering)
Quickstart installs exporter as systemd service under user's control.
By default users needs to have active session, such that can be achieved with `screen`, `tmux`.
To automatically run user's services after (re)boot, make use of `lingering`.

```bash
sudo loginctl enable-linger $(whoami)

systemctl enable qbittorrent_exporter.service --now --user
```
[Read more here](https://manpages.debian.org/bullseye/systemd/loginctl.1.en.html)

#### Optional(Grafana Dashboard)

- Dashboard ID: `23261`
- Dashboard Link: [Dashboard](https://grafana.com/grafana/dashboards/23261-qbittorrent-dashboard/?tab=reviews)

### Debug

Try lowering `-log-level` flag in `$HOME/.config/systemd/user/qbittorrent_exporter.service` to `debug`

#### Viewing logs
```bash
journalctl --user -u qbittorrent_exporter.service
```

### Credits

Inspired by [@caseyscarborough](https://github.com/caseyscarborough/qbittorrent-exporter)
