package deployment

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deployment struct {
	Name       string
	Namespace  string
	ServerType string
	SecretName string
	Port       int
}

func NewDeployment(name, namespace, serverType, secretName string, port int) *appsv1.Deployment {
	dep := Deployment{
		Name:       name,
		Namespace:  namespace,
		ServerType: serverType,
		SecretName: secretName,
		Port:       port,
	}

	return dep.newDeployment()

}

func (d *Deployment) newDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: d.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": d.Name,
				},
			},
			Template: d.getTemplate(),
		},
	}
}

func (d *Deployment) getTemplate() corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.Name,
			Labels: map[string]string{
				"app": d.Name,
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Image:           "hooksie1/cmsnr:latest",
					ImagePullPolicy: "Always",
					Name:            d.Name,
					Args:            []string{"server", "start", d.ServerType, fmt.Sprintf("-n=%s", d.Namespace)},
					Ports: []corev1.ContainerPort{
						{
							Name:          "https",
							ContainerPort: int32(d.Port),
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "webhook-certs",
							MountPath: "/var/lib/cmsnr",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "webhook-certs",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: d.SecretName,
						},
					},
				},
			},
		},
	}
}
