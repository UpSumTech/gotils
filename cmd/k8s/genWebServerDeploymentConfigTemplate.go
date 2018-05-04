package k8s

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

func genWebServerDeploymentTemplate(input WebServerDeployment) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.Join([]string{imageName, "deployment"}, "-"),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"name": imageName,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"name": imageName,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  imageName,
							Image: strings.Join([]string{imageName, imageTag}, ":"),
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: *utils.Int32Ptr(appPort),
								},
							},
						},
					},
				},
			},
		},
	}
}
