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
	"encoding/json"
	"fmt"
	"reflect"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	vllm "github.com/revving-ai/vLLM-k8s-operator/api/v1alpha1"
)

const controllerName = "vllmDeployment-controller"

// VllmDeploymentReconciler reconciles a VllmDeployment object
type VllmDeploymentReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=core.vllmoperator.org,resources=vllmdeployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.vllmoperator.org,resources=vllmdeployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.vllmoperator.org,resources=vllmdeployments/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the VllmDeployment object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.0/pkg/reconcile
func (r *VllmDeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := log.FromContext(ctx).WithValues("vllmdeployment", req.NamespacedName)
	log.Info("Starting Reconciliation")

	// Fetch the VllmDeployment instance
	var vllmDeployment vllm.VllmDeployment
	if err := r.Get(ctx, req.NamespacedName, &vllmDeployment); err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("VllmDeployment resource not found. Ignoring since object might be deleted")
			return ctrl.Result{}, err
		}
		// requeue the request
		log.Error(err, "Failed to get VllmDeployment")
		return ctrl.Result{}, err
	}
	// TODO: placeholder to handle deletion later
	if !vllmDeployment.ObjectMeta.DeletionTimestamp.IsZero() {
		log.Info("VllmDeployment is being deleted")
		log.Info(fmt.Sprintf("VllmDeployment is being deleted: Name=%s, Namespace=%s", req.NamespacedName.Name, req.NamespacedName.Namespace))
		return ctrl.Result{}, nil
	}
	desiredDeployment := constructDeployment(&vllmDeployment)

	jsonOutput, err := json.MarshalIndent(desiredDeployment, "", " ")
	if err != nil {
		log.Error(err, "Failed to convert to json")
	} else {
		fmt.Println("**************************")
		fmt.Println(string(jsonOutput))
	}

	// Set the owner reference
	if err := ctrl.SetControllerReference(&vllmDeployment, desiredDeployment, r.Scheme); err != nil {
		log.Error(err, "Failed to set owner reference on Deploynent")
		return ctrl.Result{}, nil
	}

	// checking if the deployment already exists

	var existingDeployment appsv1.Deployment
	err = r.Get(ctx, types.NamespacedName{
		Name:      desiredDeployment.Name,
		Namespace: desiredDeployment.Namespace,
	}, &existingDeployment)
	if err != nil && apierrors.IsNotFound(err) {
		// No deployment found, go ahead and create it.
		log.Info("Creating a new Deployment", "Deployment.Namespace", desiredDeployment.Namespace, "Deployment.Name", desiredDeployment.Name)
		if err := r.Create(ctx, desiredDeployment); err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", desiredDeployment.Namespace, "Deployment.Name", desiredDeployment.Name)
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		// Error reading the Deployment - requeue
		log.Error(err, "Failed to get deployment")
		return ctrl.Result{}, err
	}
	// Existing Deployment found. Check if there is a need to update it
	// Compare the desired and existing Deployment specs

	if !reflect.DeepEqual(existingDeployment.Spec, desiredDeployment.Spec) {
		log.Info("Updating existing deployment")
		// Create a copy of the existing Deployment to avoid modifying the cache
		updatedDep := existingDeployment.DeepCopy()
		updatedDep.Spec = desiredDeployment.Spec
		//Update the deployment
		if err := r.Update(ctx, updatedDep); err != nil {
			log.Error(err, "Failed to update the Deployment")
			return ctrl.Result{}, err
		}
		// Deployment updated successfully - requeue for status update
		return ctrl.Result{Requeue: true}, nil

	}

	// Update the status of VllmDeployment if necessary
	// Fetch the latest Deployment status
	if err := r.Get(ctx, types.NamespacedName{Name: existingDeployment.Name, Namespace: existingDeployment.Namespace}, &existingDeployment); err != nil {
		log.Error(err, "Failed to re-fetch Deployment for status update")
		return ctrl.Result{}, err
	}

	// Update the VllmDeployment status
	updatedStatus := vllmDeployment.Status.DeepCopy()
	//TODO: CHANGE SOON updatedStatus.Replicas = existingDeployment.Status.ReadyReplicas

	if !reflect.DeepEqual(vllmDeployment.Status, *updatedStatus) {
		// Update status
		vllmDeployment.Status = *updatedStatus
		if err := r.Status().Update(ctx, &vllmDeployment); err != nil {
			log.Error(err, "Failed to update VllmDeployment status")
			return ctrl.Result{}, err
		}
	}

	log.Info("Reconciliation complete")
	// Reconciliation successful - don't requeue
	return ctrl.Result{}, nil
}

