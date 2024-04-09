## :ledger: Index

- [:ledger: Index](#ledger-index)
- [:beginner: About](#beginner-about)
- [:zap: Usage](#zap-usage)
  - [:electric\_plug: Installation](#electric_plug-installation)
  - [:package: Commands](#package-commands)
- [:wrench: Development](#wrench-development)
  - [:notebook: Pre-Requisites](#notebook-pre-requisites)
  - [:nut\_and\_bolt: Development Environment](#nut_and_bolt-development-environment)
  - [:file\_folder: File Structure](#file_folder-file-structure)
  - [:rocket: Deployment](#rocket-deployment)
  - [:cactus: Branches](#cactus-branches)
- [:exclamation: Local Resources](#exclamation-local-resources)
  - [:cherry\_blossom: Grafana](#cherry_blossom-grafana)
  - [:tulip: Splunk](#tulip-splunk)
  - [:four\_leaf\_clover: Jaeger](#four_leaf_clover-jaeger)
- [:mag\_right: Future](#mag_right-future)
  - [:clock12: Infrastructure](#clock12-infrastructure)
  - [:clock3: Adapters](#clock3-adapters)
  - [:clock6: Frameworks](#clock6-frameworks)

##  :beginner: About
This project is a testament to the pursuit of learning and personal and professional growth in the field of software development. It showcases the configurability of different components, exemplifying Hexagonal Architecture in Golang through a simple CRUD application.

Configured using docker-compose, the project offers a hands-on exploration of various infrastructure components such as logs management, metrics tracking, and traces monitoring. It serves as both a learning resource and a practical demonstration of architectural principles in action.

## :zap: Usage
Quick initialization to how the project can be used.

###  :electric_plug: Installation
To install and set up the project on Mac OSX, follow these straightforward steps:

1. Install Docker Desktop.
```
$ brew install docker
```
2. Intall make.
```
$ brew install make
```

###  :package: Commands
There exist a multitude of possibilities for initiating the project, contingent upon the specific areas of interest one wishes to delve into or experiment with.
- All.
```make
make start compose=all
```
- Splunk logging through Fluent-Bit.
```make
make start compose=fluent-bit
```
- Prometheus metrics and representation on Grafana.
```make
make start compose=grafana
```
- OTEL traicing and representation on Jaeger.
```make
make start compose=jaeger
```
- Splunk logging through docker-compose logging driver.
```make
make start compose=splunk
```
You can stop the docker-compose using the following command specifying the compose you used to start.
```make
make stop compose=all
```
Or you can restart the docker-compose using the following command specifying the compose you used to start.
```make
make restart compose=all
```

##  :wrench: Development
This is a personal project and I want to have it as kind of portfolio.

### :notebook: Pre-Requisites
List all the pre-requisites the system needs to develop this project.
- Visual Studio Code (Go and Go Nightly extensions installed correctly).
- Visual Studio Code Docker extension.
- Visual Studio Code Database Client extension (suggested for all the PostgreSQL, ElasticSearch, Redis, Kafkam, etc management). 

###  :nut_and_bolt: Development Environment
To set up the local working environment for the project just follow the instructions.
- To download the project yo can use git clone.
```
git clone https://github.com/denyssydorenko/hexagonal-architecture-utils.git
```
- To install the dependencies you can just run the following command:
```
go mod tidy
go mod vendor
```
- You can start the project with the commands provided above.

###  :file_folder: File Structure
```
.
├── .vscode
├── build
│   ├── compose
│   │   ├── all
│   │   │   └── docker-compose.yml
│   │   ├── fluent-bit
│   │   │   └── docker-compose.yml
│   │   ├── grafana
│   │   │   └── docker-compose.yml
│   │   ├── jaeger
│   │   │   └── docker-compose.yml
│   │   └── splunk
│   │   │   └── docker-compose.yml
│   ├── ddl
│   │   ├── db
│   │   |   └── schema
│   │   │   │   └── init.sql
│   │   ├── fluent-bit
│   │   │   └── fluent-bit.conf
│   │   ├── grafana
│   │   |   └── provisioning
│   │   |       ├── dashboards
│   │   │       │   └── dashboard.yml
│   │   |       └── datasources
│   │   │           └── datasource.yml
│   │   ├── otel-collector
│   │   │   └── otel-collector-config.yaml
│   │   ├── prometheus
│   │   │   └── prometheus.yml
│   │   └── splunk
│   │   │   └── splunk.yml
│   ├── packages
│   │   └── api
│   │       └── Dockerfile
├── cmd
│   └── api
│       └── main.go
├── config
│   ├── app.go
│   ├── config.go
│   ├── config.yaml
│   └── infra.go
├── dashboards
│   └── api.json
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── adapters
│   │   ├── db
│   │   │   ├── db_test.go
│   │   │   ├── db.go
│   │   │   └── dbmock.go
│   │   └── http
│   │       ├── metrics
│   │       │   └── metrics.go
│   │       └── http.go
│   ├── application
│   │   └── core
│   │       └── api
│   │           ├── application_test.go
│   │           └── application.go
│   ├── domains
│   │   ├── db
│   │   │   └── db.go
│   │   ├── health
│   │   │   └── health.go
│   │   └── api.go
│   ├── pkg
│   │   ├── logging
│   │   │   └── logging.go
│   │   └── otel
│   │       └── otel.go
│   └── ports
│       ├── api.go
│       └── db.go
├── postman
│   ├── Hexagonal Architecture Utils.postman_collection.json
│   └── localhost.postman_environment.json
├── vendor
├── .dockerignore
├── .env
├── .gitignore
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

### :rocket: Deployment
In the future will create a local minikube. (pending)

### :cactus: Branches

I use an agile continuous integration methodology, so the version is frequently updated and development is really fast.

1. **`master`** is the production branch.

2. No other permanent branches should be created in the main repository, you can create feature branches but they should get merged with the master.

**Steps to work with feature branch**

1. To start working on a new feature, create a new branch prefixed with `features` and followed by feature name. (ie. `features/FEATURE-NAME`)
2. Once you are done with your changes, you can raise PR.

**Steps to create a pull request**

1. Make a PR to `master` branch.
2. Comply with the best practices and guidelines e.g. where the PR concerns visual elements it should have an image showing the effect.
3. It must pass all continuous integration checks and get positive reviews. (pending)

After this, changes can be merged.

## :exclamation: Local Resources
### :cherry_blossom: Grafana
- URL: https://localhost:3000
- Username: admin
- Password: admin
### :tulip: Splunk
- URL: https://localhost:8000
- Username: admin
- Password: changeme
### :four_leaf_clover: Jaeger
- URL: https://localhost:16686

##  :mag_right: Future
I want to add more things to this project in order to have all the technologies I know how to use in the same place.
### :clock12: Infrastructure
- Logging ElasticSearch/Kibana
- Logging Grafana Loki
- Traces Grafana Tempo
- Metrics Grafana Mimir
### :clock3: Adapters
- Redis
- ElasticSearch
- Amazon S3
- Nats producer/consumer
- Nats Jetstream producer/consumer
- Kafka producer/consumer
### :clock6: Frameworks
- gRPC
- Websocket
