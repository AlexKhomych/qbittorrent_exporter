# Metrics

> Note: `qb_` is a prefix and can be different if you changed it.

| Name                           | Description                                         |
| ------------------------------ | --------------------------------------------------- |
| # TorrentsInfo                 |                                                     |
| qb_torrent_name                | contains torrent's name as a label                  |
| qb_torrent_state               | contains torrent's state as a label                 |
| qb_torrent_progress            | [0.0 to 1.0] float value                            |
| qb_torrent_dlspeed             | float value in bytes(SI)                            |
| qb_torrent_upspeed             | float value in bytes(SI)                            |
| qb_torrent_download            | float value in bytes(SI)                            |
| qb_torrent_amount_left         | float value in bytes(SI)                            |
| qb_torrent_ration              | float value                                         |
| qb_torrent_eta                 | float value in seconds                              |
| qb_torrent_num_seeds           | float value                                         |
| qb_torrent_num_leechs          | float value                                         |
| # TransferInfo                 |                                                     |
| qb_transfer_status             | qBittorrent's connectivity status as label          |
| qb_transfer_dl_info_speed      | qBittorrent's global download speed in bytes(SI)    |
| qb_transfer_dl_info_data       | qBittorrent's current session download in bytes(SI) |
| qb_transfer_up_info_speed      | qBittorrent's global upload speed in bytes(SI)      |
| qb_transfer_up_info_data       | qBittorrent's current session upload in bytes(SI)   |
| qb_transfer_dl_rate_limit      | qBittorrent's download rate limit                   |
| qb_transfer_up_rate_limit      | qBittorrent's upload rate limit                     |
| qb_transfer_dht_nodes          | qBittorrent's number of dht nodes                   |
| qb_transfer_dl_info_data_total | qBittorrent's total download in bytes(SI)           |
| qb_transfer_up_info_data_total | qBittorrent's total upload in bytes(SI)             |
| # Version                      |                                                     |
| qb_app_version                 | qBittorrent's version as a label                    |

**Table 1:** exported metrics


