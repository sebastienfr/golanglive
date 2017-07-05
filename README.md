# golanglive
golang live coding session

## Key and certificate
### Generate private key (.key)

```sh
### Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048

### Key considerations for algorithm "ECDSA" ≥ secp384r1
#### List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key
```

### Generation of self-signed(x509) public key (PEM-encodings `.pem`|`.crt`) based on the private (`.key`)

```sh
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

## Slides

Use the Golang `present` tool in present folder for offline or see [link](http://go-talks.appspot.com/github.com/sebastienfr/golanglive/present/golanglive.slide#1) for the online one
   
## Run
   
```go
go run live.go
```

Runs an application with 2 web servers : *HTTPS* on port `8080`, *HTTP* on port `8081`
Go to [HTTPS/2 Go live](https://localhost:8081/index.html) and check performances against
[HTTP/1 Go live version](https://localhost:8081/index.html)
   
