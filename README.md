# Provider ArgoCD

## Overview

This is a Kubernetes Operator (Crossplane provider) that creates API tokens for specific [ArgoCD](https://argo-cd.readthedocs.io/) users.

The provider that is built from the source code in this repository adds the following new functionality:

- a Custom Resource Definition (CRD) that model ArgoCD auth tokens for specific users

## Getting Started

With Crossplane and ArgoCD installed in your cluster:

```sh
$ helm repo add crossplane-stable https://charts.crossplane.io/stable
$ helm repo update
$ helm install crossplane --namespace crossplane-system crossplane-stable/crossplane
```

```sh
$ kubectl create namespace argo-system
$ kubectl apply -n argo-system -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
```

### How to install this provider

```sh
$ kubectl apply -f ./examples/provider.yaml
```

Replace `VERSION` tag with the desired release.

### Configure this operator with `serverUrl` pointing to an ArgoCD instance

```sh
$ cat <<EOF | kubectl apply -f -
apiVersion: argocd.krateoplatformops.io/v1alpha1
kind: ProviderConfig
metadata:
  name: provider-argocd-token-config
spec:
  serverUrl: https://argocd-server.argo-system.svc:443
  credentials:
    source: Secret
    secretRef:
      namespace: argo-system
      name: argocd-initial-admin-secret
      key: password
EOF
```

### Create a new ArgoCD account

Following the steps in the [official ArgoCD documentation](https://argo-cd.readthedocs.io/en/stable/operator-manual/user-management/#create-new-user) you can create a new user defining it in the `argo-cm` ConfigMap:

```sh
$ cat <<EOF | kubectl apply -f -
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
EOF
```

Each user might have two capabilities:

- apiKey: allows generating authentication tokens for API access
- login: allows to login using UI

### Create an API token without expiration that can be used by the defined user

```sh
$ cat <<EOF | kubectl apply -f -
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
    name: provider-argocd-token-config
EOF
```

### After a while check if the API token is created

```sh
$ kubectl get secrets/krateo-dashboard-argocd-token -n krateo-system \
   --template='{{.data.authToken | base64decode}}'

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJqdGkiOiJkOWZkNDJiYi05ZGU4LTRmMGUtYTA...
```

---



## Contributing

This is a community driven project and we welcome contributions.

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please open an [issue](https://github.com/krateoplatformops/provider-argocd-token/issues).