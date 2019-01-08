package shardingredis

import (
	"context"
	"github.com/lancelot1989/RedisOperator/pkg/apis/cache/v1beta1"
	"k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	CONDITION_TYPE_STATEFULSET           v1beta1.ConditionType = "StatefulSet"
	CONDITION_STATEFUL_SET_REASON_UPDATE                       = "Update"
	CONDITION_STATEFUL_SET_REASON_CREATE                       = "Create"
	CONDITION_STATEFUL_SET_REASON_READY                        = "Ready"
)

func ensureStatefulSet(r *ReconcileShardingRedis, instance *v1beta1.ShardingRedis) (statefulSetCondition *v1beta1.Condition, err error) {
	statefulset := buildStatefulSet(instance)

	if err := controllerutil.SetControllerReference(instance, statefulset, r.scheme); err != nil {
		return nil, err
	}

	found := &v1.StatefulSet{}

	err = r.Get(context.TODO(), types.NamespacedName{Name: statefulset.Name, Namespace: statefulset.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating StatefulSet", "namespace", statefulset.Namespace, "name", statefulset)
		err = r.Create(context.TODO(), statefulset)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(statefulset.Spec, found.Spec) {
		found.Spec = statefulset.Spec
		log.Info("Updating StatefulSet", "namespace", statefulset.Namespace, "name", statefulset.Name)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return nil, err
		}
		return &v1beta1.Condition{
			Type:   CONDITION_TYPE_STATEFULSET,
			Reason: CONDITION_STATEFUL_SET_REASON_UPDATE,
		}, nil
	} else {
		setStatus := found.Status
		if setStatus.ReadyReplicas != setStatus.Replicas {
			return &v1beta1.Condition{
				Type:   CONDITION_TYPE_STATEFULSET,
				Reason: CONDITION_STATEFUL_SET_REASON_UPDATE,
			}, nil
		}
	}
	return &v1beta1.Condition{
		Type:   CONDITION_TYPE_STATEFULSET,
		Reason: CONDITION_STATEFUL_SET_REASON_READY,
	}, nil
}
