FROM golang:1.22.3-alpine3.20 as tailscale

RUN mkdir /app

WORKDIR /app

ENV TSFILE=tailscale_1.66.1_amd64.tgz

RUN wget https://pkgs.tailscale.com/stable/${TSFILE} && \
  tar xzf ${TSFILE} --strip-components=1 && \
  rm ${TSFILE}

# Download Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

# Build go program
RUN CGO_ENABLED=0 GOOS=linux go build -o /proxy-pass

from alpine:3.20
COPY --from=tailscale /proxy-pass /proxy-pass
COPY --from=tailscale /app/tailscaled /tailscaled
COPY --from=tailscale /app/tailscale /tailscale
RUN mkdir -p /var/run/tailscale /var/cache/tailscale /var/lib/tailscale

# Create a new user with UID 10014
RUN addgroup -g 10014 choreo && \
    adduser  --disabled-password --uid 10014 --ingroup choreo choreouser

RUN chown -R 10014 /var/run/tailscale /var/cache/tailscale /var/lib/tailscale

EXPOSE 1055 8080

RUN mkdir /home/wso2

WORKDIR /home/wso2

COPY start.sh .

RUN chmod +x start.sh

USER 10014

CMD ["/home/wso2/start.sh"]
