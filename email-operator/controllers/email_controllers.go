package controllers

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "bytes"

    mailv1alpha1 "github.com/mailer/email-operator/api/v1alpha1"
    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/types"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
    "sigs.k8s.io/controller-runtime/pkg/log"
)

// Helper function to get the API token from the Kubernetes Secret
func (r *EmailReconciler) getAPIToken(ctx context.Context, secretRef string, namespace string) (string, error) {
    var secret corev1.Secret
    if err := r.Get(ctx, types.NamespacedName{Name: secretRef, Namespace: namespace}, &secret); err != nil {
        return "", err
    }
    tokenBytes, ok := secret.Data["apiToken"]
    if !ok {
        return "", fmt.Errorf("apiToken not found in secret %s", secretRef)
    }
    return string(tokenBytes), nil
}

// Function to send email using MailerSend
func (r *EmailReconciler) sendEmail(token string, email mailv1alpha1.Email) (string, string, error) {
    url := "https://api.mailersend.com/v1/email"
    requestPayload := map[string]interface{}{
        "from":    {"email": email.Spec.SenderConfigRef},
        "to":      []map[string]string{{"email": email.Spec.RecipientEmail}},
        "subject": email.Spec.Subject,
        "text":    email.Spec.Body,
    }
    payloadBytes, err := json.Marshal(requestPayload)
    if err != nil {
        return "", "", err
    }

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return "", "", err
    }

    req.Header.Set("Authorization", "Bearer "+token)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return "", "", err
    }
    defer resp.Body.Close()

    var respBody map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
        return "", "", err
    }

    if resp.StatusCode != http.StatusOK {
        return "", "", fmt.Errorf("failed to send email: %s", respBody["message"])
    }

    messageId, _ := respBody["message_id"].(string)
    return "sent", messageId, nil
}

// Reconcile function modified for Email resource
func (r *EmailReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    _ = log.FromContext(ctx)

    var email mailv1alpha1.Email
    if err := r.Get(ctx, req.NamespacedName, &email); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // Implementation to handle email sending logic
    token, err := r.getAPIToken(ctx, email.Spec.SenderConfigRef, req.Namespace)
    if err != nil {
        // Handle errors properly
        return ctrl.Result{}, err
    }

    status, messageId, err := r.sendEmail(token, email)
    if err != nil {
        email.Status.DeliveryStatus = "failed"
        email.Status.Error = err.Error()
    } else {
        email.Status.DeliveryStatus = status
        email.Status.Message
		}
	}