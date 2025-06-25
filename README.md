# Custom-k8s-Controller

**Custom Kubernetes Controller written in Go**

---

[![[Build Status]](https://github.com/silhouetteUA/k8s-controller/actions/workflows/ci.yaml/badge.svg)](...)
[![Trivy Scan](https://github.com/silhouetteUA/k8s-controller/actions/workflows/ci.yaml/badge.svg)](...)
[![Go Version](https://img.shields.io/badge/go-1.24.4-blue)](https://golang.org/)
[![License](https://img.shields.io/github/license/silhouetteUA/k8s-controller)](https://github.com/silhouetteUA/k8s-controller/blob/feature/step5-ci/LICENSE)
[![Docker](https://img.shields.io/badge/docker-ghcr.io%2Fsilhouetteua%2Fk8s--controller-blue)](...)


---

## 🚀 Overview

This project contains a custom Kubernetes controller implemented in Go.  
It is designed for learning, experimentation, or extension into production-grade components.

---

## 🛠️ Development Environment Setup

To set up a development control plane environment, you can use the following repository:

🔗 [silhouetteUA/kubernetes-controlplane](https://github.com/silhouetteUA/kubernetes-controlplane)

This repository provides a lightweight Kubernetes control plane setup suitable for local testing and development of Kubernetes controllers.

### 🔄 Alternative Environments

Alternatively, **any standard Kubernetes environment** can be used for development, including:

- [k3s](https://k3s.io/) or [k3d](https://k3d.io/)
- [KIND (Kubernetes IN Docker)](https://kind.sigs.k8s.io/)
- [Minikube](https://minikube.sigs.k8s.io/)
- Any **managed Kubernetes** service provided by major cloud platforms (e.g., EKS, GKE, AKS)

Choose the setup that best fits your local or cloud-based workflow.

---

## 📦 Features

- Written in Go
- Follows Kubernetes controller-runtime patterns
- Easily extensible and testable
- Ideal for CRD experimentation and controller logic development

---

## 📄 License

MIT License  
© 2025 silhouetteUA

---
