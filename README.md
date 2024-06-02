# email-operator
This is a sample mail operator implementation for email operator

## Description
Kubernetes operator that manages custom resources for configuring email sending and sending of emails via a transactional email provider like MailerSend.

## Getting Started

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster or KubeMini.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/email-operator:tag

ex:
make docker-build docker-push IMG=mustaphasaeed/email-operator:latest
```

**Create the desired namespace**

```sh
kubectl apply -f config/default/namespace.yaml 
```

**Deploy the operator to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/email-operator:tag MAIL_SEND_TOKEN=<BASE64_OF_MAILER_SEND_TOKEN>
```
