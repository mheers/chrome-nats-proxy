# Chrome-NATS-Proxy

Tunnels HTTP(S) requests through a local NATS server through a Chrome extension.

```
+-----------------+    +-----------------+    +-----------------+    +-----------------+    +-----------------+
|  public         |    |  this           |    |  local          |    |  this           |    |                 |
|  HTTP(S) Server |<-->|  Chrome         |<-->|  NATS           |<-->|  HTTP(S) Proxy  |<-->|  HTTP(S) Client |
|                 |    |  Extension      |    |  Server         |    |                 |    |                 |
|                 |    |                 |    |                 |    |                 |    |                 |
+-----------------+    +-----------------+    +-----------------+    +-----------------+    +-----------------+
```

## Background

In cases where a HTTP(S) server is not directly reachable from the client, a NATS server can be used to tunnel the requests. This is useful in cases where the client is behind a firewall or a NAT.

The Chrome extension provides a NATS Request/Reply interface to the client and publishes this service via Websockets to the local NATS server.

This HTTP(S) proxy translates incoming HTTP(S) requests to NATS requests and sends them to the NATS server that the extension processes by sending the request to the HTTP(S) server and returning the response to the client.

## Features

- [x] HTTP(S) Proxy
- [x] [Clipboard-Sync](https://github.com/mheers/clipboard-sync) forward support

## Usage

### Install the Chrome extension

1. Clone the repository
1. Build the extension in the `chrome-extension` directory (`npm run build`)
1. Open Chrome
1. Go to `chrome://extensions/`
1. Enable `Developer mode`
1. Click on `Load unpacked`
1. Select the `chrome-extension/build` directory

### Create seeds for the docker-compose file

```bash
make create-seeds
```

Add the resulting lines to the `.env`

### Create a credentials file

```bash
make create-credentials
```

```bash
echo '{"url": "https://marcelheers.de", "method": "GET", "headers": { "Accept": "text/html" }, "body": null }' | nats --creds mq.creds --server ws://localhost:9222 request instance.chrome-nats-proxy.http
```

```bash
HTTP_PROXY=http://127.0.0.1:8080 HTTPS_PROXY=http://127.0.0.1:8080 curl -k https://www.marcelheers.de/
```

# TODO
- [x] implement a client that works with setting the HTTP_PROXY environment variable
- [x] fix cors issue in extension (extension needs a contentScript.js)
- [ ] allow server nats to be hosted on a path
- [ ] allow configuration without .env (just setting `--nats-uri` and `--creds-file` flags)
- [ ] make clipboard-sync optional
- [ ] integrate nats-server
- [ ] integrate clipboard-sync
- [ ] increase timeout to 30sec
- [ ] increase payload size to 10MB
