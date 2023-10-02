package mail

import (
	"fmt"
)

type verificationSender struct{}

func NewVerificationSender() *verificationSender {
	return &verificationSender{}
}

func (vs verificationSender) Send(email, username, id string) error {
	// TODO: Send simple html formatted page instead
	verificationLink := fmt.Sprintf("http://localhost:3000/api/verification/%s/%s", username, id)
	return Sender.Send(email, "Account verification", verificationLink)
}
