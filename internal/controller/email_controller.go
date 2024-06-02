package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	emailv1 "email-operator/api/v1"
)

type EmailReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("email", req.NamespacedName)

	log.Info("Starting reconciliation for Email resource")

	var email emailv1.Email
	if err := r.Get(ctx, req.NamespacedName, &email); err != nil {
		log.Error(err, "unable to fetch Email")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Info("Fetched Email resource", "email", email)

	if email.Status.DeliveryStatus != "" {
		log.Info("Email already processed", "email", req.NamespacedName)
		return ctrl.Result{}, nil
	}

	var senderConfig emailv1.EmailSenderConfig
	if err := r.Get(ctx, client.ObjectKey{Name: email.Spec.SenderConfigRef, Namespace: req.Namespace}, &senderConfig); err != nil {
		log.Error(err, "unable to fetch EmailSenderConfig")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = "SenderConfig not found"
		if err := r.Status().Update(ctx, &email); err != nil {
			log.Error(err, "unable to update Email status")
		}
		return ctrl.Result{}, err
	}

	log.Info("Fetched EmailSenderConfig", "senderConfig", senderConfig)

	var secret corev1.Secret
	if err := r.Get(ctx, client.ObjectKey{Name: senderConfig.Spec.ApiTokenSecretRef, Namespace: req.Namespace}, &secret); err != nil {
		log.Error(err, "unable to fetch Secret for API token")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = "API token secret not found"
		if err := r.Status().Update(ctx, &email); err != nil {
			log.Error(err, "unable to update Email status")
		}
		return ctrl.Result{}, err
	}

	apiToken, ok := secret.Data["apiToken"]
	if !ok {
		log.Error(fmt.Errorf("apiToken key not found in secret"), "unable to fetch API token from Secret")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = "API token not found in secret"
		if err := r.Status().Update(ctx, &email); err != nil {
			log.Error(err, "unable to update Email status")
		}
		return ctrl.Result{}, fmt.Errorf("apiToken key not found in secret %s", senderConfig.Spec.ApiTokenSecretRef)
	}

	messageID, err := sendEmail(string(apiToken), senderConfig.Spec.SenderEmail, email.Spec.RecipientEmail, email.Spec.Subject, email.Spec.Body)
	if err != nil {
		log.Error(err, "failed to send email")
		email.Status.DeliveryStatus = "Failed"
		email.Status.Error = err.Error()
	} else {
		log.Info("Email sent successfully")
		email.Status.DeliveryStatus = "Success"
		email.Status.MessageID = messageID
	}

	// Update the Email status
	if err := r.Status().Update(ctx, &email); err != nil {
		log.Error(err, "unable to update Email status")
		return ctrl.Result{}, err
	}

	log.Info("Updated Email status", "status", email.Status)

	return ctrl.Result{}, nil
}

func sendEmail(apiToken, senderEmail, recipientEmail, subject, body string) (string, error) {
	url := "https://api.mailersend.com/v1/email"

	type EmailPayload struct {
		From struct {
			Email string `json:"email"`
		} `json:"from"`
		To []struct {
			Email string `json:"email"`
		} `json:"to"`
		Subject string `json:"subject"`
		Text    string `json:"text,omitempty"`
	}

	payload := EmailPayload{
		Subject: subject,
		Text:    body,
	}
	payload.From.Email = senderEmail
	payload.To = append(payload.To, struct {
		Email string `json:"email"`
	}{Email: recipientEmail})

	jsonPayload, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", url, ioutil.NopCloser(bytes.NewReader(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseBody, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response: %s\n", responseBody)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return "", fmt.Errorf("failed to send email: %s", resp.Status)
	}

	messageID := resp.Header.Get("X-Message-Id")
	return messageID, nil
}

func (r *EmailReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&emailv1.Email{}).
		Complete(r)
}
