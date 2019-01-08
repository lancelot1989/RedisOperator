/*
Copyright 2019 Thomas Liang.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package shardingredis

import (
	"context"
	cachev1beta1 "github.com/lancelot1989/RedisOperator/pkg/apis/cache/v1beta1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

const (
	LABEL_SELECTOR_KEY string = "sr-label"
)

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new ShardingRedis Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileShardingRedis{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("shardingredis-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to ShardingRedis
	err = c.Watch(&source.Kind{Type: &cachev1beta1.ShardingRedis{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// watch a StatefulSet created by ShardingRedis - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.StatefulSet{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1beta1.ShardingRedis{},
	})
	if err != nil {
		return err
	}

	// watch a Service created by ShardingRedis - change this for objects you create
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &cachev1beta1.ShardingRedis{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileShardingRedis{}

// ReconcileShardingRedis reconciles a ShardingRedis object
type ReconcileShardingRedis struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a ShardingRedis object and makes changes based on the state read
// and what is in the ShardingRedis.Spec
// Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps,resources=deployments/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=cache.thomas.com,resources=shardingredis,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=cache.thomas.com,resources=shardingredis/status,verbs=get;update;patch
func (r *ReconcileShardingRedis) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the ShardingRedis instance
	instance := &cachev1beta1.ShardingRedis{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	serviceCondition, err := ensureShardingRedisService(r, instance)

	if err != nil {
		return reconcile.Result{}, err
	}

	ssCondition, err := ensureStatefulSet(r, instance)

	if err != nil {
		return reconcile.Result{}, err
	}

	var phase cachev1beta1.Phase = "PREPARING"
	if ssCondition.Reason == CONDITION_STATEFUL_SET_REASON_READY {
		phase = "READY"
	}

	conditions := make([]cachev1beta1.Condition, 0)

	conditions = append(conditions, *serviceCondition, *ssCondition)

	status := cachev1beta1.ShardingRedisStatus{
		Phase:      phase,
		Conditions: conditions,
	}

	instance.Status = status

	r.Status().Update(context.TODO(), instance)
	return reconcile.Result{}, nil
}
