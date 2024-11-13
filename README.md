# vllm-k8s-operator 🚀 [WIP!]

Welcome to the **vllm-k8s-operator**! This open-source Kubernetes operator is designed to streamline the deployment, scaling, and management of LLMs served via vLLM in Kubernetes clusters. The goal is to build an inferencing operator that DevOps teams can use to deploy and manage open-source models. The operator should integrate seamlessly with your existing infrastructure, avoiding the need to set up something entirely new.

---

## Table of Contents 📚

- [Overview](#overview)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Quick Start](#quick-start)
- [Configuration](#configuration)
  - [Spec Reference](#spec-reference)
- [Contributing](#contributing)
- [Community](#community)
- [License](#license)

---

## Overview ✨

The **vllm-k8s-operator** provides Kubernetes-native automation for deploying and managing vLLM models with a focus on scalability, configurability, and automated inference optimization. You can define your model requirements as Kubernetes manifests, and the operator will deploy those models with the optimal settings to bring them to life

---

## Features 🔍

- **Seamless Deployment**: Quickly deploy vLLM models on Kubernetes using CRDs.
- **Autoscaling**: Scale replicas effortlessly to meet demand.
- **GPU Optimization**: Efficiently use GPU and memory resources.

---

## Getting Started 🛠️

### Prerequisites

Ensure you have the following prerequisites installed:

- [Kubernetes](https://kubernetes.io) (v1.21 or later)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [Helm](https://helm.sh/) (optional for installation)

### Installation

Install the **vllm-k8s-operator** using `kubectl`:

```bash
make install
```

### Quick Start 🚀

Create a VllmDeployment resource to deploy a vLLM model:

```bash
apiVersion: core.vllmoperator.org/v1alpha1
kind: VllmDeployment
metadata:
  name: example-model
spec:
  replicas: 1
  model:
    name: "huggingface-model-name"
    hf_url: "https://huggingface.co/model-url"
  vLLMConfig:
    port: 8072
    gpu-memory-utilization: "0.75"
    log-level: "info"
    block-size: 16
    max-model-len: 2000
    enforce-eager: true
  containers:
    - name: vllm
      image: vllm/vllm-openai:v0.6.2
      ports:
        - containerPort: 8072
```

Apply the resource:

```bash
kubectl apply -f example-vllmdeployment.yaml
```

### Configuration ⚙️

**VllmDeployment Fields**

    •	replicas (integer): Number of replicas for the vLLM deployment.
    •	model (object):
    •	name (string): Name of the model.
    •	hf_url (string): URL to the model on Hugging Face or similar.
    •	vLLMConfig (object):
    •	port (integer): Port for the vLLM service.
    •	gpu-memory-utilization (string): GPU utilization ratio.
    •	log-level (string): Log level for the service.
    •	block-size (integer): Block size.
    •	max-model-len (integer): Maximum model length.
    •	enforce-eager (boolean): Enforce eager execution.
    •	containers (array): List of container specifications.
    •	name (string): Name of the container.
    •	image (string): Container image.
    •	ports (array): List of ports exposed by the container.

### Contributing 🤝

We ❤️ contributions! If you’d like to contribute to the **vllm-k8s-operator**, please take a look at our contribution guidelines. Contributions can include:

- Bug fixes 🐛
- New features 🌟
- Documentation updates 📚

To get started:

1. Fork the repository.

2. Create a new branch for your feature or bug fix.

3. Commit your changes with clear commit messages.

4. Open a pull request describing your changes.

### Community 💬

Join the discussion and connect with other developers:

- GitHub Issues: Report bugs or suggest features on our issue tracker.
- GitHub Discussions: Engage with the community on GitHub Discussions.
- Slack: Join our Slack community to chat with other users and contributors.

### License 📜

This project is licensed under the Apache License 2.0.
