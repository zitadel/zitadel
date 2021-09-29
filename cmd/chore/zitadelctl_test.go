package chore_test

import (
	"context"
	"fmt"
	"github.com/caos/zitadel/cmd/chore"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

var _ = Describe("orbctl", func() {

	const (
		envPrefix       = "ZITADEL_E2E_"
		tagEnv          = envPrefix + "TAG"
		orbEnv          = envPrefix + "ORBURL"
		ghpatEnv        = envPrefix + "GITHUB_ACCESS_TOKEN"
		cleanupAfterEnv = envPrefix + "CLEANUP_AFTER"
		reuseOrbEnv     = envPrefix + "REUSE_ORB"
	)
	var (
		tag, orbconfig, workfolder, orbURL, ghTokenPath, accessToken string
		zitadelctlGitops                                             zitadelctlGitopsCmd
	)
	BeforeSuite(func() {
		workfolder = "./artifacts"
		orbconfig = filepath.Join(workfolder, "orbconfig")
		ghTokenPath = filepath.Join(workfolder, "ghtoken")
		tag = prefixedEnv("TAG")
		orbURL = prefixedEnv("ORBURL")
		accessToken = prefixedEnv("GITHUB_ACCESS_TOKEN")
		zitadelctlGitops = zitadelctlGitopsFunc(orbconfig)

		Expect(tag).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", tagEnv))
		Expect(orbURL).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", orbEnv))
		Expect(accessToken).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", ghpatEnv))
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
				Eventually(session, 2*time.Minute, 1*time.Second).Should(gbytes.Say(regexp.QuoteMeta(tag)))
			})
		})
	})

	Context("repository initialization", func() {
		When("initializing local files", func() {
			It("ensures the ghtoken cache file so that the oauth flow is skipped", func() {

				ghtoken, err := os.Create(ghTokenPath)
				Expect(err).ToNot(HaveOccurred())
				defer ghtoken.Close()
				Expect(ghtoken.WriteString(fmt.Sprintf(`IDToken: ""
IDTokenClaims: null
access_token: %s
expiry: "0001-01-01T00:00:00Z"
token_type: bearer`, accessToken))).To(BeNumerically(">", 0))
			})
			When("configure command is executed for the first time", func() {
				It("creates a new orbconfig containing a new masterkey and a new ssh private key and adds the public key to the repository", func() {

					masterKeySession, err := gexec.Start(exec.Command("openssl", "rand", "-base64", "21"), nil, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())
					Eventually(masterKeySession).Should(gexec.Exit(0))

					configureSession, err := gexec.Start(zitadelctlGitops("configure", "--repourl", orbURL, "--masterkey", string(masterKeySession.Out.Contents())), GinkgoWriter, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())
					Eventually(configureSession, 2*time.Minute, 1*time.Second).Should(gexec.Exit())
				})
			})

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
	})
})
