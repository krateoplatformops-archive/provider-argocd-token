package clients

import (
	"context"
	"strings"

	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/krateoplatformops/provider-argocd-token/apis/v1alpha1"
	"github.com/krateoplatformops/provider-argocd-token/pkg/clients/accounts"
)

const (
	argocdInititalAdminSecret = "argocd-initial-admin-secret"
)

// GetConfig constructs a ClientOptions configuration that can be used to authenticate to argocd
// API by the argocd Go client
func GetConfig(ctx context.Context, c client.Client, mg resource.Managed) (*accounts.TokenProviderOptions, error) {
	switch {
	case mg.GetProviderConfigReference() != nil:
		return UseProviderConfig(ctx, c, mg)
	default:
		return nil, errors.New("providerConfigRef is not given")
	}
}

// UseProviderConfig to produce a config that can be used to create an ArgoCD client.
func UseProviderConfig(ctx context.Context, k client.Client, mg resource.Managed) (*accounts.TokenProviderOptions, error) {
	pc := &v1alpha1.ProviderConfig{}
	if err := k.Get(ctx, types.NamespacedName{Name: mg.GetProviderConfigReference().Name}, pc); err != nil {
		return nil, errors.Wrap(err, "cannot get referenced Provider")
	}

	t := resource.NewProviderConfigUsageTracker(k, &v1alpha1.ProviderConfigUsage{})
	if err := t.Track(ctx, mg); err != nil {
		return nil, errors.Wrap(err, "cannot track ProviderConfig usage")
	}

	opts := &accounts.TokenProviderOptions{
		ServerAddr: pc.Spec.ServerAddr,
	}

	pass, err := GetInitialAdminPassword(ctx, k, pc)
	if err != nil {
		return nil, err
	}

	token, err := accounts.Login(opts, "admin", pass)
	if err != nil {
		return nil, err
	}

	opts.AuthToken = token

	return opts, nil
}

// GetInitialAdminPassword returns the ArgoCD initial admin password.
func GetInitialAdminPassword(ctx context.Context, k client.Client, pc *v1alpha1.ProviderConfig) (string, error) {
	ref := &xpv1.SecretKeySelector{
		SecretReference: xpv1.SecretReference{
			Name: argocdInititalAdminSecret,
		},
		Key: corev1.BasicAuthPasswordKey,
	}

	if pc.Spec.Credentials != nil {
		if s := pc.Spec.Credentials.Source; s != xpv1.CredentialsSourceSecret {
			return "", errors.Errorf("credentials source %s is not currently supported", s)
		}

		csr := pc.Spec.Credentials.SecretRef
		if csr != nil {
			if name := strings.TrimSpace(csr.SecretReference.Name); len(name) > 0 {
				ref.SecretReference.Name = name
			}

			if namespace := strings.TrimSpace(csr.SecretReference.Namespace); len(namespace) > 0 {
				ref.SecretReference.Namespace = namespace
			}

			if key := strings.TrimSpace(csr.Key); len(key) > 0 {
				ref.Key = csr.Key
			}
		}
	}

	return GetSecret(ctx, k, ref)
}

func SaveAdminToken(ctx context.Context, k client.Client, pc *v1alpha1.ProviderConfig, token string) error {
	if s := pc.Spec.Credentials.Source; s != xpv1.CredentialsSourceSecret {
		return errors.Errorf("credentials source %s is not currently supported", s)
	}

	csr := pc.Spec.Credentials.SecretRef
	if csr == nil {
		return errors.New("no credentials secret referenced")
	}

	return SetSecret(ctx, k, csr.DeepCopy(), token)
}
