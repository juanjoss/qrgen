# qrgen

qrgen is a QR generator service written in Go.

Env variables required to run the service:
- APP_PORT (service port)
- GRPC_SERVER_PORT (grpc server port)

Endpoints
- GET / (main page)
- POST /qr (generate QR)

Pull and run the image with docker:

```bash
docker pull jujoss/qrgen:latest
```

```bash
docker run -dp port:port \ 
    --env APP_PORT=port \ 
    --env GRPC_SERVER_PORT=port \ 
    jujoss/qrgen:latest
```