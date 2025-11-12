package config

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/grpc/oidc/v2"
	"github.com/zitadel/zitadel/internal/api/saml"
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/config/network"
	"github.com/zitadel/zitadel/internal/database"
)

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

type Config struct {
	ExternalPort    uint16
	ExternalDomain  string
	ExternalSecure  bool
	Database        database.Config
	TLS             network.TLS
	SystemAPIUsers  map[string]*authz.SystemAPIUser
	OIDC            oidc.Config
	SAML            saml.Config
	DefaultInstance command.InstanceSetup
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "manage configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			/*			stdin, err := io.ReadAll(os.Stdin)
						if err != nil {
							return fmt.Errorf("Error reading from stdin: %v", err)
						}*/
			node := yaml.Node{
				Kind: yaml.DocumentNode,
				Content: []*yaml.Node{
					{
						Kind: yaml.MappingNode,
					},
				},
			}
			/*			if err := yaml.Unmarshal(stdin, &node); err != nil {
						return fmt.Errorf("Error unmarshaling YAML: %v", err)
					}*/
			var cfg Config
			if err := node.Content[0].Decode(&cfg); err != nil {
				return fmt.Errorf("Error decoding YAML to Config: %v", err)
			}
			note, err := addSystemUser(&cfg, node.Content[0])
			if err != nil {
				return fmt.Errorf("Error adding System user: %v", err)
			}
			result, err := yaml.Marshal(&node)
			if err != nil {
				return fmt.Errorf("Error encoding Config to YAML: %v", err)
			}
			_, err = os.Stdout.Write(result)
			if err != nil {
				return fmt.Errorf("Error writing YAML to stdout: %v", err)
			}

			// Print follow-up instructions to stderr with colors
			os.Stderr.Write([]byte("\n" + colorCyan + strings.Repeat("=", 80) + colorReset + "\n"))
			os.Stderr.Write([]byte(colorBold + colorGreen + "NEXT STEPS" + colorReset + "\n"))
			os.Stderr.Write([]byte(colorCyan + strings.Repeat("=", 80) + colorReset + "\n\n"))

			os.Stderr.Write([]byte(colorBold + "1. Generate a master key and store it securely:" + colorReset + "\n\n"))
			os.Stderr.Write([]byte("   " + colorYellow + "$ zitadel config generate-masterkey > masterkey.txt" + colorReset + "\n\n"))
			os.Stderr.Write([]byte("   Then pass it to ZITADEL with: " + colorCyan + "--masterkeyFile=./masterkey.txt" + colorReset + "\n"))

			if note != "" {
				os.Stderr.Write([]byte("\n" + colorCyan + strings.Repeat("-", 80) + colorReset + "\n\n"))
				_, err = os.Stderr.Write([]byte(note))
			}

			os.Stderr.Write([]byte("\n" + colorCyan + strings.Repeat("=", 80) + colorReset + "\n"))
			return nil
		},
	}
	cmd.AddCommand(newGenerateMasterKeyCmd())
	cmd.AddCommand(newGenerateRSAPrivateKeyCmd())
	cmd.AddCommand(newGenerateRSAPublicKeyCmd())
	return cmd
}

