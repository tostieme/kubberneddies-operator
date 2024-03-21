/*
Copyright 2024.

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
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiv1 "github.com/tostieme/kubberneddies-operator/api/v1"
)

// AutoUnsealerReconciler reconciles a AutoUnsealer object
type AutoUnsealerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.tostieme.me,resources=autounsealers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.tostieme.me,resources=autounsealers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.tostieme.me,resources=autounsealers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AutoUnsealer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *AutoUnsealerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// TODO(user): your logic here
	log.Info("Request Namespace: ", req.Namespace)
	log.Info("Request Name: ", req.Name)
	// l := log.FromContext(ctx)
	log.Info("Reconcile called")
	unsealer := &apiv1.AutoUnsealer{}
	if err := r.Get(ctx, req.NamespacedName, unsealer); err != nil {
		return ctrl.Result{}, nil
	}
	adress := unsealer.Spec.Adress //   "/v1/sys/seal-status"

	response, err := http.Get(adress + "/v1/sys/seal-status")
	if err != nil {
		log.Error("Error:", err)
		return ctrl.Result{}, nil
	}
	defer response.Body.Close()

	// Read response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Error("Error reading response body:", err)
		return ctrl.Result{}, nil
	}

	log.Info(string(body))
	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoUnsealerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.AutoUnsealer{}).
		Complete(r)
}
