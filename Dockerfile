#########################################
# Build stage
#########################################
FROM node:lts-alpine AS ui-builder

# Setup env
RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH
ENV REACT_APP_API_ROOT /api
ENV REACT_APP_REDIRECT_URL /login

# Install dependencies
COPY ui/package.json /usr/src/app/package.json
COPY ui/package-lock.json /usr/src/app/package-lock.json
RUN npm install --silent

# Build website
COPY ./ui /usr/src/app
RUN npm run build

#########################################
# Build stage
#########################################
FROM golang:1.18 AS builder

# Repository location
ARG REPOSITORY=github.com/ncarlier

# Artifact name
ARG ARTIFACT=readflow

# Copy sources into the container
ADD . /go/src/$REPOSITORY/$ARTIFACT

# Set working directory
WORKDIR /go/src/$REPOSITORY/$ARTIFACT

# Build the binary
RUN make

#########################################
# Distribution stage
#########################################
FROM gcr.io/distroless/base-debian11

# Repository location
ARG REPOSITORY=github.com/ncarlier

# Artifact name
ARG ARTIFACT=readflow

# Install binary
COPY --from=builder /go/src/$REPOSITORY/$ARTIFACT/release/$ARTIFACT /usr/local/bin/$ARTIFACT
# Install website
COPY --from=ui-builder /usr/src/app/build /usr/share/html

# Add configuration file
ADD ./pkg/config/readflow.toml /etc/readflow.toml

# Set configuration file
ENV READFLOW_CONFIG /etc/readflow.toml

# Exposed ports
EXPOSE 8080 9090

# Define entrypoint
ENTRYPOINT [ "readflow" ]
