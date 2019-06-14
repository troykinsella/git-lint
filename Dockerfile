FROM golang:latest

ENV GO111MODULE on

WORKDIR $GOPATH/src/github.com/troykinsella/git-lint
COPY . .

RUN set -eux; \
    apt-get update; \
    apt-get install -y make; \
    make dist; \
    mv git-lint_linux_amd64 /usr/local/bin/git-lint; \
    rm -rf $GOPATH/*; \
    apt-get remove -y make; \
    apt-get clean all; \
    rm -rf /var/lib/apt/lists/*;

CMD /usr/local/bin/git-lint
