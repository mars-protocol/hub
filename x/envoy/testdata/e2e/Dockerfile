FROM archlinux:base-devel

RUN pacman -Sy git go cargo deno yq vim --needed --noconfirm

# hermes

RUN cargo install ibc-relayer-cli --bin hermes --locked

# wasmd

WORKDIR /root

RUN git clone --depth 1 --branch v0.30.0 https://github.com/CosmWasm/wasmd

WORKDIR /root/wasmd

RUN make build

# marsd

WORKDIR /root

RUN git clone --depth 1 --branch ica2 https://github.com/mars-protocol/hub

WORKDIR /root/hub

RUN make build

# setup

RUN mv /root/.cargo/bin/hermes /bin/hermes
RUN mv /root/wasmd/build/wasmd /bin/wasmd
RUN mv /root/hub/build/marsd /bin/marsd

RUN tar -xzf /root/hub/x/envoy/testdata/configs.tar.gz -C /root

WORKDIR /root/hub/x/envoy/testdata/e2e
