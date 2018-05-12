/*
cat input.json

{
  "bintray_user": "sumanmukherjee03"
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

func genSecretBintrayTokenTemplate(input BintrayTokenSecretTemplate) *corev1.Secret {
	return &corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      BINTRAY_TOKEN_SECRET_NAME,
			Namespace: namespace,
			Labels: map[string]string{
				"type": "bintraytoken",
			},
			Annotations: map[string]string{
				"description": fmt.Sprintf("Contains the bintray token"),
			},
		},
		Data: map[string][]byte{
			BINTRAY_TOKEN_SECRET_KEY: []byte(input.BintrayToken),
		},
		StringData: map[string]string{
			BINTRAY_USERNAME_SECRET_KEY:  input.BintrayUser,
			BINTRAY_REPO_NAME_SECRET_KEY: input.BintrayRepo,
		},
		Type: corev1.SecretTypeOpaque,
	}
}
