package shardingredis

import (
	"context"
	"github.com/lancelot1989/RedisOperator/pkg/apis/cache/v1beta1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	CONDITION_TYPE_SERVICE          v1beta1.ConditionType = "Service"
	CONDITION_SERVICE_REASON_UPDATE                       = "Update"
	CONDITION_SERVICE_REASON_CREATE                       = "Create"
)

func ensureShardingRedisService(r *ReconcileShardingRedis, instance *v1beta1.ShardingRedis) (condition *v1beta1.Condition, err error) {
	service := buildService(instance)
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return nil, err
	}

	found := &v1.Service{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating Service", "namespace", service.Namespace, "name", service)
		err = r.Create(context.TODO(), service)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(service.Spec, found.Spec) {
		found.Spec = service.Spec
		log.Info("Updating Service", "namespace", service.Namespace, "name", service.Name)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return nil, err
		}
		return &v1beta1.Condition{
			Type:   CONDITION_TYPE_SERVICE,
			Reason: CONDITION_SERVICE_REASON_UPDATE,
		}, nil
	}
	return &v1beta1.Condition{
		Type:   CONDITION_TYPE_SERVICE,
		Reason: CONDITION_SERVICE_REASON_CREATE,
	}, nil
}

func buildService(instance *v1beta1.ShardingRedis) *v1.Service {
	return &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: instance.Namespace,
			Name:      buildServiceName(instance),
		},
		Spec: v1.ServiceSpec{
			Selector:  buildLabelSelectorMap(instance),
			ClusterIP: v1.ClusterIPNone,
		},
	}
}
