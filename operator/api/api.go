package api

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/common"
	"gopkg.in/yaml.v3"
)

const (
	zitadelFile  = "zitadel.yml"
	databaseFile = "database.yml"
)

type PushDesiredFunc func(monitor mntr.Monitor) error

func ExistsZitadelYml(gitClient *git.Client) (bool, error) {
	return existsFileInGit(gitClient, zitadelFile)
}

func ReadZitadelYml(gitClient *git.Client) (*tree.Tree, error) {
	return readFileInGit(gitClient, zitadelFile)
}

func PushZitadelYml(monitor mntr.Monitor, msg string, gitClient *git.Client, desired *tree.Tree) (err error) {
	return pushFileInGit(monitor, msg, gitClient, desired, zitadelFile)
}

func PushZitadelDesiredFunc(gitClient *git.Client, desired *tree.Tree) PushDesiredFunc {
	return func(monitor mntr.Monitor) error {
		monitor.Info("Writing zitadel desired state")
		return PushZitadelYml(monitor, "Zitadel desired state written", gitClient, desired)
	}
}

func ExistsDatabaseYml(gitClient *git.Client) (bool, error) {
	return existsFileInGit(gitClient, databaseFile)
}

func ReadDatabaseYml(gitClient *git.Client) (*tree.Tree, error) {
	return readFileInGit(gitClient, databaseFile)
}

func PushDatabaseYml(monitor mntr.Monitor, msg string, gitClient *git.Client, desired *tree.Tree) (err error) {
	return pushFileInGit(monitor, msg, gitClient, desired, databaseFile)
}

func PushDatabaseDesiredFunc(gitClient *git.Client, desired *tree.Tree) PushDesiredFunc {
	return func(monitor mntr.Monitor) error {
		monitor.Info("Writing database desired state")
		return PushDatabaseYml(monitor, "Database desired state written", gitClient, desired)
	}
}

func pushFileInGit(monitor mntr.Monitor, msg string, gitClient *git.Client, desired *tree.Tree, path string) (err error) {
	monitor.OnChange = func(_ string, fields map[string]string) {
		err = gitClient.UpdateRemote(mntr.SprintCommit(msg, fields), git.File{
			Path:    path,
			Content: common.MarshalYAML(desired),
		})
		mntr.LogMessage(msg, fields)
	}
	monitor.Changed(msg)
	return err
}

func existsFileInGit(gitClient *git.Client, path string) (bool, error) {
	of := gitClient.Read(path)
	if of != nil && len(of) > 0 {
		return true, nil
	}
	return false, nil
}

func readFileInGit(gitClient *git.Client, path string) (*tree.Tree, error) {
	tree := &tree.Tree{}
	if err := yaml.Unmarshal(gitClient.Read(path), tree); err != nil {
		return nil, err
	}

	return tree, nil
}
