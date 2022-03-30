package key

import (
	"errors"
	"io/ioutil"

	"github.com/spf13/cobra"
)

const (
	flagMasterKey      = "masterkeyFile"
	flagMasterKeyShort = "m"
	flagMasterKeyArg   = "masterkey"
)

var (
	ErrNotSingleFlag = errors.New("masterkey must either be provided by file path or value")
)

func AddMasterKeyFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringP(flagMasterKey, flagMasterKeyShort, "", "path to the masterkey for en/decryption keys")
	cmd.PersistentFlags().String(flagMasterKeyArg, "", "masterkey as argument for en/decryption keys")
}

func MasterKey(cmd *cobra.Command) (string, error) {
	masterKeyFile, _ := cmd.Flags().GetString(flagMasterKey)
	masterKeyFromArg, _ := cmd.Flags().GetString(flagMasterKeyArg)
	if err := checkSingleFlag(masterKeyFile, masterKeyFromArg); err != nil {
		return "", err
	}
	if masterKeyFromArg != "" {
		return masterKeyFromArg, nil
	}
	data, err := ioutil.ReadFile(masterKeyFile)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func checkSingleFlag(masterKeyFile, masterKeyFromArg string) error {
	if masterKeyFile != "" && masterKeyFromArg != "" {
		return ErrNotSingleFlag
	}
	if masterKeyFile == "" && masterKeyFromArg == "" {
		return ErrNotSingleFlag
	}
	return nil
}
