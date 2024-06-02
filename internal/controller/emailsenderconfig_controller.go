package controller

import (
	"context"
	"github.com/go-logr/logr"
	emailv1 "email-operator/api/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type EmailSenderConfigReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *EmailSenderConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	var emailSenderConfig emailv1.EmailSenderConfig
	if err := r.Get(ctx, req.NamespacedName, &emailSenderConfig); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if emailSenderConfig.ObjectMeta.Generation == 1 {
		r.Log.Info("EmailSenderConfig created", "name", req.Name)
	} else {
		r.Log.Info("EmailSenderConfig updated", "name", req.Name)
	}

	return ctrl.Result{}, nil
}

func (r *EmailSenderConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&emailv1.EmailSenderConfig{}).
		Complete(r)
}
