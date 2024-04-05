## This is my solution for the DevOps/SRE Engineer Test Task

### My considerations and observations of the solution

- Security
  - The final Docker image runs by non-root user
- Reliability
  - The health checks and readiness probes are in the kubernetes manifest file to ensure that pods are healthy and ready to serve traffic before receiving requests.
  - There are appropriate resource requests and limits to prevent resource contention and ensure optimal utilization of cluster resources.
  - There is rolling updates strategy implemented - see below - to minimize downtime during application updates and ensure smooth transitions between different versions of your application.
- Scalability
  - There is a Kubernetes deployment and horizontal pod autoscaling (HPA) to automatically scale the number of replicas based on workload demand, ensuring that the microservice can handle increased traffic and load spikes. (replica count from 3 to 10 based on CPU utilization)
- How to monitor and log the performance and errors of your service
  - There is `/metrics` endpoint to retrieve basic metrics of the microservice
  - `kubectl logs` can provide logs from the container
- Assumptions made for simplicity
  - My own metrics implementation despite [prometheus go lang package](https://pkg.go.dev/github.com/prometheus/client_golang/prometheus)
  - Manual application versioning and docker image publishing to my personal docker hub registry
  - No CI/CD tool therefore manual commands execution is required (see guide below)
  - No unit tests within golang application

### Prerequisites

- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [minikube addons enable metrics-server](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale-walkthrough/)
- [Docker Desktop](https://www.docker.com/products/docker-desktop/)
- [Python3](https://www.python.org/downloads/)

### Initial Deployment

- set docker image path (replace with your own)
  `export IMAGE_PATH=famaten/my-test-task`
- build image
  `docker build -t $IMAGE_PATH .`
- push image into docker registry
  `docker push $IMAGE_PATH`
- create deployment from the kubernetes manifest
  `kubectl apply -f kubernetes/microservice.yaml`

### Updates Strategy

- set docker image path with a new version tag (v2) (replace with your own)
  `export IMAGE_PATH=famaten/my-test-task:v2`
- build image
  `docker build -t $IMAGE_PATH .`
- push image into docker registry
  `docker push $IMAGE_PATH`
- update kubernetes deployment image: [Kubernetes rolling update strategy](https://kubernetes.io/docs/tutorials/kubernetes-basics/update/update-intro/)
  `kubectl set image deploy/microservice-deploy microservice=$IMAGE_PATH`

## How to access the microservice

- access using [minikube tunnel](https://minikube.sigs.k8s.io/docs/handbook/accessing/) (in a separate terminal window)
- check _EXTERNAL-IP_
  `kubectl get svc/microservice-svc`
- accesss microservice endpoints
  - `curl http://EXTERNAL-IP:3000/health`
  - `curl http://EXTERNAL-IP:3000/ready`
  - `curl http://EXTERNAL-IP:3000/payload`
  - `curl http://EXTERNAL-IP:3000/metrics`

## Load Testing

- create a new [python3 virtual env](https://packaging.python.org/en/latest/guides/installing-using-pip-and-virtual-environments/)
  `python3 -m venv loadtest/.venv`
- activate a new virtual env
  `source loadtest/.venv/bin/activate`
- install dependencies: the (`requests`) package
  `python3 -m pip install requests`
- execute load testing script
  `python3 loadtest/loadtest.py`
- check results - see sample below

```
Endpoint                 RequestCount             TotalDuration(ms)        AverageLatency(ms)
health                   37                       0.148                    0.004
metrics                  31                       0.637                    0.021
payload                  27                       0.479                    0.018
ready                    34                       0.192                    0.006
```

### Default configuration parameters (hard-coded in the source code)

- maximum random number for fibonacci sequence
  `MAX_RANDOM_NUMBER = 20`
- port number exposed by microservice
  `PORT = 3000`
