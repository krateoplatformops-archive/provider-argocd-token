package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// TokenObservation are the observable fields of a Token.
type TokenObservation struct {
	ID        string `json:"id,omitempty"`
	ExpiresIn string `json:"expiresIn,omitempty"`
}

// TokenParameters are the configurable fields of of a Token.
type TokenParameters struct {
	// ID optional token id. Fall back to uuid if not value specified
	// +optional
	ID string `json:"id,omitempty"`

	// Account name
	Account string `json:"account"`

	// ExpiresIn duration before the token will expire. (Default: No expiration)
	// +optional
	// ExpiresIn string `json:"expiresIn,omitempty"`

	WriteTokenSecretToRef xpv1.SecretKeySelector `json:"writeTokenSecretToRef"`
}

// A TokenSpec defines the desired state of a Token.
type TokenSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       TokenParameters `json:"forProvider"`
}

// A TokenStatus represents the observed state of a Token.
type TokenStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          TokenObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,argocd}
// +kubebuilder:subresource:status
type Token struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TokenSpec   `json:"spec"`
	Status TokenStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TokenList contains a list of Token
type TokenList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Token `json:"items"`
}
