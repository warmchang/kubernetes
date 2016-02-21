FROM debian:jessie

RUN DEBIAN_FRONTEND=noninteractive apt-get update -y \
    && DEBIAN_FRONTEND=noninteractive apt-get -yy -q \
    install \
    iptables \
    ca-certificates \
    file \
    util-linux \
    socat \
    curl \
    ethtool \
    && DEBIAN_FRONTEND=noninteractive apt-get autoremove -y \
    && DEBIAN_FRONTEND=noninteractive apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN curl -sfSL -o /usr/bin/share-mnt https://github.com/rancher/runc/releases/download/share-mnt-v0.0.4/share-mnt && \
    chmod u+x /usr/bin/share-mnt

RUN cp /usr/bin/nsenter /nsenter
ADD kubelet /usr/bin/

ADD kube-proxy /kube-proxy
ADD kubelet /kubelet
ADD entry.sh /usr/bin

ENTRYPOINT ["/usr/bin/entry.sh"]
