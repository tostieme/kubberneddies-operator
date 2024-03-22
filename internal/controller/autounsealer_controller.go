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
	"bytes"
	"context"
	"encoding/json"

	// "io"
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
	body := map[string]any{}
	json.NewDecoder(response.Body).Decode(&body)
	if body["sealed"] == true {
		log.Info("Shaise dat Ding is sealed")
		// Hier müsse mer reparieren dann
		// http://127.0.0.1:8200/v1/sys/unseal
		// Get unseal treshold
		unsealTreshold := body["t"]
		for i := 0; i < unsealTreshold.(int); i++ {
            // get Unseal key
			data := []byte(`{"key":"abc"`)
			http.Post(adress+"/v1/sys/unseal", "application/json", bytes.NewBuffer(data))
		}
	}

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

func getUnsealKey() []byte{

    return []byte(`{"key":"abc"`)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AutoUnsealerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1.AutoUnsealer{}).
		Complete(r)
}
