package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/zitadel/zitadel/apps/cli/internal/auth"
	"github.com/zitadel/zitadel/apps/cli/internal/config"
)

func newLoginCmd() *cobra.Command {
	var (
		instance string
		clientID string
		ctxName  string
	)

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with a ZITADEL instance using the OIDC Device Authorization flow",
		RunE: func(cmd *cobra.Command, args []string) error {
			if token := os.Getenv("ZITADEL_TOKEN"); token != "" {
				fmt.Fprintln(os.Stderr, "ZITADEL_TOKEN is set — PAT mode is active. Browser login skipped.")
				return nil
			}

			// If any flag is provided, fall back to non-interactive mode.
			if instance != "" || clientID != "" {
				if instance == "" {
					return fmt.Errorf("--instance is required. Run 'zitadel-cli login setup' for instructions.")
				}
				if clientID == "" {
					return fmt.Errorf("--client-id is required. Run 'zitadel-cli login setup' for instructions.")
				}
				return runLogin(cmd, instance, clientID, ctxName)
			}

			// Otherwise, run interactive wizard
			fmt.Fprintln(os.Stderr, "Interactive Login Wizard")
			fmt.Fprintln(os.Stderr, "------------------------")

			qs := []*survey.Question{
				{
					Name: "instance",
					Prompt: &survey.Input{
						Message: "What is your ZITADEL instance URL? (e.g., mycompany.zitadel.cloud)",
					},
					Validate: survey.Required,
				},
				{
					Name: "authMethod",
					Prompt: &survey.Select{
						Message: "How would you like to authenticate?",
						Options: []string{"Browser (Device Flow)", "Personal Access Token (PAT)"},
						Default: "Browser (Device Flow)",
					},
				},
			}

			answers := struct {
				Instance   string
				AuthMethod string
			}{}

			if err := survey.Ask(qs, &answers); err != nil {
				return err
			}

			instance = answers.Instance

			if answers.AuthMethod == "Personal Access Token (PAT)" {
				var pat string
				err := survey.AskOne(&survey.Password{
					Message: "Enter your Personal Access Token:",
				}, &pat, survey.WithValidator(survey.Required))
				if err != nil {
					return err
				}

				err = survey.AskOne(&survey.Input{
					Message: "Name this context (so you can switch to it later):",
					Default: cleanContextName(instance),
				}, &ctxName)
				if err != nil {
					return err
				}

				return runPATLogin(instance, pat, ctxName)
			}

			// Device Flow path
			err := survey.AskOne(&survey.Input{
				Message: "What is your Application Client ID? (Run 'zitadel-cli login setup' if you don't have one)",
			}, &clientID, survey.WithValidator(survey.Required))
			if err != nil {
				return err
			}

			err = survey.AskOne(&survey.Input{
				Message: "Name this context (so you can switch to it later):",
				Default: cleanContextName(instance),
			}, &ctxName)
			if err != nil {
				return err
			}

			return runLogin(cmd, instance, clientID, ctxName)
		},
	}

	cmd.Flags().StringVar(&instance, "instance", "", "ZITADEL instance host (e.g. mycompany.zitadel.cloud)")
	cmd.Flags().StringVar(&clientID, "client-id", "", "OIDC client ID of a native application")
	cmd.Flags().StringVar(&ctxName, "context", "", "name for this context (defaults to instance host)")

	cmd.AddCommand(newSetupCmd())

	return cmd
}

func cleanContextName(s string) string {
	s = strings.TrimPrefix(s, "https://")
	s = strings.TrimPrefix(s, "http://")
	return s
}

func runLogin(cmd *cobra.Command, instance, clientID, ctxName string) error {
	if ctxName == "" {
		ctxName = cleanContextName(instance)
	}

	token, err := auth.Login(cmd.Context(), instance, clientID)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	ctx := config.Context{
		Instance:   instance,
		AuthMethod: "interactive",
		ClientID:   clientID,
		Token:      token.AccessToken,
	}
	if token.RefreshToken != "" {
		ctx.RefreshToken = token.RefreshToken
	}
	if !token.Expiry.IsZero() {
		ctx.TokenExpiry = token.Expiry.Format(time.RFC3339)
	}
	
	// Load config before saving
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	
	cfg.Contexts[ctxName] = ctx
	cfg.ActiveContext = ctxName

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Logged in to %s (context: %s)\n", instance, ctxName)
	return nil
}

func runPATLogin(instance, pat, ctxName string) error {
	if ctxName == "" {
		ctxName = cleanContextName(instance)
	}

	ctx := config.Context{
		Instance:   instance,
		AuthMethod: "pat",
		PAT:        pat,
	}

	// Load config before saving
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	
	cfg.Contexts[ctxName] = ctx
	cfg.ActiveContext = ctxName

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	fmt.Fprintf(os.Stderr, "✓ Added PAT for %s (context: %s)\n", instance, ctxName)
	return nil
}
