<<<<<<< HEAD
package policy

import (
	"testing"

	kyvernov1 "github.com/kyverno/kyverno/api/kyverno/v1"
	kyvernov2 "github.com/kyverno/kyverno/api/kyverno/v2"
	common "github.com/kyverno/kyverno/pkg/background/common"
	"github.com/kyverno/kyverno/pkg/config"
)

func Test_newMutateUR(t *testing.T) {
	tests := []struct {
		name     string
		policy   kyvernov1.PolicyInterface
		trigger  kyvernov1.ResourceSpec
		ruleName string
	}{
		{
			name:     "Successfully creates a mutate UpdateRequest",
			policy:   makeClusterPolicy("test-policy", nil),
			ruleName: "check-pod-labels",
			trigger: kyvernov1.ResourceSpec{
				Kind:       "Pod",
				Namespace:  "default",
				Name:       "test-pod",
				APIVersion: "v1",
				UID:        "abc-123",
			},
		},
		{
			name:     "empty trigger fields enforcing policy without panicking",
			policy:   makeClusterPolicy("test-policy", nil),
			ruleName: "check-empty-fields",
			trigger:  kyvernov1.ResourceSpec{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := newMutateUR(tt.policy, tt.trigger, tt.ruleName)

			if result.Spec.Type != kyvernov2.Mutate {
				t.Errorf("wrong type. want: mutate, got: %v", result.Spec.Type)
			}
			if result.Spec.Rule != tt.ruleName {
				t.Errorf("wrong rule name. want: %q, got: %q", tt.ruleName, result.Spec.Rule)
			}
			if result.Spec.Policy != policyKey(tt.policy) {
				t.Errorf("Spec.Policy: %q, want: %q", result.Spec.Policy, policyKey(tt.policy))
			}

			// labels
			wantLabels := common.MutateLabelsSet(policyKey(tt.policy), tt.trigger)
			for k, wantVal := range wantLabels {
				if gotVal, ok := result.Labels[k]; !ok {
					t.Errorf("Labels missing key %q", k)
				} else if gotVal != wantVal {
					t.Errorf("Labels[%q]: %q, want: %q", k, gotVal, wantVal)
				}
			}

			// every field from trigger
			res := result.Spec.Resource
			if res.Kind != tt.trigger.GetKind() {
				t.Errorf("Spec.Resource.Kind: %q, want: %q", res.Kind, tt.trigger.GetKind())
			}
			if res.Namespace != tt.trigger.GetNamespace() {
				t.Errorf("Spec.Resource.Namespace: %q, want: %q", res.Namespace, tt.trigger.GetNamespace())
			}
			if res.Name != tt.trigger.GetName() {
				t.Errorf("Spec.Resource.Name: %q, want: %q", res.Name, tt.trigger.GetName())
			}
			if res.APIVersion != tt.trigger.GetAPIVersion() {
				t.Errorf("Spec.Resource.APIVersion: %q, want: %q", res.APIVersion, tt.trigger.GetAPIVersion())
			}
			if res.UID != tt.trigger.GetUID() {
				t.Errorf("Spec.Resource.UID: %q, want: %q", res.UID, tt.trigger.GetUID())
			}
		})
	}
}

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