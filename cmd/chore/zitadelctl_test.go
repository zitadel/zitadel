package chore_test

import (
	. "github.com/onsi/ginkgo"
	"path/filepath"
	"regexp"
	"time"
)

var _ = Describe("orbctl", func() {

	var (
		tag, orbURL, workfolder, orbconfig, ghTokenPath, accessToken string
	)
	BeforeSuite(func() {
		workfolder = "./artifacts"
		orbconfig = filepath.Join(workfolder, "orbconfig")
		ghTokenPath = filepath.Join(workfolder, "ghtoken")
		tag = prefixedEnv("TAG")
		orbURL = prefixedEnv("ORBURL")
		accessToken = prefixedEnv("GITHUB_ACCESS_TOKEN")
	})

	Context("version", func() {
		When("the orbctl is downloaded from github releases", func() {
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
}
