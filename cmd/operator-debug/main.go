package main

import (
	"flag"
	"io/ioutil"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator/crtlgitops"
	"github.com/caos/zitadel/operator/helpers"
)

func main() {
	orbconfig := flag.String("orbconfig", "~/.orb/config", "The orbconfig file to use")
	kubeconfig := flag.String("kubeconfig", "~/.kube/config", "The kubeconfig file to use")
	verbose := flag.Bool("verbose", false, "Print debug levelled logs")
	gitops := flag.Bool("gitops", false, "defines if the ZITADEL operator should run in gitops mode")

	flag.Parse()

	monitor := mntr.Monitor{
		OnInfo:   mntr.LogMessage,
		OnChange: mntr.LogMessage,
		OnError:  mntr.LogError,
	}

	if *verbose {
		monitor = monitor.Verbose()
	}

	kc, err := ioutil.ReadFile(helpers.PruneHome(*kubeconfig))
	if err != nil {
		panic(err)
	}

	if err := crtlgitops.Operator(
		monitor,
		helpers.PruneHome(*orbconfig),
		kubernetes.NewK8sClient(monitor, strPtr(string(kc))),
		strPtr("local-debugging"),
		*gitops,
	); err != nil {
		panic(err)
	}
}

func strPtr(str string) *string {
	return &str
}
