/*
cat input.json

{
  "limits": {
    "cpu": "100m",
    "memory": "100Mi"
  },
  "termination_grace_period": 60,
  "deadline": 300,
  "docker_user": "sumanmukherjee03",
}
*/

package k8s

import (
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func genServiceTemplate(input ServiceTemplate) *corev1.Service {
	serviceName := strings.Join([]string{"service", imageName, imageTag}, "-")
	runLabel := strings.Join([]string{imageName, imageTag}, "-")
	return &corev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
			Labels: map[string]string{
				LABEL_RUN_KEY: runLabel,
				LABEL_APP_KEY: imageName,
			},
		},
		Spec: corev1.ServiceSpec{
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Protocol:   corev1.ProtocolTCP,
					Port:       int32(appPort),
					TargetPort: intstr.FromInt(appPort),
				},
			},
			Selector: map[string]string{
				LABEL_RUN_KEY: runLabel,
				LABEL_APP_KEY: imageName,
			},
			ClusterIP:             CLUSTER_IP_DEFAULT,
			Type:                  corev1.ServiceTypeLoadBalancer,
			SessionAffinity:       corev1.ServiceAffinityNone,
			ExternalTrafficPolicy: corev1.ServiceExternalTrafficPolicyTypeCluster,
		},
	}
}
