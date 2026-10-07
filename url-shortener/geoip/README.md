# Country lookup

The backend downloads DB-IP Country Lite MMDB at each container startup
and on the first of every month at 00:00 Moscow time (UTC+3).
Source: https://db-ip.com/db/download/ip-to-country-lite
The writable directory persists the database across container restarts.
The database is not included in the image or Git.

Data: DB-IP, licensed under CC BY 4.0. Attribution is displayed in Analytics.
Downloads are size-limited and verified with the MMDB reader before atomic
file replacement and a locked reader swap. No backend restart is needed.
Failures keep the old database and retry hourly (including delayed releases).
Until the first successful download, country lookup returns unknown.

`GEOIP_DB_PATH` selects the file. `TRUSTED_PROXY_HOST` identifies the reverse
proxy by DNS (Compose: frontend); only this peer may supply `X-Real-IP`.
Nginx overwrites that header with its connection's remote address. If another
proxy/CDN is added, configure nginx real_ip with that proxy's trusted ranges.
Never trust client-supplied forwarded headers directly.

Docker Desktop/localhost traffic may expose only a private gateway address;
such traffic has no country. Old hardcoded RU records are not rewritten.
