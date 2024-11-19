package senders

type CodeGenerator interface {
	VerifyCode(verificationID, code string) error
}

type CodeGeneratorInfo struct {
	Success        bool   `json:"success,omitempty"`
	ErrorMessage   string `json:"errorMessage,omitempty"`
	ID             string `json:"id,omitempty"`
	VerificationID string `json:"verificationId,omitempty"`
}

func (c *CodeGeneratorInfo) GetSuccess() bool {
	if c == nil {
		return false
	}
	return c.Success
}

func (c *CodeGeneratorInfo) GetErrorMessage() string {
	if c == nil {
		return ""
	}
	return c.ErrorMessage
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
