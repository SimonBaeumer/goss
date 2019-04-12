FROM ubuntu:18.04

RUN apt-get update
RUN apt-get install -y \
        bind9 \
        vim \
        dnsutils

COPY bind9 /etc/default/bind9
COPY named.conf.options /etc/bind/named.conf.options
COPY named.conf.local /etc/bind/named.conf.local

# Create zone file
RUN mkdir /etc/bind/zones
COPY db.dnstesting.com /etc/bind/zones/db.dnstesting.com

# Copy reverse zone config
COPY db.20.172 /etc/bind/zones/db.20.172

# Run syntax check
RUN named-checkconf
RUN named-checkzone dnstesting.com /etc/bind/zones/db.dnstesting.com
RUN named-checkzone 20.172.in-addr.arpa /etc/bind/zones/db.20.172

CMD ["/usr/sbin/named", "-g"]