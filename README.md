[![License][License-Image]][License-URL] [![ReportCard][ReportCard-Image]][ReportCard-URL] [![Build][Build-Status-Image]][Build-Status-URL]
# ws-tcp-proxy
Simple websocket tcp proxy.

```
Usage:
  ws-tcp-proxy <address:port> [flags]

Flags:
  -a, --auto-cert string         register hostname with LetsEncrypt
  -c, --cert string              path to cert.pem for TLS
  -k, --key string               path to key.pem for TLS
  -p, --port int                 server port (default 8080)
      --tcp-tls                  connect to TCP address via TLS
      --tcp-tls-cert string      path to client.crt for TCP TLS
      --tcp-tls-key string       path to client.key for TCP TLS
      --tcp-tls-root-ca string   path to ca.crt for TCP TLS
  -t, --text-mode                text mode
  -v, --version                  display version

```

## Install

```
go get -v github.com/zquestz/ws-tcp-proxy
cd $GOPATH/src/github.com/zquestz/ws-tcp-proxy
make
make install
```

If you have issues building ws-tcp-proxy, you can vendor the dependencies by using [gvt](https://github.com/FiloSottile/gvt):

```
go get -u github.com/FiloSottile/gvt
cd $GOPATH/src/github.com/zquestz/ws-tcp-proxy
gvt restore
```

## Contributors

* [Josh Ellithorpe (zquestz)](https://github.com/zquestz/)

## License

ws-tcp-proxy is released under the MIT license.

[License-URL]: http://opensource.org/licenses/MIT
[License-Image]: https://img.shields.io/npm/l/express.svg
[ReportCard-URL]: http://goreportcard.com/report/zquestz/ws-tcp-proxy
[ReportCard-Image]: https://goreportcard.com/badge/github.com/zquestz/ws-tcp-proxy
[Build-Status-URL]: http://travis-ci.com/zquestz/ws-tcp-proxy
[Build-Status-Image]: https://travis-ci.com/zquestz/ws-tcp-proxy.svg?branch=master
