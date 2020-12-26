FROM golang:1.15-alpine3.12 as backend

WORKDIR /go/src/app
COPY . .

RUN go build -o /dist/notification-server cmd/notification-server/notification-server.go


FROM node:lts-alpine3.12 as frontend

WORKDIR /src
COPY ui .

RUN npm install && npm run build


FROM alpine:3.12.1
COPY --from=backend /dist/notification-server /app/notification-server
COPY --from=frontend /src/dist /ui

RUN chmod +x /app/notification-server

ENTRYPOINT []
CMD ["/app/notification-server", "-host", "0.0.0.0", "-static", "/ui"]
