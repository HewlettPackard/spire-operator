/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"os/exec"

	"gopkg.in/yaml.v2"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	spirev1 "github.com/glcp/spire-k8s-operator/api/v1"
)

// SpireServerReconciler reconciles a SpireServer object
type SpireServerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

type TruatBundelMetadata struct {
	Name      string `yaml:"name"`
	Namespace string `yaml:"namespace"`
}

type TrustBundle struct {
	APIVersion string              `yaml:"apiVersion"`
	Kind       string              `yaml:"kind"`
	Metadata   TruatBundelMetadata `yaml:"metadata"`
}

//+kubebuilder:rbac:groups=spire.hpe.com,resources=spireservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=spire.hpe.com,resources=spireservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=spire.hpe.com,resources=spireservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the SpireServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *SpireServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	createSpireBundle()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SpireServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&spirev1.SpireServer{}).
		Complete(r)
}

func createSpireBundle() {

	data := TruatBundelMetadata{
		Name:      "spire-bundle",
		Namespace: "sprie",
	}
	trustBundle := TrustBundle{
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Metadata:   data,
	}
	yamlData, err := yaml.Marshal(&trustBundle)
	if err != nil {
		fmt.Printf("Error while Marshaling. %v", err)
	}
	exec.Command(fmt.Sprintf("kubectl apply -f <<EOF %s EOF", string(yamlData)), "./")
}
