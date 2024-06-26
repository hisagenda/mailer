# Email Operator

The Email Operator is a Kubernetes operator designed to manage the sending of emails through specified configurations. It automates the management of email configurations and sending processes using Custom Resource Definitions (CRDs) in Kubernetes.

## Prerequisites

Before deploying the Email Operator, you need to have the following installed:

- **Kubernetes Cluster**: Version 1.19 or higher.
- **kubectl**: Configured to communicate with your Kubernetes cluster.
- **Docker**: Required if you need to build your own images.
- **Kubebuilder**: For development and testing of the operator.

## Installation

Follow these steps to get your Email Operator up and running:

### 1. Clone the Repository

Clone this repository to your local machine:

```bash
git clone https://github.com/yourgithubrepo/email-operator.git
cd email-operator
```

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make docker-build docker-push IMG="abiola89/email-operator:v1.0.0"
```

Deploy the Operator
Deploy the operator to your Kubernetes cluster:

``` bash
make deploy IMG="abiola89/email-operator:v1.0.0"
```

**Usage and testing**

First, create a Kubernetes secret that contains the API token for tne email service provider.

Create an EmailSenderConfig resource to specify the sender's email and API token

Create a YAML file for the email you want to send. An example can be found in test email

```yaml
apiVersion: mail.mydomain.com/v1alpha1
kind: Email
metadata:
  name: test-email
  namespace: email-operator-system
spec:
  senderConfigRef: default-sender
  recipientEmail: recipient@example.com
  subject: "Test Email from Kubernetes"
  body: "Hello, this is a test email sent via the Email Operator."
```


Apply the Email resource to your cluster.
```bash
kubectl apply -f email.yaml
```
Logs are crucial for understanding the behavior of your operator. They can provide insights into how the operator processes the CRDs and communicates with the email service.

To view the logs of your Email Operator, use the following kubectl command:

```bash
kubectl logs -l name=email-operator -n email-operator-system
```
This command fetches the logs from all pods labeled with name=email-operator within the email-operator-system namespace.
You can also test from testmail.yaml

You can describe the crd using

`kubectl describe crd emails.mail.mailertest.com `

When a new Email is created, the operator should send the email using the provided sender configuration and MailerSend API


## K8s deployment

Kubernetes cluster was deployed using Terraform

# Terraform Configuration Files

## main.tf:

This file contains the Terraform configuration for creating an Azure Kubernetes Service (AKS) cluster. It includes the following resources:

- `random_pet.rg_name`: Generates a random name for the resource group.
- `azurerm_resource_group.rg`: Creates an Azure resource group using the random name generated by `random_pet.rg_name`.
- `azurerm_kubernetes_cluster.k8s`: Creates an AKS cluster in the resource group, using the generated resource group name and other specified configurations.

The AKS cluster is configured with the following settings:

- `default_node_pool`: Defines the default node pool for the cluster, including the name, VM size, and node count.
- `identity`: Enables system-assigned identity for the cluster.
- `linux_profile`: Specifies the Linux profile for the cluster, including the admin username and SSH key.
- `network_profile`: Configures the network settings for the cluster, including the network plugin and load balancer SKU.

## Initialize and Apply Terraform

To use this Terraform configuration, follow these steps:

1. Run `terraform init` to initialize the Terraform working directory. 
2. Run `terraform plan` to view planned changes and `terraform apply` to create or update the Azure resources defined in the configuration.

Make sure you have the necessary permissions and credentials to create resources in your Azure subscription.
