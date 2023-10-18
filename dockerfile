# syntax=docker/dockerfile:1

FROM golang:1.19 as builder

LABEL stage=gobuilder

ENV CGO_ENABLED 0
ENV GOOS linux
ENV GOPROXY https://goproxy.cn,direct

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
#COPY *.go ./

COPY . .

# Build
RUN go build -o dot



# 使用alpine这个轻量级镜像为基础镜像--运行阶段
FROM alpine AS runner
WORKDIR  /app

COPY --from=builder /app /app


# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8091

# Run
CMD ["./dot"]

