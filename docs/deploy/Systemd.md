# Systemd deployment

### Quickstart

```bash
curl --proto '=https' --tlsv1.3 -sSfL  https://raw.githubusercontent.com/AlexKhomych/qbittorrent_exporter/refs/tags/1.0.2/deploy/quickstart.sh | bash
```

1) Edit `config.yaml` under `$HOME/.config/qbittorrent_exporter/config.yaml` path.
2) And restart `systemctl restart qbittorrent_exporter.service --user`.

## Configuration

### Optional(Enable Lingering)
Quickstart installs exporter as systemd service under user's control.
By default users needs to have active session, such that can be achieved with `screen`, `tmux`.
To automatically run user's services after (re)boot, make use of `lingering`.

```bash
sudo loginctl enable-linger $(whoami)

systemctl enable qbittorrent_exporter.service --now --user
```
[Read more here](https://manpages.debian.org/bullseye/systemd/loginctl.1.en.html)

## Debug

### Adjust log level

Edit `$HOME/.config/qbittorrent_exporter/env`

```service
QBE_LOG_LEVEL=debug
```

### Viewing logs
```bash
journalctl --user -u qbittorrent_exporter.service
```

