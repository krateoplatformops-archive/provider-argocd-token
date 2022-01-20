module github.com/krateoplatformops/provider-argocd-token

go 1.16

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/crossplane/crossplane-runtime v0.15.1-0.20211202230900-d43d510ec578
	github.com/crossplane/crossplane-tools v0.0.0-20210916125540-071de511ae8e
	github.com/onsi/gomega v1.15.0 // indirect
	github.com/pkg/errors v0.9.1
	golang.org/x/mod v0.5.1-0.20210830214625-1b1db11ec8f4 // indirect
	golang.org/x/net v0.0.0-20210614182718-04defd469f4e // indirect
	golang.org/x/oauth2 v0.0.0-20210628180205-a41e5a781914 // indirect
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver v0.22.2 // indirect
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	sigs.k8s.io/controller-runtime v0.9.6
	sigs.k8s.io/controller-tools v0.6.2
)

replace (
	k8s.io/api => k8s.io/api v0.22.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.22.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.22.4-rc.0
	k8s.io/apiserver => k8s.io/apiserver v0.22.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.22.2
	k8s.io/client-go => k8s.io/client-go v0.22.2 // Required by prometheus-operator
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.22.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.22.2
	k8s.io/code-generator => k8s.io/code-generator v0.22.4-rc.0
	k8s.io/component-base => k8s.io/component-base v0.22.2
	k8s.io/component-helpers => k8s.io/component-helpers v0.22.2
	k8s.io/controller-manager => k8s.io/controller-manager v0.22.2
	k8s.io/cri-api => k8s.io/cri-api v0.23.0-alpha.0
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.22.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.22.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.22.2
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.22.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.22.2
	k8s.io/kubectl => k8s.io/kubectl v0.22.2
	k8s.io/kubelet => k8s.io/kubelet v0.22.2
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.22.2
	k8s.io/metrics => k8s.io/metrics v0.22.2
	k8s.io/mount-utils => k8s.io/mount-utils v0.22.4-rc.0
	k8s.io/node-api => k8s.io/node-api v0.21.1
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.22.2
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.22.2
	k8s.io/sample-controller => k8s.io/sample-controller v0.22.2
)

replace k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.22.2

//replace k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20210323165736-1a6458611d18
