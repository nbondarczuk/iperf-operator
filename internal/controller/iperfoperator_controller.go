/*
Copyright 2025.

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

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorv1alpha1 "gitgub.com/nbondarczuk/iperf-operator/api/v1alpha1"
	"gitgub.com/nbondarczuk/iperf-operator/assets"
	appsv1 "k8s.io/api/apps/v1"
)

// IPerfOperatorReconciler reconciles a IPerfOperator object
type IPerfOperatorReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=operator.nbt.pl,resources=iperfoperators,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=operator.nbt.pl,resources=iperfoperators/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=operator.nbt.pl,resources=iperfoperators/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the IPerfOperator object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *IPerfOperatorReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Get operator CR and skip all if not found.
	operatorCR := &operatorv1alpha1.IPerfOperator{}
	err := r.Get(ctx, req.NamespacedName, operatorCR)
	if err != nil && errors.IsNotFound(err) {
		logger.Info("Operator CR object not found.")
		return ctrl.Result{}, nil
	} else if err != nil {
		logger.Error(err, "Error getting operator CR object")
		return ctrl.Result{}, err
	}

	// Get server deplyment from manifest file.
	serverDeployment := &appsv1.Deployment{}
	serverCreate := false
	err = r.Get(ctx, req.NamespacedName, serverDeployment)
	if err != nil && errors.IsNotFound(err) {
		serverCreate = true
		serverDeployment = assets.GetDeploymentFromFile("assets/manifests/iperf3-server-deployment.yaml")
	} else if err != nil {
		logger.Error(err, "Error getting existing IPerf server deployment manifest.")
		return ctrl.Result{}, err
	}

	// Set properties of server controller
	serverDeployment.Namespace = req.Namespace
	serverDeployment.Name = req.Name
	serverDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = operatorCR.Spec.Port

	// Create or Update server controller
	ctrl.SetControllerReference(operatorCR, serverDeployment, r.Scheme)
	if serverCreate {
		err = r.Create(ctx, serverDeployment)
	} else {
		err = r.Update(ctx, serverDeployment)
	}

   // Get client deplyment from manifest file.
    clientDeployment := &appsv1.Deployment{}
    clientCreate := false
    err = r.Get(ctx, req.NamespacedName, clientDeployment)
    if err != nil && errors.IsNotFound(err) {
        clientCreate = true
        clientDeployment = assets.GetDeploymentFromFile("assets/manifests/iperf3-client-deployment.yaml")
    } else if err != nil {
        logger.Error(err, "Error getting existing IPerf client deployment manifest.")
        return ctrl.Result{}, err
    }

    // Set properties of server controller
    clientDeployment.Namespace = req.Namespace
	clientDeployment.Spec.Template.Spec.Containers[0].Ports[0].ContainerPort = operatorCR.Spec.Port

    // Create or Update server controller
    if clientCreate {
        err = r.Create(ctx, clientDeployment)
    } else {
        err = r.Update(ctx, clientDeployment)
    }

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *IPerfOperatorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorv1alpha1.IPerfOperator{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
