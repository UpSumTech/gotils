/*
cat input.json

{
  "github_user": "sumanmukherjee03"
}
*/

package k8s

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// Uncomment the following line to load the gcp plugin (only required to authenticate against GKE clusters).
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

func genSecretGithubTokenTemplate(input GithubTokenSecretTemplate) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      GITHUB_TOKEN_SECRET_NAME,
			Namespace: namespace,
			Labels: map[string]string{
				"type": "githubtoken",
			},
			Annotations: map[string]string{
				"description": fmt.Sprintf("Contains the github token"),
			},
		},
		Data: map[string][]byte{
			GITHUB_TOKEN_SECRET_KEY: []byte(input.GithubToken),
		},
		StringData: map[string]string{
			GITHUB_USERNAME_SECRET_KEY: input.GithubUser,
		},
		Type: corev1.SecretTypeOpaque,
	}
}
