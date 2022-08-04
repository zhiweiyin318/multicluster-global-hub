/*
Copyright 2022.

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

package hubofhubs

import (
	"context"
	"embed"

	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/event"
	ctrllog "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	"github.com/go-logr/logr"
	hubofhubsv1alpha1 "github.com/stolostron/hub-of-hubs/operator/apis/hubofhubs/v1alpha1"
	"github.com/stolostron/hub-of-hubs/operator/pkg/config"
	"github.com/stolostron/hub-of-hubs/operator/pkg/constants"
	leafhubscontroller "github.com/stolostron/hub-of-hubs/operator/pkg/controllers/leafhub"

	// pmcontroller "github.com/stolostron/hub-of-hubs/operator/pkg/controllers/packagemanifest"
	"github.com/stolostron/hub-of-hubs/operator/pkg/deployer"
	"github.com/stolostron/hub-of-hubs/operator/pkg/renderer"
	"github.com/stolostron/hub-of-hubs/operator/pkg/utils"
)

//go:embed manifests
var fs embed.FS

var isLeafHubControllerRunnning = false

const NAMESPACE = "open-cluster-management"

// var isPackageManifestControllerRunnning = false

// ConfigReconciler reconciles a Config object
type ConfigReconciler struct {
	manager.Manager
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=hubofhubs.open-cluster-management.io,resources=configs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=hubofhubs.open-cluster-management.io,resources=configs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=hubofhubs.open-cluster-management.io,resources=configs/finalizers,verbs=update

//+kubebuilder:rbac:groups="",resources=services,verbs=get;create;update;delete
//+kubebuilder:rbac:groups="",resources=serviceaccounts,verbs=get;create;update;delete
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;create;update;delete
//+kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;create;update;delete
//+kubebuilder:rbac:groups="apps",resources=deployments,verbs=get;create;update;delete
//+kubebuilder:rbac:groups="networking.k8s.io",resources=ingresses,verbs=get;create;update;delete
//+kubebuilder:rbac:groups="networking.k8s.io",resources=networkpolicies,verbs=get;create;update;delete
//+kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=roles,verbs=get;create;update;delete
//+kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=rolebindings,verbs=get;create;update;delete
//+kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterroles,verbs=get;list;create;update;delete
//+kubebuilder:rbac:groups="rbac.authorization.k8s.io",resources=clusterrolebindings,verbs=get;list;create;update;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Config object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.0/pkg/reconcile
func (r *ConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrllog.FromContext(ctx)

	// Fetch the hub-of-hubs config instance
	hohConfig := &hubofhubsv1alpha1.Config{}
	err := r.Get(ctx, req.NamespacedName, hohConfig)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("Config resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Config")
		return ctrl.Result{}, err
	}

	err = r.setConditionResourceFound(ctx, hohConfig)
	if err != nil {
		log.Error(err, "Failed to set condition resource found")
		return ctrl.Result{}, err
	}

	// try to start leafhub controller if it is not running
	if !isLeafHubControllerRunnning {
		if err := (&leafhubscontroller.LeafHubReconciler{
			Client: r.Client,
			Scheme: r.Scheme,
		}).SetupWithManager(r.Manager); err != nil {
			log.Error(err, "unable to create controller", "controller", "LeafHub")
			return ctrl.Result{}, err
		}
		log.Info("leafhub controller is started")
		isLeafHubControllerRunnning = true
	}

	// // try to start packagemanifest controller if it is not running
	// if !isPackageManifestControllerRunnning {
	// 	if err := (&pmcontroller.PackageManifestReconciler{
	// 		Client: r.Client,
	// 		Scheme: r.Scheme,
	// 	}).SetupWithManager(r.Manager); err != nil {
	// 		log.Error(err, "unable to create controller", "controller", "PackageManifest")
	// 		return ctrl.Result{}, err
	// 	}
	// 	log.Info("packagemanifest controller is started")
	// 	isPackageManifestControllerRunnning = true
	// }

	// handle gc
	isTerminating, err := r.initFinalization(ctx, hohConfig, log)
	if err != nil {
		return ctrl.Result{}, err
	}

	if isTerminating {
		log.Info("hoh config is terminating, skip the reconcile")
		return ctrl.Result{}, err
	}

	// init DB and transport here
	err = r.reconcileDatabase(ctx, hohConfig, types.NamespacedName{Name: hohConfig.Spec.PostgreSQL.Name, Namespace: NAMESPACE})
	if err != nil {
		return ctrl.Result{}, err
	}

	// create new HoHRenderer and HoHDeployer
	hohRenderer := renderer.NewHoHRenderer(fs)
	hohDeployer := deployer.NewHoHDeployer(r.Client)

	annotations := hohConfig.GetAnnotations()
	hohRBACObjects, err := hohRenderer.Render("manifests/rbac", func(component string) (interface{}, error) {
		hohRBACConfig := struct {
			Image string
		}{
			Image: config.GetImage(annotations, "hub_of_hubs_rbac"),
		}

		return hohRBACConfig, err
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, obj := range hohRBACObjects {
		// TODO: more solid way to check if the object is namespace scoped resource
		if obj.GetNamespace() != "" {
			// set ownerreference of controller
			if err := controllerutil.SetControllerReference(hohConfig, obj, r.Scheme); err != nil {
				log.Error(err, "failed to set controller reference", "kind", obj.GetKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
			}
		}
		// set owner labels
		labels := obj.GetLabels()
		labels[constants.HoHOperatorOwnerLabelKey] = hohConfig.GetName()
		obj.SetLabels(labels)

		log.Info("Creating or updating object", "object", obj)
		err := hohDeployer.Deploy(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	managerObjects, err := hohRenderer.Render("manifests/manager", func(component string) (interface{}, error) {
		managerConfig := struct {
			Image string
		}{
			Image: config.GetImage(annotations, "hub_of_hubs_manager"),
		}

		return managerConfig, err
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	for _, obj := range managerObjects {
		// TODO: more solid way to check if the object is namespace scoped resource
		if obj.GetNamespace() != "" {
			// set ownerreference of controller
			if err := controllerutil.SetControllerReference(hohConfig, obj, r.Scheme); err != nil {
				log.Error(err, "failed to set controller reference", "kind", obj.GetKind(), "namespace", obj.GetNamespace(), "name", obj.GetName())
			}
		}
		// set owner labels
		labels := obj.GetLabels()
		labels[constants.HoHOperatorOwnerLabelKey] = hohConfig.GetName()
		obj.SetLabels(labels)

		log.Info("Creating or updating object", "object", obj)
		err := hohDeployer.Deploy(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *ConfigReconciler) initFinalization(ctx context.Context, hohConfig *hubofhubsv1alpha1.Config, log logr.Logger) (bool, error) {
	if hohConfig.GetDeletionTimestamp() != nil && utils.Contains(hohConfig.GetFinalizers(), constants.HoHOperatorFinalizer) {
		log.Info("to delete hoh resources")
		// clean up the cluster resources, eg. clusterrole, clusterrolebinding, etc
		if err := r.pruneGlobalResources(ctx, hohConfig); err != nil {
			log.Error(err, "failed to remove cluster scoped resources")
			return false, err
		}

		hohConfig.SetFinalizers(utils.Remove(hohConfig.GetFinalizers(), constants.HoHOperatorFinalizer))
		err := r.Client.Update(context.TODO(), hohConfig)
		if err != nil {
			log.Error(err, "failed to remove finalizer from hoh config resource")
			return false, err
		}
		log.Info("finalizer is removed from hoh config resource")

		return true, nil
	}
	if !utils.Contains(hohConfig.GetFinalizers(), constants.HoHOperatorFinalizer) {
		hohConfig.SetFinalizers(append(hohConfig.GetFinalizers(), constants.HoHOperatorFinalizer))
		err := r.Client.Update(context.TODO(), hohConfig)
		if err != nil {
			log.Error(err, "failed to add finalizer to hoh config resource")
			return false, err
		}
		log.Info("finalizer is added to hoh config resource")
	}

	return false, nil
}

// pruneGlobalResources deletes the cluster scoped resources created by the hub-of-hubs-operator
// cluster scoped resources need to be deleted manually because they don't have ownerrefenence set
func (r *ConfigReconciler) pruneGlobalResources(ctx context.Context, hohConfig *hubofhubsv1alpha1.Config) error {
	listOpts := []client.ListOption{
		client.MatchingLabels(map[string]string{constants.HoHOperatorOwnerLabelKey: hohConfig.GetName()}),
	}

	clusterRoleList := &rbacv1.ClusterRoleList{}
	err := r.Client.List(ctx, clusterRoleList, listOpts...)
	if err != nil {
		return err
	}
	for idx := range clusterRoleList.Items {
		err := r.Client.Delete(ctx, &clusterRoleList.Items[idx], &client.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	clusterRoleBindingList := &rbacv1.ClusterRoleBindingList{}
	err = r.Client.List(ctx, clusterRoleBindingList, listOpts...)
	if err != nil {
		return err
	}
	for idx := range clusterRoleBindingList.Items {
		err := r.Client.Delete(ctx, &clusterRoleBindingList.Items[idx], &client.DeleteOptions{})
		if err != nil {
			return err
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	hohConfigPred := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			// set request name to be used in leafhub controller
			config.SetHoHConfigNamespacedName(types.NamespacedName{Namespace: e.Object.GetNamespace(), Name: e.Object.GetName()})
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			return e.ObjectOld.GetResourceVersion() != e.ObjectNew.GetResourceVersion()
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return !e.DeleteStateUnknown
		},
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&hubofhubsv1alpha1.Config{}, builder.WithPredicates(hohConfigPred)).
		Complete(r)
}