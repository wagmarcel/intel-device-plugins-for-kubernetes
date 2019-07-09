FROM golang:1.11 as builder
ARG DIR=/go/src/github.com/intel/intel-device-plugins-for-kubernetes
WORKDIR $DIR
COPY . .
RUN cd cmd/iot_plugin; go install

FROM gcr.io/distroless/base
COPY --from=builder /go/bin/iot_plugin /usr/bin/iot_device_plugin
COPY --from=builder /go/src/github.com/intel/intel-device-plugins-for-kubernetes/cmd/iot_plugin/deviceConfigExample.json /etc/deviceconfig.json
CMD ["/usr/bin/iot_device_plugin"]
