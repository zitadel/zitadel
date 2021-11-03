package chore_test

import (
	"context"
	"fmt"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/cli"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/zitadel/cmd/chore/helpers"
	k8s_test "github.com/caos/zitadel/cmd/chore/helpers/k8s"
	"github.com/caos/zitadel/cmd/chore/helpers/orbctl"
	"github.com/caos/zitadel/pkg/backup"
	"github.com/caos/zitadel/pkg/databases"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var _ = Describe("zitadelctl", func() {

	const (
		zitadelNamespace    = "caos-zitadel"
		envPrefix           = "ZITADEL_E2E_"
		tagEnv              = "TAG"
		shaEnv              = "SHA"
		userEnv             = "DBUSER"
		backupSAJSONPathEnv = "BACKUPSAJSON"
		backupAKIDPathEnv   = "BACKUPAKID"
		backupSAKPathEnv    = "BACKUPSAK"
		orbconfigEnv        = "ORBCONFIG"
		cleanupAfterEnv     = envPrefix + "CLEANUP_AFTER"
		reuseOrbEnv         = envPrefix + "REUSE_ORB"
	)

	var (
		orbctlVersion                                                                                  string
		databaseFile, zitadelFile                                                                      string
		tag, sha, orbconfigPath, kubeconfigPath, workfolder, user, backupSAJson, backupAKID, backupSAK string
		gcsBackupName, s3BackupName                                                                    string
		gcsBackupBucket, s3BackupBucket, s3Endpoint, s3Region                                          string
		backupName, changedBackupName                                                                  string
		monitor                                                                                        mntr.Monitor
		k8sClient                                                                                      kubernetes.ClientInt
		gitClient                                                                                      *git.Client
		kubectl                                                                                        k8s_test.KubectlCmd
		zitadelctlBuildGitops                                                                          helpers_test.ZitadelctlGitopsCmd
		zitadelctlGitops                                                                               helpers_test.ZitadelctlGitopsCmd
		orbctlGitops                                                                                   helpers_test.OrbctlGitopsCmd
		GetLogsOfPod                                                                                   helpers_test.GetLogsOfPod
		ApplyFile                                                                                      helpers_test.ApplyFile
		DeleteFile                                                                                     helpers_test.DeleteFile
		AwaitCompletedPodFromJob                                                                       helpers_test.AwaitCompletedPodFromJob
		AwaitCompletedPod                                                                              helpers_test.AwaitCompletedPod
		AwaitReadyPods                                                                                 helpers_test.AwaitReadyPods
		AwaitSecret                                                                                    helpers_test.AwaitSecret
		AwaitCronJobScheduled                                                                          helpers_test.AwaitCronJobScheduled
		AwaitReadyNodes                                                                                helpers_test.AwaitReadyNodes
		DeleteResource                                                                                 helpers_test.DeleteResource
		DeleteNamespacedResource                                                                       helpers_test.DeleteNamespacedResource
	)
	BeforeSuite(func() {

		databaseFile = "database.yml"
		zitadelFile = "zitadel.yml"
		backupName = "e2e-test"
		changedBackupName = "e2e-test2"
		gcsBackupName = "bucket"
		s3BackupName = "csbucket"
		gcsBackupBucket = "caos-zitadel-e2e-backup"
		s3BackupBucket = "caos-zitadel-e2e-backup"
		s3Endpoint = "https://objects.lpg.cloudscale.ch"
		s3Region = "LPG"
		workfolder = "./artifacts"
		kubeconfigPath = filepath.Join(workfolder, "kubeconfig")
		orbconfigPath = helpers_test.PrefixedEnv(orbconfigEnv)
		tag = helpers_test.PrefixedEnv(tagEnv)
		sha = helpers_test.PrefixedEnv(shaEnv)
		user = helpers_test.PrefixedEnv(userEnv)
		backupSAJson = helpers_test.PrefixedEnv(backupSAJSONPathEnv)
		backupAKID = helpers_test.PrefixedEnv(backupAKIDPathEnv)
		backupSAK = helpers_test.PrefixedEnv(backupSAKPathEnv)

		kubectl = k8s_test.KubectlCmdFunc(kubeconfigPath)
		monitor = mntr.Monitor{
			OnInfo:         mntr.LogMessage,
			OnChange:       mntr.LogMessage,
			OnError:        mntr.LogError,
			OnRecoverPanic: mntr.LogPanic,
		}
		AwaitReadyNodes = helpers_test.AwaitReadyNodesFunc(kubectl)
		ApplyFile = helpers_test.ApplyFileFunc(kubectl)
		DeleteFile = helpers_test.DeleteFileFunc(kubectl)
		AwaitCompletedPod = helpers_test.AwaitCompletedPodFunc(kubectl)
		AwaitCompletedPodFromJob = helpers_test.AwaitCompletedPodFromJobFunc(kubectl)
		AwaitReadyPods = helpers_test.AwaitReadyPodsFunc(kubectl)
		AwaitSecret = helpers_test.AwaitSecretFunc(kubectl)
		AwaitCronJobScheduled = helpers_test.AwaitCronJobScheduledFunc(kubectl)
		GetLogsOfPod = helpers_test.GetLogsOfPodFunc(kubectl)
		DeleteResource = helpers_test.DeleteResourceFunc(kubectl)
		DeleteNamespacedResource = helpers_test.DeleteNamespacedResourceFunc(kubectl)

		orbconfig, err := orb.ParseOrbConfig(orbconfigPath)
		Expect(err).ToNot(HaveOccurred())
		gitClient = git.New(context.Background(), monitor, "zitadelctl", "test@orbos.ch")
		err = cli.InitRepo(orbconfig, gitClient)
		Expect(err).ToNot(HaveOccurred())

		data := gitClient.Read("networking.yml")
		orbctlVersion = orbctl.GetVersion(data)
		zitadelctlGitops = helpers_test.ZitadelctlGitopsFunc(orbconfigPath, tag)
		zitadelctlBuildGitops = helpers_test.ZitadelctlBuildGitopsFunc(orbconfigPath, tag)
		orbctlGitops = helpers_test.OrbctlGitopsFunc(orbconfigPath, orbctlVersion)

		Expect(tag).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+tagEnv))
		Expect(sha).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+shaEnv))
		Expect(user).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+userEnv))
		Expect(orbconfigPath).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+orbconfigEnv))
		Expect(backupSAJson).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+backupSAJSONPathEnv))
		Expect(backupAKID).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+backupAKIDPathEnv))
		Expect(backupSAK).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+backupSAKPathEnv))
	})
	Context("version", func() {
		When("the tests are started", func() {
			It("the docker images should be existing in a specific version", func() {
				cmd := exec.Command("/bin/bash", "-c", "docker pull ghcr.io/caos/zitadel-crbackup:"+tag)
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 1*time.Second).Should(gexec.Exit(0))

				cmd = exec.Command("/bin/bash", "-c", "docker pull ghcr.io/caos/zitadel-operator:"+tag)
				session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 1*time.Second).Should(gexec.Exit(0))

				cmd = exec.Command("/bin/bash", "-c", "docker pull ghcr.io/caos/zitadel:"+tag)
				session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 1*time.Second).Should(gexec.Exit(0))
			})
		})
		When("the orbctl is downloaded from github releases", func() {
			It("contains the tag read from networking.yml", func() {
				cmdFunc, err := orbctl.Command(false, false, true, orbctlVersion)
				Expect(err).ToNot(HaveOccurred())

				cmd := cmdFunc(context.Background())
				cmd.Args = append(cmd.Args, "--version")

				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 1*time.Second).Should(gexec.Exit(0))
				Eventually(session, 2*time.Minute, 1*time.Second).Should(gbytes.Say(regexp.QuoteMeta(orbctlVersion)))
			})
		})
		When("the zitadelctl is downloaded from github releases", func() {
			It("contains the tag read from environment variable", func() {

				session, err := gexec.Start(zitadelctlBuildGitops("--version"), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 1*time.Second).Should(gexec.Exit(0))
				Eventually(string(session.Out.Contents()), 1*time.Minute, 1*time.Second).Should(ContainSubstring(regexp.QuoteMeta(fmt.Sprintf("zitadelctl version %s", tag))))
			})
		})
	})

	Context("orbos", func() {
		When("orbctl takeoff", func() {
			It("either is e running cluster or starts a cluster", func() {
				session, err := gexec.Start(orbctlGitops("takeoff"), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 15*time.Minute, 5*time.Second).Should(gexec.Exit(0))
			})
			It("gives a kubeconfig for the cluster", func() {
				session, err := gexec.Start(orbctlGitops("readsecret", "orbiter.k8s.kubeconfig.encrypted"), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 2*time.Minute, 5*time.Second).Should(gexec.Exit(0))

				absPath, err := filepath.Abs(kubeconfigPath)
				Expect(err).ToNot(HaveOccurred())
				err = ioutil.WriteFile(absPath, session.Out.Contents(), os.ModePerm)
				Expect(err).ToNot(HaveOccurred())

			})
			It("kubernetes is connectable", func() {
				cmd := kubectl("get", "ns")
				session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

				kubeconfig, _ := ioutil.ReadFile(kubeconfigPath)
				kubeconfigStr := string(kubeconfig)
				k8sClientT, _ := kubernetes.NewK8sClient(monitor, &kubeconfigStr)
				k8sClient = k8sClientT
			})
			It("scales up the nodes", func() {
				AwaitReadyNodes(4, 10*time.Minute)
			})
			It("cleanup caos-zitadel", func() {
				DeleteNamespacedResource("deploy", "caos-system", "database-operator")
				DeleteNamespacedResource("deploy", "caos-system", "zitadel-operator")
				DeleteResource("ns", "caos-zitadel")
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
			It("deletes already existing files", func() {
				session, err := gexec.Start(zitadelctlGitops("file", "remove", databaseFile), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 2*time.Minute, 1*time.Second).Should(gexec.Exit(0))
				session, err = gexec.Start(zitadelctlGitops("file", "remove", zitadelFile), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 2*time.Minute, 1*time.Second).Should(gexec.Exit(0))
			})
		})
	})

	Context("database", func() {
		When("bootstraping", func() {
			It("succeeds when creating the initial "+databaseFile, func() {
				helpers_test.LocalToRemoteFile(zitadelctlGitops, databaseFile, "./templates/database.yml", os.Getenv)
			})
			It("deploy cockroachdb with 1 node", func() {
				session, err := gexec.Start(zitadelctlGitops("takeoff"), os.Stdout, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())

				Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))
				AwaitReadyPods(zitadelNamespace, "app.kubernetes.io/name=cockroachdb", 1, 5*time.Minute)
			})
		})
		When("get connect string to cockroach", func() {
			It("logic should give back connection information", func() {
				host, port, err := databases.GitOpsGetConnectionInfo(monitor, k8sClient, gitClient)
				Expect(err).ToNot(HaveOccurred())
				Expect(host).ToNot(BeEmpty())
				Expect(port).ToNot(BeEmpty())
				err = os.Setenv(envPrefix+"DBHOST", host)
				err = os.Setenv(envPrefix+"DBPORT", port)
			})
			It("deploys job to test cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-root.yaml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-root", 2*time.Minute)
			})
		})

		When("scaling to 3 nodes", func() {
			It("succeeds to scale cockroachdb to 3 nodes", func() {
				count := 3
				session, err := gexec.Start(zitadelctlGitops("file", "patch", databaseFile, "database.spec.replicaCount", "--exact", "--value", strconv.Itoa(count)), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))

				AwaitReadyPods(zitadelNamespace, "app.kubernetes.io/name=cockroachdb", count, 3*time.Minute)
			})
			It("deploys job to test cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-root.yaml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-root", 2*time.Minute)
			})
		})

		When("add and delete users", func() {
			It("add user to DB", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-add.yaml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-add", 2*time.Minute)
			})
			It("generate certificate for user", func() {
				err := databases.GitOpsAddUser(monitor, user, k8sClient, gitClient)
				Expect(err).ToNot(HaveOccurred())
				AwaitSecret(zitadelNamespace, "cockroachdb.client."+user, []string{"ca.crt", "client." + user + ".crt", "client." + user + ".key"}, 1*time.Minute)
			})
			It("deploys job to test added user to connect to cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-user.yaml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-user", 2*time.Minute)
			})
			It("delete user from DB", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-del.yaml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-del", 2*time.Minute)
			})
			It("deploys job to test deleted user to connect to cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-user-fail.yaml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-user-fail", 2*time.Minute)
			})
			It("delete user secret", func() {
				err := databases.GitOpsDeleteUser(monitor, user, k8sClient, gitClient)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("zitadel", func() {
		When("bootstraping", func() {
			It("succeeds when creating the initial "+zitadelFile, func() {
				helpers_test.LocalToRemoteFile(zitadelctlGitops, zitadelFile, "./templates/zitadel.yml", os.Getenv)
			})
			It("generates missing secrets successfully", func() {
				configureSession, err := gexec.Start(zitadelctlGitops("configure"), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(configureSession, 2*time.Minute, 1*time.Second).Should(gexec.Exit(0))
			})
			It("deploys zitadel with 1 node", func() {
				session, err := gexec.Start(zitadelctlGitops("takeoff"), os.Stdout, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 2*time.Minute, 5*time.Second).Should(gexec.Exit(0))
			})
			It("runs migrations", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=cockroachdb-cluster-migration-ensure", 15*time.Minute)
			})
			It("runs setup", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=zitadel-setup-ensure", 15*time.Minute)
			})
			It("runs 1 zitadel pod", func() {
				AwaitReadyPods(zitadelNamespace, "app.kubernetes.io/name=zitadel", 1, 2*time.Minute)
			})
			It("flyway migrations should have run a defined version", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-root-flyway.yaml")
				Expect(err).ToNot(HaveOccurred())
				ApplyFile(bytes)
				AwaitCompletedPod(zitadelNamespace, "job-name=cockroachdb-connect-flyway", 1*time.Minute)

				logs := GetLogsOfPod(zitadelNamespace, "job-name=cockroachdb-connect-flyway")

				DeleteFile(bytes)

				outLines := strings.Split(logs, "\n")
				versions := outLines[:len(outLines)-1]
				latestVersion := helpers_test.LastVersionOfMigrations("../../migrations/cockroach")

				Ω(versions[len(versions)-1]).Should(Equal("1." + strconv.Itoa(latestVersion)))

			})
		})
	})

	Context("backup", func() {

		When("backup is defined", func() {
			It("uses defined backups", func() {
				backups, err := ioutil.ReadFile("./templates/backups.yml")
				Expect(err).ToNot(HaveOccurred())

				session, err := gexec.Start(zitadelctlGitops("file", "patch", databaseFile, "database.spec.backups", "--exact", "--value", string(backups)), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))
			})
			It("applies secrets for GCS backup", func() {
				session, err := gexec.Start(zitadelctlGitops("writesecret", "database."+gcsBackupName+".serviceaccountjson.encrypted", "--file", backupSAJson), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))
			})
			It("applies secrets for S3 backup", func() {
				session, err := gexec.Start(zitadelctlGitops("writesecret", "database."+s3BackupName+".accesskeyid.encrypted", "--file", backupAKID), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))

				session, err = gexec.Start(zitadelctlGitops("writesecret", "database."+s3BackupName+".secretaccesskey.encrypted", "--file", backupSAK), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))
			})
			It("applies defined cronjobs for backups", func() {
				AwaitCronJobScheduled(zitadelNamespace, "backup-"+s3BackupName, "0 0 1 1 *", 5*time.Minute)
				AwaitCronJobScheduled(zitadelNamespace, "backup-"+gcsBackupName, "0 0 1 1 *", 5*time.Minute)
			})
		})

		When("instant-backup", func() {
			It("cleanup backups if already existing GCS", func() {
				err := backup.DeleteGCSFolder(backupSAJson, gcsBackupBucket, filepath.Join(gcsBackupName, backupName))
				Expect(err).ToNot(HaveOccurred())
			})
			It("cleanup backups if already existing S3", func() {
				akid, err := ioutil.ReadFile(backupAKID)
				Expect(err).ToNot(HaveOccurred())
				sak, err := ioutil.ReadFile(backupSAK)
				Expect(err).ToNot(HaveOccurred())
				err = backup.DeleteS3Folder(s3Endpoint, string(akid), string(sak), s3BackupBucket, filepath.Join(s3BackupName, backupName), s3Region)
				Expect(err).ToNot(HaveOccurred())
			})
			It("has no existing backup in GCS", func() {
				backups, err := backup.ListGCSFolders(backupSAJson, gcsBackupBucket, gcsBackupName)
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).ToNot(ContainElement(backupName))
			})
			It("has no existing backup in S3", func() {
				akid, err := ioutil.ReadFile(backupAKID)
				Expect(err).ToNot(HaveOccurred())
				sak, err := ioutil.ReadFile(backupSAK)
				Expect(err).ToNot(HaveOccurred())
				backups, err := backup.ListS3Folders(s3Endpoint, string(akid), string(sak), s3BackupBucket, s3BackupName, s3Region)
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).ToNot(ContainElement(backupName))
			})
			It("starts command to backup data", func() {
				session, err := gexec.Start(zitadelctlGitops("backup", "--backup", backupName), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 8*time.Minute, 5*time.Second).Should(gexec.Exit(0))
			})
			/*commented as command cleans up completed pods
			It("deploys job to backup data to gcs", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+gcsBackupName, 8*time.Minute)
			})
			It("deploys job to backup data to s3", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+s3BackupName, 8*time.Minute)
			})*/
			It("created a backup on the GCS bucket", func() {
				backups, err := backup.ListGCSFolders(backupSAJson, gcsBackupBucket, gcsBackupName)
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).To(ContainElement(backupName))
			})
			It("created a backup on the S3 bucket", func() {
				akid, err := ioutil.ReadFile(backupAKID)
				Expect(err).ToNot(HaveOccurred())
				sak, err := ioutil.ReadFile(backupSAK)
				Expect(err).ToNot(HaveOccurred())
				backups, err := backup.ListS3Folders(s3Endpoint, string(akid), string(sak), s3BackupBucket, s3BackupName, s3Region)
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).To(ContainElement(backupName))
			})
		})
		When("instant-backup of changed data", func() {
			It("cleanup backups if already existing GCS", func() {
				err := backup.DeleteGCSFolder(backupSAJson, gcsBackupBucket, filepath.Join(gcsBackupName, changedBackupName))
				Expect(err).ToNot(HaveOccurred())
			})
			It("cleanup backups if already existing S3", func() {
				akid, err := ioutil.ReadFile(backupAKID)
				Expect(err).ToNot(HaveOccurred())
				sak, err := ioutil.ReadFile(backupSAK)
				Expect(err).ToNot(HaveOccurred())
				err = backup.DeleteS3Folder(s3Endpoint, string(akid), string(sak), s3BackupBucket, filepath.Join(s3BackupName, changedBackupName), s3Region)
				Expect(err).ToNot(HaveOccurred())
			})
			It("has no existing backup in GCS", func() {
				backups, err := backup.ListGCSFolders(backupSAJson, gcsBackupBucket, gcsBackupName)
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).ToNot(ContainElement(changedBackupName))
			})
			It("has no existing backup in S3", func() {
				akid, err := ioutil.ReadFile(backupAKID)
				Expect(err).ToNot(HaveOccurred())
				sak, err := ioutil.ReadFile(backupSAK)
				Expect(err).ToNot(HaveOccurred())
				backups, err := backup.ListS3Folders(s3Endpoint, string(akid), string(sak), s3BackupBucket, s3BackupName, s3Region)
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).ToNot(ContainElement(changedBackupName))
			})
			It("add user to DB", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-add.yaml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-add", 2*time.Minute)
			})
			It("generate certificate for user", func() {
				err := databases.GitOpsAddUser(monitor, user, k8sClient, gitClient)
				Expect(err).ToNot(HaveOccurred())
				AwaitSecret(zitadelNamespace, "cockroachdb.client."+user, []string{"ca.crt", "client." + user + ".crt", "client." + user + ".key"}, 1*time.Minute)
			})
			It("deploys job to test added user to connect to cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-user.yaml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-user", 2*time.Minute)
			})
			It("starts command to backup changed data", func() {
				session, err := gexec.Start(zitadelctlGitops("backup", "--backup", changedBackupName), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 8*time.Minute, 5*time.Second).Should(gexec.Exit(0))
			})
			/*commented as command cleans up completed pods
			It("deploys job to backup changed data to GCS", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+gcsBackupName, 8*time.Minute)
			})
			It("deploys job to backup changed data to S3", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+s3BackupName, 8*time.Minute)
			})*/
			It("created a backup on the GCS bucket", func() {
				backups, err := backup.ListGCSFolders(backupSAJson, gcsBackupBucket, gcsBackupName)
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).To(ContainElement(changedBackupName))
			})
			It("created a backup on the S3 bucket", func() {
				akid, err := ioutil.ReadFile(backupAKID)
				Expect(err).ToNot(HaveOccurred())
				sak, err := ioutil.ReadFile(backupSAK)
				Expect(err).ToNot(HaveOccurred())
				backups, err := backup.ListS3Folders(s3Endpoint, string(akid), string(sak), s3BackupBucket, s3BackupName, s3Region)
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).To(ContainElement(changedBackupName))
			})
		})
		When("backuplist", func() {
			It("starts command to list backups", func() {
				session, err := gexec.Start(zitadelctlGitops("backuplist"), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 8*time.Minute, 5*time.Second).Should(gexec.Exit(0))

				backups := strings.Split(string(session.Out.Contents()), "\n")
				expectedBackups := []string{
					gcsBackupName + "." + backupName,
					s3BackupName + "." + backupName,
					gcsBackupName + "." + changedBackupName,
					s3BackupName + "." + changedBackupName,
				}
				Ω(backups).To(ContainElements(expectedBackups))
			})
		})

		When("restore", func() {
			It("starts command to restore data", func() {
				session, err := gexec.Start(zitadelctlGitops("restore", "--backup", gcsBackupName+"."+changedBackupName), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 8*time.Minute, 5*time.Second).Should(gexec.Exit(0))
			})
			It("deploys job to restore data", func() {
				// commented as command cleanup deletes job
				//AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+gcsBackupName+"-restore", 10*time.Minute)
				AwaitReadyPods(zitadelNamespace, "app.kubernetes.io/name=cockroachdb", 3, 3*time.Minute)
			})
			/*
				example to test db functionality
				It("generate certificate for user", func() {
					err := databases.GitOpsAddUser(monitor, user, k8sClient, gitClient)
					Expect(err).ToNot(HaveOccurred())
					AwaitSecret(zitadelNamespace, "cockroachdb.client."+user, []string{"ca.crt", "client." + user + ".crt", "client." + user + ".key"}, 1*time.Minute)
				})
				It("deploys job to test added user to connect to cockroach", func() {
					bytes, err := ioutil.ReadFile("./templates/cockroachdb-user.yaml")
					Expect(err).ToNot(HaveOccurred())

					AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-user", 2*time.Minute)
				})
			*/
		})
	})
	/*
		Context("cleanup", func() {
			When("destroy is called", func() {
				It("succeeds with destroy command", func() {
					session, err := gexec.Start(zitadelctlGitops("destroy"), GinkgoWriter, GinkgoWriter)
					Expect(err).ToNot(HaveOccurred())
					Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))
				})
			})
		})*/
})
