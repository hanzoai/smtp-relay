# smtp-relay
Simple HTTP API to relay messages to SMTP servers.


## Installation
Install deps:
```
$ make deps
```

Install the following command-line tools:

- gcloud is used to create and delete Kubernetes Engine clusters. gcloud is included in the Google Cloud SDK.
- kubectl is used to manage Kubernetes, the cluster orchestration system used by Kubernetes Engine. You can install kubectl using gcloud:

```
$ gcloud components install kubectl
```

Set some defaults:
```
gcloud config set project crowdstart-us
gcloud config set compute/zone us-central1-b
```

## Usage
```
$ make run
```

You can test the server is running by sending an email payload to it:

```
$ curl -i --user user:pass -H "Content-Type: application/json" -H "Accept: application/json" -X POST -d '{"username":"foo", "password":"asdf", "host":"", "port":"","mailfrom":"","mailto":[],"msg":""}' https://smtp-relay.hanzo.ai
```

## Deployment

Create cluster named `smtp-relay` with 2 nodes:
```
$ gcloud container clusters create smtp-relay --num-nodes=2
Creating cluster smtp-relay...done.                                                                    Created [https://container.googleapis.com/v1/projects/crowdstart-us/zones/us-central1-b/clusters/smtp-relay].
kubeconfig entry generated for smtp-relay.
NAME        ZONE           MASTER_VERSION  MASTER_IP       MACHINE_TYPE   NODE_VERSION  NUM_NODES  STATUS
smtp-relay  us-central1-b  1.7.8-gke.0     35.202.117.238  n1-standard-1  1.7.8-gke.0   2          RUNNING
```

You can list clusters for in your project or get the details for a single cluster using the following commands:

```
$ gcloud container clusters list
$ gcloud container clusters describe smtp-relay
```

To create a deployment:

```
$ kubectl create -f deployment.yaml
```

Verify the three replicas are running by querying the list of the labels that identify the web frontend:

```
kubectl get pods -l app=smtp-relay
```
