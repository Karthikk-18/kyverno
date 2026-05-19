package policy

import (
    "testing"
 
	kyvernov2 "github.com/kyverno/kyverno/api/kyverno/v2"
	"github.com/kyverno/kyverno/pkg/config"
)

func Test_newUrMeta(t *testing.T) {
	result := newUrMeta()
    
	wantKind := "UpdateRequest"
	if result.Kind != wantKind {
		t.Errorf("newUrMeta() kind: %q, want: %q", result.Kind, wantKind)
	}
	
	wantGenerateName := "ur-"
    if result.GenerateName != wantGenerateName {
		t.Errorf("newUrMeta() GenerateName = %q, want %q", result.GenerateName, wantGenerateName)
	}

	wantNamespace := config.KyvernoNamespace()
	if result.Namespace != wantNamespace {
		t.Errorf("newUrMeta() Namespace = %q, want %q", result.Namespace, wantNamespace)
	}

	wantAPIVersion := kyvernov2.SchemeGroupVersion.String()
	if result.APIVersion != wantAPIVersion {
		t.Errorf("newUrMeta() APIVersion = %q, want %q", result.APIVersion, wantAPIVersion)
	}
}