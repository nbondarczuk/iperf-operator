# iperf3 in a container
#
# Run as Server:
# docker run  -it --rm --name=iperf3-srv -p 5201:5201 iperf3 -s
#
# Run as Client (first get server IP address):
# docker inspect --format "{{ .NetworkSettings.IPAddress }}" iperf3-srv
# docker run  -it --rm iperf3 -c <SERVER_IP>
#
FROM debian:bookworm-slim

# install binary and remove cache
RUN apt-get update \
    && apt-get install -y iperf3 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

COPY ./iperf3-*.sh /

# Expose the default iperf3 server port
EXPOSE 5201

# docker run -it <IMAGE> --help' is like running 'iperf3 --help'
ENTRYPOINT ["iperf3"]
