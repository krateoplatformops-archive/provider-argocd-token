package token

import (
	"context"

	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/controller"

	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/ratelimiter"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	tokensv1alpha1 "github.com/krateoplatformops/provider-argocd-token/apis/tokens/v1alpha1"
	"github.com/krateoplatformops/provider-argocd-token/pkg/clients"
	"github.com/krateoplatformops/provider-argocd-token/pkg/clients/accounts"
)

const (
	errNotToken = "managed resource is not an argocd token custom resource"
	//errGetPC          = "cannot get ProviderConfig"
	//errFmtKeyNotFound = "key %s is not found in referenced Kubernetes secret"
)

// Setup adds a controller that reconciles Token managed resources.
func Setup(mgr ctrl.Manager, o controller.Options) error {
	name := managed.ControllerName(tokensv1alpha1.TokenGroupKind)

	log := o.Logger.WithValues("controller", name)

	recorder := mgr.GetEventRecorderFor(name)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(tokensv1alpha1.TokenGroupVersionKind),
		managed.WithExternalConnecter(&connector{
			kube: mgr.GetClient(),
			log:  log,
		}),
		managed.WithLogger(log),
		managed.WithRecorder(event.NewAPIRecorder(recorder)))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		WithOptions(o.ForControllerRuntime()).
		For(&tokensv1alpha1.Token{}).
		Complete(ratelimiter.NewReconciler(name, r, o.GlobalRateLimiter))
}

type connector struct {
	kube client.Client
	log  logging.Logger
}

func (c *connector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	cr, ok := mg.(*tokensv1alpha1.Token)
	if !ok {
		return nil, errors.New(errNotToken)
	}

	cfg, err := clients.GetConfig(ctx, c.kube, cr)
	if err != nil {
		return nil, err
	}

	c.log.Debug("Created session", "token", cfg.AuthToken)

	return &external{kube: c.kube, log: c.log, cfg: cfg}, nil
}

// An ExternalClient observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type external struct {
	kube client.Client
	log  logging.Logger
	cfg  *accounts.TokenProviderOptions
}

func (e *external) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*tokensv1alpha1.Token)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotToken)
	}

	spec := cr.Spec.ForProvider.DeepCopy()

	token, err := clients.GetSecret(ctx, e.kube, &spec.WriteTokenSecretToRef)
	if err != nil {
		if clients.ErrorIsNotFound(err) {
			// TODO handle token expiration?
			return managed.ExternalObservation{
				ResourceExists:   false,
				ResourceUpToDate: true,
			}, nil
		}
		return managed.ExternalObservation{}, err
	}

	if len(token) > 0 {
		// TODO handle token expiration?
		return managed.ExternalObservation{
			ResourceExists:   true,
			ResourceUpToDate: true,
		}, nil
	}

	return managed.ExternalObservation{
		ResourceExists:   false,
		ResourceUpToDate: true,
	}, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*tokensv1alpha1.Token)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotToken)
	}

	spec := cr.Spec.ForProvider.DeepCopy()

	token, err := accounts.GenerateToken(e.cfg, spec.Account, 0)
	if err != nil {
		return managed.ExternalCreation{}, err
	}
	e.log.Debug("Generated token", "account", spec.Account)

	if err := clients.SetSecret(ctx, e.kube, &spec.WriteTokenSecretToRef, token); err != nil {
		return managed.ExternalCreation{}, err
	}
	e.log.Debug("Saved token as secret", "account", spec.Account, "secret", spec.WriteTokenSecretToRef.Name)

	return managed.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	return managed.ExternalUpdate{}, nil // noop
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*tokensv1alpha1.Token)
	if !ok {
		return errors.New(errNotToken)
	}

	spec := cr.Spec.ForProvider.DeepCopy()

	e.log.Debug("Deleting token secret", "account", spec.Account, "secret", spec.WriteTokenSecretToRef.Name)

	return clients.DeleteSecret(ctx, e.kube, &spec.WriteTokenSecretToRef)
}
