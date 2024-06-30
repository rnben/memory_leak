# README

## 生成私钥和证书

- openssl 方式

```bash
openssl genrsa -out key.pem 2048
openssl req -new -x509 -key key.pem -out cert.pem -days 3650
```

- `crypto/tls`

```bash
go run $GOROOT/src/crypto/tls/generate_cert.go --host localhost
```

## 启动服务

```bash
go run main.go
```

## 测试 socks5

```bash
curl --socks5 foo:bar@127.0.0.1:1080 http://10.16.238.2:8314/ping
``