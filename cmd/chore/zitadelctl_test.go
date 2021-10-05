package chore_test

import (
	"context"
	"fmt"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/cli"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/cmd/chore"
	"github.com/caos/zitadel/pkg/databases"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

var _ = Describe("orbctl", func() {

	const (
		envPrefix       = "ZITADEL_E2E_"
		tagEnv          = envPrefix + "TAG"
		userEnv         = envPrefix + "USER"
		cleanupAfterEnv = envPrefix + "CLEANUP_AFTER"
		reuseOrbEnv     = envPrefix + "REUSE_ORB"
	)
	var (
		tag, orbconfigPath, workfolder, user string
		monitor                              mntr.Monitor
		k8sClient                            kubernetes.ClientInt
		gitClient                            *git.Client
		kubectl                              kubectlCmd
		zitadelctlGitops                     zitadelctlGitopsCmd
		AwaitCompletedPodFromJob             awaitCompletedPodFromJob
		AwaitReadyPods                       awaitReadyPods
		AwaitSecret                          awaitSecret
	)
	BeforeSuite(func() {
		workfolder = "./artifacts"
		kubeconfigPath := filepath.Join(workfolder, "kubeconfig")
		orbconfigPath = filepath.Join(workfolder, "orbconfig")
		tag = prefixedEnv("TAG")
		user = prefixedEnv("USER")
		zitadelctlGitops = zitadelctlGitopsFunc(orbconfigPath)
		kubectl = kubectlCmdFunc(kubeconfigPath)
		monitor = mntr.Monitor{}
		AwaitCompletedPodFromJob = awaitCompletedPodFromJobFunc(kubectl)
		AwaitReadyPods = awaitReadyPodsFunc(kubectl)
		AwaitSecret = awaitSecretFunc(kubectl)

		orbconfig, err := orb.ParseOrbConfig(orbconfigPath)
		Expect(err).ToNot(HaveOccurred())
		kubeconfig, _ := ioutil.ReadFile(kubeconfigPath)
		kubeconfigStr := string(kubeconfig)
		k8sClientT, _ := kubernetes.NewK8sClient(monitor, &kubeconfigStr)
		k8sClient = k8sClientT
		gitClient = git.New(context.Background(), monitor, "zitadelctl", "test@orbos.ch")
		err = cli.InitRepo(orbconfig, gitClient)
		Expect(err).ToNot(HaveOccurred())

		Expect(tag).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", tagEnv))
		Expect(user).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", userEnv))
	})

	Context("version", func() {
		When("the zitadelctl is downloaded from github releases", func() {
			It("contains the tag read from environment variable", func() {

				cmdFunc, err := chore.Command(false, false, true, tag)
				Expect(err).ToNot(HaveOccurred())

				cmd := cmdFunc(context.Background())
				cmd.Args = append(cmd.Args, "--version")

				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 1*time.Second).Should(gexec.Exit(0))
				Eventually(session, 1*time.Minute, 1*time.Second).Should(gbytes.Say(regexp.QuoteMeta(fmt.Sprintf("zitadelctl version %s", tag))))
			})
		})
	})

	Context("repository initialization", func() {
		When("initialized repository access", func() {
			It("configures successfully", func() {
				configureSession, err := gexec.Start(zitadelctlGitops("configure"), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(configureSession, 2*time.Minute, 1*time.Second).Should(gexec.Exit(0))
			})
		})
	})
	Context("database", func() {
		When("bootstraping", func() {
			It("succeeds when creating the initial database.yml", func() {
				localToRemoteFile(zitadelctlGitops, "database.yml", "./templates/database.yml", os.Getenv)
			})
			It("deploy cockroachdb with 1 node", func() {
				session, err := gexec.Start(zitadelctlGitops("takeoff"), os.Stdout, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())

				Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))
				AwaitReadyPods("caos-zitadel", "app.kubernetes.io/name=cockroachdb", 1, 5)
			})
			It("deploys job to test cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-root.yml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, "caos-zitadel", "job-name=cockroachdb-root", 2*time.Minute)
			})
		})
		When("scaling to 3 nodes", func() {
			It("succeeds to scale cockroachdb to 3 nodes", func() {
				count := 3
				session, err := gexec.Start(zitadelctlGitops("file", "patch", "database.yml", "database.spec.replicaCount", "--exact", "--value", strconv.Itoa(count)), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

				AwaitReadyPods("caos-zitadel", "app.kubernetes.io/name=cockroachdb", count, 3)
			})
			It("deploys job to test cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-root.yml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, "caos-zitadel", "job-name=cockroachdb-root", 2*time.Minute)
			})
		})
		When("add and delete users", func() {
			It("add user to DB", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-add.yml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, "caos-zitadel", "job-name=cockroachdb-add", 2*time.Minute)
			})
			It("generate certifcate for user", func() {
				err := databases.GitOpsAddUser(monitor, user, k8sClient, gitClient)
				Expect(err).ToNot(HaveOccurred())
				AwaitSecret("caos-zitadel", "cockroachdb.client."+user, []string{"ca.crt", "client." + user + ".crt", "client." + user + ".key"}, 1*time.Minute)
			})
			It("deploys job to test added user to connect to cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-add.yml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, "caos-zitadel", "job-name=cockroachdb-user", 2*time.Minute)
			})
			It("delete user from DB", func() {
				err := databases.GitOpsDeleteUser(monitor, user, k8sClient, gitClient)
				Expect(err).ToNot(HaveOccurred())
			})
			It("deploys job to test deleted user to connect to cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-add.yml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, "caos-zitadel", "job-name=cockroachdb-user", 2*time.Minute)
			})
		})
	})

	Context("zitadel", func() {
		When("bootstraping", func() {
			It("succeeds when creating the initial zitadel.yml", func() {
				localToRemoteFile(zitadelctlGitops, "zitadel.yml", "./templates/zitadel.yml", os.Getenv)
			})
		})
	})
})
