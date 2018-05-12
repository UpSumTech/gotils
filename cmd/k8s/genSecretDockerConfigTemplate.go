/*
cat input.json

{
  "docker_user": "sumanmukherjee03"
}
*/

package k8s

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/sumanmukherjee03/gotils/cmd/utils"
)

func genSecretDockerConfigTemplate(input DockerConfigSecretTemplate) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      DOCKER_CONFIG_SECRET_NAME,
			Namespace: namespace,
			Labels: map[string]string{
				"type": "dockercfg",
			},
			Annotations: map[string]string{
				"description": fmt.Sprintf("Contains the dockercfg file"),
			},
		},
		Data: map[string][]byte{
			corev1.DockerConfigKey: utils.GetDockerConfig(input.DockerConfigFile),
		},
		Type: corev1.SecretTypeDockercfg,
	}
}
