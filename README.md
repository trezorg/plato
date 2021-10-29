# Palto devops test

## Installations

1. GO111MODULE=on go install github.com/trezorg/plato/cmd/fetcher@latest

### Usage example

    fetcher http://example.com https://google.com https://facebook.com https://some_wrong_host

## k8s app architecture

### Requirements

You need to design the architecture of the system that will run in kubernetes. You can
use any desired format (text, UML diagrams, k8s manifests, hand drawings, etc.) to
share the result.

1. SPA frontend
2. API backend
3. Postgres cluster
4. S3 bucket
5. external data provider (JSON, HTTP/1.1)
6. SQL script for initializing the database
7. binary that creates fixtures

### Solution

Suppose we have bare metal k8s installation

1. SPA. Development. Access via nginx-ingress. Node port service
2. API backend. Development. k8s service, cluster IP
3. Postgres cluster. Postgres operator. headless service
4. S3 bucket. StatefulSet.
5. external data provider (JSON, HTTP/1.1). External name service
6. SQL script for initializing the database. Helm hook with post-install, pre-upgrade
7. binary that creates fixtures. SizeCar container for API backend pod. Share volume with API backend container.