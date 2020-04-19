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

	dbs "sampledb/api/v1alpha1"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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
	ctx := context.Background()
	log := r.Log.WithValues("statefuldb", req.NamespacedName)

	log.Info("Starting recon")
	// Get the StatefulDB resource
	statefulDB := &dbs.StatefulDB{}
	if err := r.Get(ctx, req.NamespacedName, statefulDB); err != nil {
		log.Error(err, "unable to fetch statefulDB")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Getting Desired StatefulSet resource")
	// Get the desired StatefulSet
	desiredSS, err := r.desiredStatefulset(*statefulDB)
	if err != nil {
		log.Error(err, "unable to get the desiredStatefulset")
		return ctrl.Result{}, err
	}

	log.Info("Server side apply of this spec")
	// Patch the statefulSet
	applyOpts := []client.PatchOption{client.ForceOwnership, client.FieldOwner("statefuldb-controller")}
	err = r.Patch(ctx, &desiredSS, client.Apply, applyOpts...)
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("Update the statefulDB status with the StatefulSet status")
	// Update Status
	statefulDB.Status.SSStatus = desiredSS.Status

	err = r.Status().Update(ctx, statefulDB)
	if err != nil {
		log.Info("Error while updating status", "Error", err)
		return ctrl.Result{}, err
	}

	log.Info("Recon successful!")

	return ctrl.Result{}, nil
}

func (r *StatefulDBReconciler) desiredStatefulset(statefulDB dbs.StatefulDB) (appsv1.StatefulSet, error) {
	ss := appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: appsv1.SchemeGroupVersion.String(),
			Kind:       "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulDB.Name + "-ss",
			Namespace: statefulDB.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: statefulDB.Spec.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"statefuldb": statefulDB.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"statefuldb": statefulDB.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    "busybox",
							Image:   "k8s.gcr.io/busybox",
							Command: []string{"/bin/sh", "-c", "sleep 10000"},
						},
					},
				},
			},
		},
	}
	if err := ctrl.SetControllerReference(&statefulDB, &ss, r.Scheme); err != nil {
		return ss, err
	}

	return ss, nil
}

func (r *StatefulDBReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&dbs.StatefulDB{}).
		Owns(&appsv1.StatefulSet{}).
		Complete(r)
}
