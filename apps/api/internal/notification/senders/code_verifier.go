package senders

type CodeGenerator interface {
	VerifyCode(verificationID, code string) error
}

type CodeGeneratorInfo struct {
	ID             string `json:"id,omitempty"`
	VerificationID string `json:"verificationId,omitempty"`
}

func (c *CodeGeneratorInfo) GetID() string {
	if c == nil {
		return ""
	}
	return c.ID
}

func (c *CodeGeneratorInfo) GetVerificationID() string {
	if c == nil {
		return ""
	}
	return c.VerificationID
}
