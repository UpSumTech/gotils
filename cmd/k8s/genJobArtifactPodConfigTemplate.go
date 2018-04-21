package k8s

import (
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

func GenArtifactBuilderPodTemplate() (string, error) {
	var data string

	if appPort == 0 {
		return data, utils.RaiseErr("Missing the port number")
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: strings.Join([]string{imageName, "deployment"}, "-"),
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: utils.Int32Ptr(2),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": imageName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": imageName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  imageName,
							Image: strings.Join([]string{imageName, imageTag}, ":"),
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: *utils.Int32Ptr(appPort),
								},
							},
						},
					},
				},
			},
		},
	}

	data, err := utils.ToJson(deployment)
	if err != nil {
		return data, err
	}
	return data, nil
}
