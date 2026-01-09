# Reverse Proxy

This document covers common reverse-proxy setups when placing a proxy in front of qBittorrent and notes how `insecureSkipVerify` and certificates affect the exporter.

> Note:
>
> "Reverse proxy" here may refer to a proxy in front of the qBittorrent Web UI or a proxy terminating TLS in front of the exporter. Configure `qBittorrent.baseUrl` (or `QBE_URL`) to match how the exporter reaches the Web UI.

## Typical setups

- TLS termination at the reverse proxy: the proxy serves HTTPS to clients and forwards requests to qBittorrent over HTTP on localhost. In this case, the exporter should use `http://` in `qBittorrent.baseUrl` when connecting directly to the backend, or `https://` if the exporter connects to the proxy endpoint.
- TLS passthrough or direct TLS to qBittorrent: the exporter connects over HTTPS directly to qBittorrent and must validate the certificate presented by qBittorrent (or the proxy if it terminates TLS).

## Nginx example (TLS termination)

```nginx
server {
    listen 443 ssl;
    server_name qb.example.com;

    ssl_certificate /etc/ssl/certs/example.crt;
    ssl_certificate_key /etc/ssl/private/example.key;

    location / {
        proxy_pass http://127.0.0.1:8080/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

If you expose the Web UI via a path prefix (for example `/qb/`), ensure `qBittorrent.baseUrl` uses the same path so the exporter can access the correct endpoints.

## Headers

Ensure the proxy forwards `Host` and `X-Forwarded-Proto` when rewriting or terminating TLS. Some qBittorrent Web UI features rely on `Host` or the original scheme.

## Certificates and `insecureSkipVerify`

- Trusted certificates: If qBittorrent (or the proxy the exporter connects to) presents a certificate signed by a CA trusted by the exporter's host OS, no additional configuration is required. Keep `insecureSkipVerify: false` (recommended).
- Self-signed or internal CA certificates: you have two secure options:
  - Add the issuing CA certificate to the host's trust store so the exporter (Go runtime) will validate the chain normally.
  - Use a certificate issued by a private/internal CA and make that CA trusted on the exporter host.

### Using `insecureSkipVerify`

If you cannot add a CA to the trust store, you can disable certificate verification. In the YAML config:

```yaml
qBittorrent:
  baseUrl: https://qb.example.local:8080
  insecureSkipVerify: true
```

Or via environment variable:

```bash
QBE_INSECURE_SKIP_VERIFY=true
```

**Important notes when setting `insecureSkipVerify: true`:**

- It disables certificate chain and hostname verification for TLS connections made by the exporter. This defeats the primary protections of TLS and makes connections vulnerable to man-in-the-middle attacks.
- Only use `insecureSkipVerify` in trusted, isolated networks (e.g., local testing, temporary debugging) or when you fully understand the security implications.
- If you use `insecureSkipVerify: true` and the server uses a self-signed certificate, you may encounter errors such as `SID cookie not found`. This can happen if the exporter fails to authenticate due to certificate or login issues.
    - **Note:** The `SID cookie not found` error can also indicate that the qBittorrent API has blacklisted the exporter's IP address after multiple failed login attempts. See [Issue #3](https://github.com/alexkhomych/qbittorrent_exporter/issues/3) for more details.

## Practical recommendations

- Prefer issuing certificates from a trusted CA or add your internal CA to the exporter host's trust store instead of using `insecureSkipVerify`.
- If you must use self-signed certs for development, set `insecureSkipVerify` only on non-production systems and document the reason.
- When the reverse proxy terminates TLS and the exporter accesses the backend over plain HTTP, certificate concerns apply only between clients and the proxy; ensure internal traffic policies are acceptable.

## Quick config examples

- Connect to a proxy-terminated TLS endpoint (trusted cert):

```yaml
qBittorrent:
  baseUrl: https://qb.example.com
  insecureSkipVerify: false
```

- Connect to a self-signed TLS endpoint (temporary):

```yaml
qBittorrent:
  baseUrl: https://qb.local:8080
  insecureSkipVerify: true
```
