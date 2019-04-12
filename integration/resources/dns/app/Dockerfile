FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y \
        dnsutils \
        iproute2 \
        iputils-ping \
        vim

COPY resolv.conf /etc/resolv.conf

CMD ["/bin/bash"]