package v1alpha1

import (
	"reflect"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// Package type metadata.
const (
	Group   = "argocd.krateoplatformops.io"
	Version = "v1alpha1"
)

var (
	// SchemeGroupVersion is group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: Group, Version: Version}
	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	SchemeBuilder = &scheme.Builder{GroupVersion: SchemeGroupVersion}
)

// Task type metadata
var (
	TokenKind             = reflect.TypeOf(Token{}).Name()
	TokenGroupKind        = schema.GroupKind{Group: Group, Kind: TokenKind}.String()
	TokenKindAPIVersion   = TokenKind + "." + SchemeGroupVersion.String()
	TokenGroupVersionKind = SchemeGroupVersion.WithKind(TokenKind)
)

func init() {
	SchemeBuilder.Register(&Token{}, &TokenList{})
}