// constructDeployment constructs a Deployment object based on the given vllmDeployment // by mapping the fields from the vllmDeployment spec to the Deployment spec
func constructDeployment(v *vllm.VllmDeployment) *appsv1.Deployment {

	labels := map[string]string{
		"app": v.Name,
	}
	if v.Labels != nil {
		for k, v := range v.Labels {
			labels[k] = v
		}
	}

	envVars := []corev1.EnvVar{}
	vllmContainer := getVllmContainer(&v.Spec)
	if vllmContainer != nil {
		envVars = vllmContainer.Env
	}

	args := convertVllmConfigToArgs(&v.Spec)

	// prepare container ports
	containerPorts := []corev1.ContainerPort{}
	if port := v.Spec.VLLMConfig.Port; port != 0 {
		containerPorts = append(containerPorts, corev1.ContainerPort{
			ContainerPort: int32(port),
			Protocol:      corev1.ProtocolTCP,
		})
	}

	// prepare container spec
	container := corev1.Container{
		Name:            vllmContainer.Name,
		Image:           vllmContainer.Image,
		ImagePullPolicy: corev1.PullPolicy(vllmContainer.ImagePullPolicy),
		Env:             envVars,
		Args:            args,
		Ports:           containerPorts,
		// TODO: Add remaining
	}
	tolerations := []corev1.Toleration{}
	if t := v.Spec.Tolerations; t != nil {
		tolerations = t
	}

	// create pod template spec
	podTemplate := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Labels: labels,
		},
		Spec: corev1.PodSpec{
			Containers:  []corev1.Container{container},
			Tolerations: tolerations,
			// TODO: add the remaining here.
		},
	}
	var replicas int32 = 1
	if v.Spec.Replicas != nil && *v.Spec.Replicas != 0 {
		replicas = *v.Spec.Replicas
	}

	//create the deployment spec
	deploymentSpec := appsv1.DeploymentSpec{
		Replicas: &replicas,
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: podTemplate,
	}
	// Create the deployment object
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-deployment", v.Name),
			Namespace: v.Namespace,
			Labels:    labels,
		},
		Spec: deploymentSpec,
	}

	return deployment

}

// Util function that will fetch vllm container where are there
// more than 1 container e.g. authentication proxy

func getVllmContainer(v *vllm.VllmDeploymentSpec) *corev1.Container {
	// If there is exactly one container, return it directly
	if len(v.Containers) == 1 {
		return &v.Containers[0]
	}
	// If there are multiple containers, look for the container with the name "vllm"
	for _, container := range v.Containers {
		if container.Name == "vllm" {
			return &container
		}
	}
	return nil
}

// Convert vllmconfig struct items into args that can be passed to the vllm container.

func convertVllmConfigToArgs(v *vllm.VllmDeploymentSpec) []string {
	vc := v.VLLMConfig
	args := []string{}

	if vc.GpuMemoryUtilization != "" {
		args = append(args, "--gpu-memory-utilization", vc.GpuMemoryUtilization)
	}
	if vc.LogLevel != "" {
		args = append(args, "--log-level", vc.LogLevel)
	}
	if vc.BlockSize != 0 {
		args = append(args, "--block-size", fmt.Sprintf("%d", vc.BlockSize))
	}
	if vc.MaxModelLen != 0 {
		args = append(args, "--max-model-len", fmt.Sprintf("%d", vc.MaxModelLen))

	}

	// Add port if specified
	if vc.Port != 0 {
		// Adding --port and converting the integer to a string
		args = append(args, "--port", fmt.Sprintf("%d", vc.Port))
	}

	return args
}

// SetupWithManager sets up the controller with the Manager.
func (r *VllmDeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.recorder = mgr.GetEventRecorderFor(controllerName)
	return ctrl.NewControllerManagedBy(mgr).
		For(&vllm.VllmDeployment{}).
		Owns(&appsv1.Deployment{}).
		Named(controllerName).
		Complete(r)
}
