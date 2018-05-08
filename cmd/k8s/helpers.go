package k8s

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/sumanmukherjee03/gotils/cmd/utils"
	validator "gopkg.in/go-playground/validator.v9"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func createDeployment(deployment *appsv1.Deployment) error {
	clientset := utils.GetK8sClientSet()
	deploymentsClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)
	_, err := deploymentsClient.Create(deployment)
	if err != nil {
		return err
	}
	return nil
}

func updateDeployment(deployment *appsv1.Deployment) error {
	clientset := utils.GetK8sClientSet()
	deploymentsClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)
	// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		result, getErr := deploymentsClient.Get("demo-deployment", metav1.GetOptions{})
		if getErr != nil {
			return fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
		}
		result.Spec.Replicas = utils.Int32Ptr(1)                     // reduce replica count
		result.Spec.Template.Spec.Containers[0].Image = "nginx:1.13" // change nginx version
		_, updateErr := deploymentsClient.Update(result)
		return updateErr
	})
	if retryErr != nil {
		return fmt.Errorf("Update failed: %v", retryErr)
	}
	return nil
}

func listDeployments() {
	clientset := utils.GetK8sClientSet()
	deploymentsClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)
	list, err := deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (%d replicas)\n", d.Name, *d.Spec.Replicas)
	}
}

func deleteDeployment(deployment *appsv1.Deployment) error {
	clientset := utils.GetK8sClientSet()
	deploymentsClient := clientset.AppsV1().Deployments(corev1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentsClient.Delete("demo-deployment", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}
	return nil
}

func getResourceRequirements(l ResourceLimitConfig, r ResourceRequestConfig) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(l.Cpu),
			corev1.ResourceMemory: resource.MustParse(l.Memory),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(r.Cpu),
			corev1.ResourceMemory: resource.MustParse(r.Memory),
		},
	}
}

func readJson(i JsonInput, src string) error {
	r, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	json.Unmarshal(r, i)
	return nil
}

func validateJsonInput(i JsonInput) error {
	validate := validator.New()
	err := validate.Struct(i)
	if err != nil {
		return err
	}
	return nil
}
