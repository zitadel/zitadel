package chore_test

import (
	"context"
	"fmt"
	"github.com/caos/zitadel/cmd/chore"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var _ = Describe("orbctl", func() {

	const (
		envPrefix       = "ZITADEL_E2E_"
		tagEnv          = envPrefix + "TAG"
		cleanupAfterEnv = envPrefix + "CLEANUP_AFTER"
		reuseOrbEnv     = envPrefix + "REUSE_ORB"
	)
	var (
		tag, orbconfig, workfolder string
		kubectl                    kubectlCmd
		zitadelctlGitops           zitadelctlGitopsCmd
		AwaitCompletedPodFromJob   awaitCompletedPodFromJob
		AwaitReadyPods             awaitReadyPods
	)
	BeforeSuite(func() {
		workfolder = "./artifacts"
		orbconfig = filepath.Join(workfolder, "orbconfig")
		tag = prefixedEnv("TAG")
		zitadelctlGitops = zitadelctlGitopsFunc(orbconfig)
		kubectl = kubectlCmdFunc(filepath.Join(workfolder, "kubeconfig"))
		AwaitCompletedPodFromJob = awaitCompletedPodFromJobFunc(kubectl)

		Expect(tag).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", tagEnv))
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
				Eventually(session, 2*time.Minute, 1*time.Second).Should(gexec.Exit(0))
				Eventually(session, 2*time.Minute, 1*time.Second).Should(gbytes.Say(regexp.QuoteMeta(fmt.Sprintf("zitadelctl version %s", tag))))
			})
		})
	})

	Context("repository initialization", func() {
		When("initializing local files", func() {
			Context("initialized repository access", func() {
				When("creating remote initial files", func() {
					It("succeeds when creating the initial database.yml", func() {
						localToRemoteFile(zitadelctlGitops, "database.yml", "./templates/database.yml", os.Getenv)
					})
					It("succeeds when creating the initial zitadel.yml", func() {
						localToRemoteFile(zitadelctlGitops, "zitadel.yml", "./templates/zitadel.yml", os.Getenv)
					})

					//TODO add secrets

					It("configures successfully", func() {
						configureSession, err := gexec.Start(zitadelctlGitops("configure"), GinkgoWriter, GinkgoWriter)
						Expect(err).ToNot(HaveOccurred())
						Eventually(configureSession, 2*time.Minute, 1*time.Second).Should(gexec.Exit(0))
					})
				})
			})
		})
		When("bootstrapping", func() {
			It("deploy cockroachdb with 1 node", func() {
				session, err := gexec.Start(zitadelctlGitops("takeoff"), os.Stdout, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())

				Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))
				AwaitReadyPods("app.kubernetes.io/name=cockroachdb", 1)
			})
			It("deploys job to test cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb.yml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, "app.kubernetes.io/name=cockroach-connect", 5*time.Minute)
			})
			It("scale cockroachdb to 3 nodes", func() {
				session, err := gexec.Start(zitadelctlGitops("file", "patch", "orbiter.yml", "database.spec.replicaCount", "--exact", "--value", "3"), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

				AwaitReadyPods("app.kubernetes.io/name=cockroachdb", 3)
			})
		})
	})
})
