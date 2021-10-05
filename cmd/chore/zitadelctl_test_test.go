package chore_test

import (
	"context"
	"fmt"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/cli"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/pkg/databases"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestChore_user(t *testing.T) {
	monitor := mntr.Monitor{}
	user := "test"
	workfolder := "./artifacts"
	kubeconfigPath := filepath.Join(workfolder, "kubeconfig")
	kubeconfig, _ := ioutil.ReadFile(kubeconfigPath)
	kubeconfigStr := string(kubeconfig)
	k8sClientT, _ := kubernetes.NewK8sClient(monitor, &kubeconfigStr)
	k8sClient := k8sClientT
	gitClient := git.New(context.Background(), monitor, "zitadelctl", "test@orbos.ch")
	orbconfigPath := filepath.Join(workfolder, "orbconfig")
	orbconfig, err := orb.ParseOrbConfig(orbconfigPath)
	assert.NoError(t, err)
	assert.NoError(t, cli.InitRepo(orbconfig, gitClient))

	assert.NoError(t, databases.GitOpsAddUser(monitor, user, k8sClient, gitClient))

	sec := getSecretKeysWithName(kubectlCmdFunc(kubeconfigPath), "caos-zitadel", "cockroachdb.client."+user)
	fmt.Println(sec)
}
