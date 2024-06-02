package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)
type EmailSpec struct {
	SenderConfigRef string `json:"senderConfigRef,omitempty"`
	RecipientEmail  string `json:"recipientEmail,omitempty"`
	Subject         string `json:"subject,omitempty"`
	Body            string `json:"body,omitempty"`
}

type EmailStatus struct {
	DeliveryStatus string `json:"deliveryStatus,omitempty"`
	MessageID      string `json:"messageId,omitempty"`
	Error          string `json:"error,omitempty"`
}
type Email struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSpec   `json:"spec,omitempty"`
	Status EmailStatus `json:"status,omitempty"`
}

type EmailList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Email `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Email{}, &EmailList{})
}
