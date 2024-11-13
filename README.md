# vllm-k8s-operator 🚀 [WIP!]

Welcome to the **vllm-k8s-operator**! This open-source Kubernetes operator is designed to streamline the deployment, scaling, and management of vLLM (vectorized Language Models) in Kubernetes clusters. Whether you're deploying large LLMs or managing complex workloads, our operator makes it simple, efficient, and scalable.

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

The **vllm-k8s-operator** provides Kubernetes-native automation for deploying and managing vLLM models with a focus on scalability, configurability, and optimized resource usage. With support for custom configuration, flexible deployment options, and seamless integration with existing Kubernetes workflows, this operator makes it easy for platform teams and DevOps engineers to manage vLLMs at scale.

---

## Features 🔍

- **Seamless Deployment**: Quickly deploy vLLM models on Kubernetes using CRDs.
- **Scalable Management**: Scale replicas effortlessly to meet demand.
- **Customizable Configurations**: Easily configure model parameters, ports, GPU utilization, and more.
- **Resource Optimization**: Efficiently use GPU and memory resources.
- **Open-Source Community**: Join and contribute to a vibrant community of developers.

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
kubectl apply -f https://raw.githubusercontent.com/<your-repo>/vllm-k8s-operator/main/deploy/crds.yaml
kubectl apply -f https://raw.githubusercontent.com/<your-repo>/vllm-k8s-operator/main/deploy/operator.yaml
```
