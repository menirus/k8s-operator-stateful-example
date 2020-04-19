/*


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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	databasesv1alpha1 "sampledb/api/v1alpha1"
)

// StatefulDBReconciler reconciles a StatefulDB object
type StatefulDBReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=databases.menirus95.dev,resources=statefuldbs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=databases.menirus95.dev,resources=statefuldbs/status,verbs=get;update;patch

func (r *StatefulDBReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("statefuldb", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *StatefulDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&databasesv1alpha1.StatefulDB{}).
		Complete(r)
}
