#!/bin/bash

set -u

HOME="${HOME:-"~"}"

mkdir -p $HOME/.local/bin $HOME/.config/qbittorrent_exporter $HOME/.local/share/qbittorrent_exporter $HOME/.config/systemd/user

curl -L https://github.com/AlexKhomych/qbittorrent_exporter/releases/download/1.0.0/qbittorrent_exporter -o $HOME/.local/bin/qbittorrent_exporter
chmod +x $HOME/.local/bin/qbittorrent_exporter

cat > $HOME/.config/qbittorrent_exporter/config.yaml << EOF
qBittorrent:
  baseUrl: http://127.0.0.1:8080
  username: admin
  password: adminpassword

metrics:
  port: 17171
  urlPath: /metrics

global:
  statePath: $HOME/.local/share/qbittorrent_exporter/state.json
EOF

cat > $HOME/.config/qbittorrent_exporter/env << EOF
QBE_LOG_LEVEL=info
EOF

cat > $HOME/.config/systemd/user/qbittorrent_exporter.service << EOF
[Unit]
Description=qBittorrent exporter
Wants=network-online.target
After=network-online.target

[Service]
EnvironmentFile=$HOME/.config/qbittorrent_exporter/env
ExecStart=$HOME/.local/bin/qbittorrent_exporter -log-format default -log-level \$QBE_LOG_LEVEL -config $HOME/.config/qbittorrent_exporter/config.yaml
ExecReload=/bin/kill -s HUP \$MAINPID
WorkingDirectory=$HOME/.local/share/qbittorrent_exporter
Restart=always
RestartSec=5
TimeoutSec=0

[Install]
WantedBy=default.target
EOF

systemctl daemon-reload --user
systemctl status qbittorrent_exporter.service --user
