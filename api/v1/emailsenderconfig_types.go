package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type EmailSenderConfigSpec struct {
	ApiTokenSecretRef string `json:"apiTokenSecretRef,omitempty"`
	SenderEmail       string `json:"senderEmail,omitempty"`
}

type EmailSenderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EmailSenderConfigSpec   `json:"spec,omitempty"`
}

type EmailSenderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EmailSenderConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EmailSenderConfig{}, &EmailSenderConfigList{})
}
