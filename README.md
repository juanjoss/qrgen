# qrgen

qrgen is a QR generator service written in Go.

env variables required to run the service
- APP_PORT (service port)
- GRPC_SERVER_PORT (grpc server port)

endpoints
- GET / (main page)
- POST /qr (generate QR)

build and run the image with docker
> `docker build -t qrgen:latest .`

> `docker run -dp <port>:<port> --env APP_PORT= --env GRPC_SERVER_PORT= qrgen`