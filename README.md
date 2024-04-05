## This is my solution for the DevOps/SRE Engineer Test Task

### Default configuration parameters (hard-coded):

```
MAX_RANDOM_NUMBER = 20
PORT = 3000
```

### Prerequisites

- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
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

### Update Strategy

- set docker image path with a new version tag (v2) (replace with your own)
  `export IMAGE_PATH=famaten/my-test-task:v2`
- build image
  `docker build -t $IMAGE_PATH .`
- push image into docker registry
  `docker push $IMAGE_PATH`
- update kubernetes deployment image: [Kubernetes rolling update strategy](https://kubernetes.io/docs/tutorials/kubernetes-basics/update/update-intro/)
  `kubectl set image deploy/microservice-deploy microservice=$IMAGE_PATH`

## How to access the microservice

- access using [minikube tunnel](https://minikube.sigs.k8s.io/docs/handbook/accessing/)
  `minikube tunnel` - in separate terminal window
- check _EXTERNAL-IP_
  `kubectl get svc/microservice-svc`
- accesss microservice endpoints
  `curl http://EXTERNAL-IP:3000/health`
  `curl http://EXTERNAL-IP:3000/ready`
  `curl http://EXTERNAL-IP:3000/payload`
  `curl http://EXTERNAL-IP:3000/metrics`

## Load Testing

- create a new [python3 virtual env](https://packaging.python.org/en/latest/guides/installing-using-pip-and-virtual-environments/)
  `python3 -m venv loadtest/.venv`
- activate a new virtual env
  `source loadtest/.venv/bin/activate`
- install dependencies - _requests_ package
  `python3 -m pip install requests`
- execute load testing script
  `python3 loadtest/loadtest.py`
