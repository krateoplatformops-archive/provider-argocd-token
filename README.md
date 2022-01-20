# Provider ArgoCD

## Overview

This is a Kubernetes Operator (Crossplane provider) that creates API tokens for specific [ArgoCD](https://argo-cd.readthedocs.io/) users.

The provider that is built from the source code in this repository adds the following new functionality:

- a Custom Resource Definition (CRD) that model ArgoCD auth tokens for specific users

## Getting Started

With ArgoCD installed in your cluster, for example...

```sh
$ kubectl create namespace argo-system
$ kubectl apply -n argo-system -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

### Configure this operator with `serverAddr` pointing to an ArgoCD instance

```sh
apiVersion: argocd.krateoplatformops.io/v1alpha1
kind: ProviderConfig
metadata:
  name: provider-argocd-config
spec:
  serverAddr: argocd-server.argocd.svc:443
  insecure: true
  portForward: true
  portForwardNamespace: argo-system
  credentials:
    source: Secret
    secretRef:
      namespace: argo-system
      name: argocd-initial-admin-secret
      key: password
```

```sh
$ kubectl apply -f examples/provider-argocd-token-config.yaml
```

### Create a new user

Following the steps in the [official ArgoCD documentation](https://argo-cd.readthedocs.io/en/stable/operator-manual/user-management/#create-new-user) you can create a new user defining it in the `argo-cm` ConfigMap:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-cm
  namespace: argo-system
  labels:
    app.kubernetes.io/name: argocd-cm
    app.kubernetes.io/part-of: argocd
data:
  accounts.krateo-dashboard: apiKey, login
```

Each user might have two capabilities:

- apiKey: allows generating authentication tokens for API access
- login: allows to login using UI

```sh
$ kubectl apply -f examples/argocd-cm.yaml
```

### Create an API token without expiration that can be used by the defined user

```yaml
apiVersion: argocd.krateoplatformops.io/v1alpha1
kind: Token
metadata:
  name: krateo-dashboard-argocd-token
spec:
  forProvider:
    account: krateo-dashboard
    writeTokenSecretToRef:
      name: krateo-dashboard-argocd-token
      key: authToken
      namespace: krateo-system
  providerConfigRef:
    name: provider-argocd-config
```

```sh
$ kubectl apply -f examples/krateo-dashboard-argocd-token.yaml
```

## Installing this provider using Helm

```sh
$ helm repo add krateo-runtime-providers https://krateoplatformops.github.io/krateo-runtime-providers 
$ helm repo update
$ helm install provider-argocd-token --namespace $(NAMESPACE) krateo-runtime-providers/argocd-token
```

---


## Contributing

This is a community driven project and we welcome contributions.

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please open an [issue](https://github.com/krateoplatformops/provider-argocd-token/issues).