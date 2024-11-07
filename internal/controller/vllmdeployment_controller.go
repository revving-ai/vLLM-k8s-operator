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
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/revving-ai/vLLM-k8s-operator/api/v1alpha1"
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
	_ = log.FromContext(ctx)

	// TODO(user): your logic here

	return ctrl.Result{}, nil
}

// constructDeployment constructs a Deployment object based on the given vllmDeployment // by mapping the fields from the vllmDeployment spec to the Deployment spec
func constructDeployment(v *v1alpha1.VllmDeployment) *appsv1.Deployment {

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
		envVars = c.Env
	}

	args := convertVllmConfigToArgs(&v.Spec)

	// prepare container ports
	containerPorts := []corev1.ContainerPort{}
	if port := v.Spec.VLLMConfig.Port; port != 0 {
		containerPorts := append(containerPorts, corev1.ContainerPort{
			ContainerPort: int32(port),
			Protocol:      corev1.ProtocolTCP,
		})
	}

	// prepare container spec
	container := corev1.Container{
		Name:            vllmContainer.Name,
		Image:           vllmContainer.Image,
		ImagePullPolicy: v1.PullPolicy(vllmContainer.ImagePullPolicy),
		Env:             envVars,
		Args:            args,
		Ports:           containerPorts,
	}

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
		args = append(args, fmt.Sprintf("--gpu-memory-utilization=%s", vc.GpuMemoryUtilization))
	}
	if vc.LogLevel != "" {
		args = append(args, fmt.Sprintf("--log-level=%s", vc.LogLevel))

	}
	if vc.BlockSize != 0 {
		args = append(args, fmt.Sprintf("--block-size=%d", vc.BlockSize))
	}
	if vc.MaxModelLen != 0 {
		args = append(args, fmt.Sprintf("--max-model-len=%d", vc.MaxModelLen))
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
