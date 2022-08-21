# qrgen

qrgen is a QR generator service written in Go.

env variables required to run the service
- APP_PORT (service port)
- GRPC_SERVER_PORT (grpc server port)

endpoints
- GET / (main page)
- POST /qr (generate QR)

pull and run the image with docker

`docker pull jujoss/qrgen:v1`

`docker run -dp <port>:<port> --env APP_PORT=<port> --env GRPC_SERVER_PORT=<port> jujoss/qrgen:<tag>`