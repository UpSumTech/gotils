package k8s

import (
	"fmt"

	"github.com/sumanmukherjee03/gotils/cmd/utils"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func createDeployment(deployment *appsv1.Deployment) error {
	clientset := utils.GetK8sClientSet()
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	_, err := deploymentsClient.Create(deployment)
	if err != nil {
		return err
	}
	return nil
}

func updateDeployment(deployment *appsv1.Deployment) error {
	clientset := utils.GetK8sClientSet()
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
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
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
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
	deploymentsClient := clientset.AppsV1().Deployments(apiv1.NamespaceDefault)
	deletePolicy := metav1.DeletePropagationForeground
	err := deploymentsClient.Delete("demo-deployment", &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
	if err != nil {
		return err
	}
	return nil
}
