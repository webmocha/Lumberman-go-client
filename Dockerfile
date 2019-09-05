FROM golang as build

WORKDIR /app
ADD . /app
RUN cd /app && \
  make build-linux64

FROM alpine

COPY --from=build /app/bin/lmc /

ENTRYPOINT ["/lmc"]
