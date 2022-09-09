# ticker

`ticker` is a simple web API that fetches upstream data regarding stocks using the stock symbol.

The referenced upstream API in this example is https://www.alphavantage.co/.

- [ticker](#ticker)
  - [Setup](#setup)
    - [Docker](#docker)
    - [Source](#source)
  - [Environment Variables](#environment-variables)
  - [Documentation](#documentation)
  - [Usage](#usage)
  - [Deployment](#deployment)
  - [Resiliency](#resiliency)
  - [Caveats](#caveats)

## Setup

### Docker

A Dockerfile is provided for building locally. The image is already built and available on Docker Hub via `eb129/ticker:v0.1.1`.

If you wish to build it yourself:

```bash
$ docker build . -t ticker:v0.1.0
```

You can adjust the image name and tag as you see fit, but be sure to update those when attempting to run via Docker.

### Source

Alternatively, you can build the binary locally without the image and run:

```bash
$ make build
$ bin/ticker_darwin_amd64
```

The application runs on port `8080` - you can adjust this in `main.go` if you see fit.

Running this from the root of the repository will work since `.env` is available if you have populated it. You can optionally use the `ticker_linux_amd64` build depending on your environment. If you require compatibility for `arm` arch, you may adjust the contents of the Makefile.

## Environment Variables

The following environment variables are required via the `.env` file if running locally. Alternatively, the `.env` will be skipped if it isn't present and environment variables will be referenced directly. See the section on deployment below for additional details. All of these are required to be set.

- `APIKEY` - The API key for the alphavantage API.
- `SYMBOL` - The stock symbol to look up - i.e. `MSFT`.
- `NDAYS` - When using the `/avg` endpoint, this is the number of days to look back for data.

## Documentation

Swagger docs are available via http://localhost:8080/swagger/index.html when the API is running.

## Usage

Once the API is available either locally or through deployment to Kubernetes, you can use the following calls to test. The endpoint will need to be adjusted depending on where you run it.

```bash
$ curl -X 'GET' \
  'http://localhost:8080/api/v1/stock' \
  -H 'accept: application/json'
```

This will return the full range of data available for the stock, bypassing `NDAYS`.

```bash
$ curl -X 'GET' \
  'http://localhost:8080/api/v1/stock/avg' \
  -H 'accept: application/json'
```

This will return any items greater than or equal to the `NDAYS` value - meaning, anything older will not be returned. This also returns the average close price over that same time period.

## Deployment

Kustomize manifests are provided for deployment. You can switch to the `deploy/kubernetes/ticker/overlays/development` folder and run the following commands to work with the files.

You must have a Kubernetes cluster running and configured for the `apply` step. This was tested in development using `minikube`.

```bash
$ kubectl kustomize .
...outputs rendered manifests
```

```bash
$ kubectl apply -k .
...applies to cluster
```

Once the application is deployed, you can port-forward directly to the service if you wish to test locally:

```bash
$ kubectl port-forward -n ticker service/ticker 8080:8080
$ curl...
```

**Note:** The manifests use a secret referencing a file called `.env.api` in the `base/` folder that will generate a `Secret` using `secretGenerator`. For the sake of this demonstration, `SYMBOL` and `NDAYS` are provided as literals for the `configMapGenerator`. You must provide your API key in `.env.api` and you may adjust the literals in the `configMapGenerator` as you wish.

## Resiliency

A `HorizontalPodAutoscaler` and a `PodDisruptionBudget` are provided here as they are fundamental constructs for managing resiliency in a real-world application. A base health check (liveness/readiness) is also provided. These settings would need to be adjusted in production to reflect the following requirements:

- Enough `Pods` at baseline to support realistic traffic on the `HPA` in terms of min/max settings.
- Appropriate configuration of the `PDB` to guarantee the application remains available during events like rollouts, cluster upgrades, etc.
- Appropriate configuration of `RollingUpdate` configuration on the `Deployment` to tolerate graceful rollouts.
- Appropriate configuration of health checks.

The configuration provided here is oversimplified, but this would warrant additional consideration when deploying an application.

## Caveats

- In production, you'd want to use a real data store for rate limiting implementation. This uses an in-memory store so restarting the app negates the usefulness of rate limiting. The rate limiting implemenation is 5/minute - the alphavantage API also instates 500/day, but that wasn't implemented here.
- Parsing of sensitive data, like the API key or anything passed in from `.env`, is done in the function calls to save time. In a real-world application, I'd recommend using a struct to instantiate the data so it can be reused properly.
- This was tested on `minikube` using `minikube v1.25.2 on Darwin 12.4`.
- I would recommend writing tests for this application in real-world usage.
- `Ingress` configuration, in my opinion, is beyond the scope of this exercise. You may need to add annotations or adjust the manifest for the `Ingress` depending on your environment.
- `HorizontalPodAutoscaler` functionality only works if metrics are available via the API. This is beyond the scope of this exercise.