func addSystemUser(in *Config, out *yaml.Node) (string, error) {
	if len(in.SystemAPIUsers) != 0 {
		// TODO: check if login-system-user exists
		return "", nil
	}
	prompt := promptui.Prompt{
		Label:     "Add login-system-user for the Login authentication (Y/n)",
		IsConfirm: true,
		Default:   "Y",
		Stdout:    os.Stderr,
	}
	_, err := prompt.Run()
	if err != nil {
		prompt.Stdout.Write([]byte(fmt.Sprintf("Not adding System user %v\n", err)))
		return "", nil
	}
	prompt.Stdout.Write([]byte("Adding System user\n"))
	var privateKeyBuf bytes.Buffer
	if err := generateRSAPrivateKey(&privateKeyBuf); err != nil {
		return "", fmt.Errorf("Error generating RSA private key: %v", err)
	}
	privateKey := privateKeyBuf.Bytes()
	base64PrivateKey := base64.StdEncoding.EncodeToString(privateKey)
	publicKeyBuf := bytes.Buffer{}
	if err := generateRSAPublicKey(privateKey, &publicKeyBuf); err != nil {
		return "", fmt.Errorf("Error generating RSA public key: %v", err)
	}
	base64PublicKey := base64.StdEncoding.EncodeToString(publicKeyBuf.Bytes())
	keyData := bytesToString([]byte(base64PublicKey))
	in.SystemAPIUsers = map[string]*authz.SystemAPIUser{
		"login-system-user": {
			KeyData: keyData,
			Memberships: []*authz.Membership{{
				MemberType: authz.MemberTypeSystem,
				Roles:      []string{"IAM_LOGIN_CLIENT"},
			}},
		},
	}

	// Create YAML nodes manually to ensure KeyData is a string
	membershipsNode := &yaml.Node{}
	if err := membershipsNode.Encode(in.SystemAPIUsers["login-system-user"].Memberships); err != nil {
		return "", fmt.Errorf("Error encoding Memberships to YAML node: %v", err)
	}

	systemUsersNode := &yaml.Node{
		Kind: yaml.MappingNode,
		Content: []*yaml.Node{
			// login-system-user key
			{Kind: yaml.ScalarNode, Value: "login-system-user"},
			// login-system-user value
			{
				Kind: yaml.MappingNode,
				Content: []*yaml.Node{
					{Kind: yaml.ScalarNode, Value: "keydata"},
					{Kind: yaml.ScalarNode, Value: base64PublicKey, Style: yaml.LiteralStyle},
					{Kind: yaml.ScalarNode, Value: "memberships"},
					membershipsNode,
				},
			},
		},
	}

	// Add SystemAPIUsers key and value to the document
	out.Content = append(out.Content,
		&yaml.Node{
			Kind:  yaml.ScalarNode,
			Value: "SystemAPIUsers",
		},
		systemUsersNode,
	)

	note := fmt.Sprintf(colorBold+"2. Store the private key securely:"+colorReset+"\n\n"+
		"   "+colorYellow+"$ echo '%s' > login-client.pem"+colorReset+"\n\n"+
		"   Then pass it to the login app with: "+colorCyan+"SYSTEM_USER_PRIVATE_KEY_FILE=login-client.pem"+colorReset,
		base64PrivateKey)

	return note, nil
}

func newGenerateMasterKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-masterkey",
		Short: "generate a random master key",
		RunE: func(_ *cobra.Command, _ []string) error {
			key := make([]byte, 32)
			_, err := rand.Read(key)
			if err != nil {
				return fmt.Errorf("Error generating random bytes: %v", err)
			}
			for _, b := range key {
				fmt.Printf("%02x", b)
			}
			return nil
		},
	}
	return cmd
}

func newGenerateRSAPrivateKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-rsa-private-key",
		Short: "generate a random RSA private key",
		RunE: func(_ *cobra.Command, _ []string) error {
			return generateRSAPrivateKey(os.Stdout)
		},
	}
	return cmd
}

func newGenerateRSAPublicKeyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-rsa-public-key",
		Short: "generate a random RSA public key",
		RunE: func(cmd *cobra.Command, args []string) error {
			stdin, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("Error reading from stdin: %v", err)
			}
			decoded, err := base64.StdEncoding.DecodeString(string(stdin))
			if err != nil {
				return fmt.Errorf("Error decoding base64 input: %v", err)
			}
			stdin = decoded
			return generateRSAPublicKey(stdin, os.Stdout)
		},
	}
	return cmd
}

func generateRSAPrivateKey(out io.Writer) error {
	cryptoKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("Error generating RSA private key: %v", err)
	}
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(cryptoKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	err = pem.Encode(out, pemBlock)
	if err != nil {
		return fmt.Errorf("Error encoding RSA private key to PEM: %v", err)
	}
	return nil
}

func generateRSAPublicKey(privateKey []byte, out io.Writer) error {
	// Decode PEM block
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return fmt.Errorf("Error decoding PEM block from private key")
	}
	pk, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("Error parsing RSA private key: %v", err)
	}
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return fmt.Errorf("Error marshaling RSA public key: %v", err)
	}
	pemBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	err = pem.Encode(out, pemBlock)
	if err != nil {
		return fmt.Errorf("Error encoding RSA public key to PEM: %v", err)
	}
	return nil
}

type bytesToString []byte

func (b *bytesToString) MarshalYAML() (interface{}, error) {
	return string(*b), nil
}
