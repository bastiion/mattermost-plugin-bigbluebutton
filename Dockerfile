#FROM mattermost-build-base:ubuntu-golang-1.13-node-14.12
FROM mattermost-build-base:latest

RUN apk add --no-cache expat expat-dev

#COPY /home/basti/daten/Projekte/bigbluebutton/workspace/go-src-cache /go/src/

RUN mkdir -p /go/src/github.com/blindsidenetworks/mattermost-plugin-bigbluebutton

WORKDIR /go/src/github.com/blindsidenetworks/mattermost-plugin-bigbluebutton

COPY . /go/src/github.com/blindsidenetworks/mattermost-plugin-bigbluebutton/

RUN make build
