FROM golang:1.17.6 as builder

WORKDIR /go/src/github.com/gavinmcnair/prometheuscheck

COPY . .

RUN make prometheuscheck

FROM scratch
MAINTAINER Gavin McNair

ARG git_repository="Unknown"
ARG git_commit="Unknown"
ARG git_branch="Unknown"
ARG built_on="Unknown"

LABEL git.repository=$git_repository
LABEL git.commit=$git_commit
LABEL git.branch=$git_branch
LABEL build.on=$built_on

COPY --from=builder /go/src/github.com/gavinmcnair/prometheuscheck/bin/linux/prometheuscheck .

CMD [ "/prometheuscheck" ]
