package shardingredis

import (
	cachev1beta1 "github.com/lancelot1989/RedisOperator/pkg/apis/cache/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func buildStatefulSet(instance *cachev1beta1.ShardingRedis) *appsv1.StatefulSet {
	labelSelectorMap := buildLabelSelectorMap(instance)
	serviceName := buildServiceName(instance)
	statefulSetName := buildStatefulSetName(instance)
	statefulset := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetName,
			Namespace: instance.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labelSelectorMap,
			},
			Replicas:    &instance.Spec.Replicas,
			ServiceName: serviceName,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labelSelectorMap,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "redis",
							Image:           instance.Spec.Image,
							ImagePullPolicy: corev1.PullIfNotPresent,
							Ports: []corev1.ContainerPort{
								{
									Name:          "redis",
									ContainerPort: 6379,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							Command: []string{
								"redis-server",
							},
							ReadinessProbe: &corev1.Probe{
								InitialDelaySeconds: 30,
								FailureThreshold:    3,
								PeriodSeconds:       6,
								TimeoutSeconds:      3,
								Handler: corev1.Handler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"sh",
											"-c",
											"redis-cli -h $(hostname) ping",
										},
									},
								},
							},
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 20,
								TimeoutSeconds:      2,
								FailureThreshold:    3,
								PeriodSeconds:       3,
								Handler: corev1.Handler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(6379),
									},
								},
							},
							Resources: instance.Spec.Resources,
						},
					},
				},
			},
		},
	}
	return statefulset
}

func buildStatefulSetName(instance *cachev1beta1.ShardingRedis) string {
	statefulSetName := instance.Name + "-statefulset"
	return statefulSetName
}

func buildLabelSelectorMap(sr *cachev1beta1.ShardingRedis) map[string]string {
	return map[string]string{LABEL_SELECTOR_KEY: sr.Namespace + "-" + sr.Name}
}

func buildServiceName(sr *cachev1beta1.ShardingRedis) string {
	return "redis-" + sr.Name
}
