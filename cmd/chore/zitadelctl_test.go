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
	"github.com/caos/zitadel/pkg/backup"
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
		cleanupAfterEnv     = envPrefix + "CLEANUP_AFTER"
		reuseOrbEnv         = envPrefix + "REUSE_ORB"
	)
	var (
		tag, sha, orbconfigPath, workfolder, user, backupSAJson, backupAKID, backupSAK              string
		gcsBackupName, gcsBackupBucket, s3BackupName, s3BackupBucket, backupName, changedBackupName string
		monitor                                                                                     mntr.Monitor
		k8sClient                                                                                   kubernetes.ClientInt
		gitClient                                                                                   *git.Client
		kubectl                                                                                     kubectlCmd
		zitadelctlGitops                                                                            zitadelctlGitopsCmd
		GetLogsOfPod                                                                                getLogsOfPod
		ApplyFile                                                                                   applyFile
		DeleteFile                                                                                  deleteFile
		AwaitCompletedPodFromJob                                                                    awaitCompletedPodFromJob
		AwaitCompletedPod                                                                           awaitCompletedPod
		AwaitReadyPods                                                                              awaitReadyPods
		AwaitSecret                                                                                 awaitSecret
		AwaitCronJobScheduled                                                                       awaitCronJobScheduled
	)
	BeforeSuite(func() {
		backupName = "e2e-test"
		changedBackupName = "e2e-test2"
		gcsBackupName = "bucket"
		gcsBackupBucket = "caos-zitadel-e2e-backup"
		s3BackupName = "csbucket"
		s3BackupBucket = "caos-zitadel-e2e-backup"
		workfolder = "./artifacts"
		kubeconfigPath := filepath.Join(workfolder, "kubeconfig")
		orbconfigPath = filepath.Join(workfolder, "orbconfig")
		tag = prefixedEnv(tagEnv)
		sha = prefixedEnv(shaEnv)
		user = prefixedEnv(userEnv)
		backupSAJson = prefixedEnv(backupSAJSONPathEnv)
		backupAKID = prefixedEnv(backupAKIDPathEnv)
		backupSAK = prefixedEnv(backupSAKPathEnv)
		zitadelctlGitops = zitadelctlGitopsFunc(orbconfigPath)
		kubectl = kubectlCmdFunc(kubeconfigPath)
		monitor = mntr.Monitor{
			OnInfo:         mntr.LogMessage,
			OnChange:       mntr.LogMessage,
			OnError:        mntr.LogError,
			OnRecoverPanic: mntr.LogPanic,
		}
		ApplyFile = applyFileFunc(kubectl)
		DeleteFile = deleteFileFunc(kubectl)
		AwaitCompletedPod = awaitCompletedPodFunc(kubectl)
		AwaitCompletedPodFromJob = awaitCompletedPodFromJobFunc(kubectl)
		AwaitReadyPods = awaitReadyPodsFunc(kubectl)
		AwaitSecret = awaitSecretFunc(kubectl)
		AwaitCronJobScheduled = awaitCronJobScheduledFunc(kubectl)
		GetLogsOfPod = getLogsOfPodFunc(kubectl)

		orbconfig, err := orb.ParseOrbConfig(orbconfigPath)
		Expect(err).ToNot(HaveOccurred())
		kubeconfig, _ := ioutil.ReadFile(kubeconfigPath)
		kubeconfigStr := string(kubeconfig)
		k8sClientT, _ := kubernetes.NewK8sClient(monitor, &kubeconfigStr)
		k8sClient = k8sClientT
		gitClient = git.New(context.Background(), monitor, "zitadelctl", "test@orbos.ch")
		err = cli.InitRepo(orbconfig, gitClient)
		Expect(err).ToNot(HaveOccurred())

		Expect(tag).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+tagEnv))
		Expect(sha).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+shaEnv))
		Expect(user).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+userEnv))
		Expect(backupSAJson).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+backupSAJSONPathEnv))
		Expect(backupAKID).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+backupAKIDPathEnv))
		Expect(backupSAK).ToNot(BeEmpty(), fmt.Sprintf("environment variable %s is required", envPrefix+backupSAKPathEnv))
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
				session, err := gexec.Start(zitadelctlGitops("file", "patch", "database.yml", "database.spec.replicaCount", "--exact", "--value", strconv.Itoa(count)), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

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
			It("succeeds when creating the initial zitadel.yml", func() {
				localToRemoteFile(zitadelctlGitops, "zitadel.yml", "./templates/zitadel.yml", os.Getenv)
			})
			It("generates missing secrets successfully", func() {
				configureSession, err := gexec.Start(zitadelctlGitops("configure"), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(configureSession, 2*time.Minute, 1*time.Second).Should(gexec.Exit(0))
			})
			It("deploys zitadel with 1 node", func() {
				session, err := gexec.Start(zitadelctlGitops("takeoff"), os.Stdout, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))
			})
			It("runs migrations", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=cockroachdb-cluster-migration-ensure", 10*time.Minute)
			})
			It("runs setup", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=zitadel-setup-ensure", 10*time.Minute)
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
				latestVersion := lastVersionOfMigrations("../../migrations/cockroach")

				Ω(versions[len(versions)-1]).Should(Equal("1." + strconv.Itoa(latestVersion)))

			})
		})
	})

	Context("backup", func() {
		When("backup is defined", func() {
			It("applies job with declared information", func() {
				backups, err := ioutil.ReadFile("./templates/backups.yml")
				Expect(err).ToNot(HaveOccurred())

				session, err := gexec.Start(zitadelctlGitops("file", "patch", "database.yml", "database.spec.backups", "--exact", "--value", string(backups)), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

				session, err = gexec.Start(zitadelctlGitops("writesecret", "database."+gcsBackupName+".serviceaccountjson.encrypted", "--file", backupSAJson), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
				AwaitCronJobScheduled(zitadelNamespace, "backup-"+gcsBackupName, "0 0 1 1 *", 5*time.Minute)

				session, err = gexec.Start(zitadelctlGitops("writesecret", "database."+s3BackupName+".accesskeyid.encrypted", "--file", backupAKID), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

				session, err = gexec.Start(zitadelctlGitops("writesecret", "database."+s3BackupName+".secretaccesskey.encrypted", "--file", backupSAK), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
				AwaitCronJobScheduled(zitadelNamespace, "backup-"+s3BackupName, "0 0 1 1 *", 5*time.Minute)
			})
		})

		When("instant-backup", func() {
			It("starts command to backup data", func() {
				session, err := gexec.Start(zitadelctlGitops("backup", "--backup", backupName), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 8*time.Minute).Should(gexec.Exit(0))
			})
			It("deploys job to backup data to gcs", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+gcsBackupName, 8*time.Minute)
			})
			It("deploys job to backup data to s3", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+s3BackupName, 8*time.Minute)
			})
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
				backups, err := backup.ListS3Folders("https://objects.lpg.cloudscale.ch", string(akid), string(sak), s3BackupBucket, s3BackupName, "LPG")
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).To(ContainElement(backupName))
			})
		})
		When("instant-backup of changed data", func() {
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
				Eventually(session, 8*time.Minute).Should(gexec.Exit(0))
			})
			It("deploys job to backup changed data to GCS", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+gcsBackupName, 8*time.Minute)
			})
			It("deploys job to backup changed data to S3", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+s3BackupName, 8*time.Minute)
			})
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
				backups, err := backup.ListS3Folders("https://objects.lpg.cloudscale.ch", string(akid), string(sak), s3BackupBucket, s3BackupName, "LPG")
				Expect(err).ToNot(HaveOccurred())
				Ω(backups).To(ContainElement(changedBackupName))
			})
		})
		When("backuplist", func() {
			It("starts command to list backups", func() {
				session, err := gexec.Start(zitadelctlGitops("backuplist"), GinkgoWriter, GinkgoWriter)
				Expect(err).ToNot(HaveOccurred())
				Eventually(session, 8*time.Minute).Should(gexec.Exit(0))
				buf := make([]byte, 0)
				_, err = session.Out.Read(buf)
				Expect(err).ToNot(HaveOccurred())
				backups := strings.Split(string(buf), "\n")

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
				Eventually(session, 8*time.Minute).Should(gexec.Exit(0))
			})
			It("deploys job to restore data", func() {
				AwaitCompletedPod(zitadelNamespace, "job-name=backup-"+gcsBackupName+"-restore", 10*time.Minute)
				AwaitReadyPods(zitadelNamespace, "app.kubernetes.io/name=cockroachdb", 3, 3*time.Minute)
			})
			It("deploys job to test added user to connect to cockroach", func() {
				bytes, err := ioutil.ReadFile("./templates/cockroachdb-user.yaml")
				Expect(err).ToNot(HaveOccurred())

				AwaitCompletedPodFromJob(bytes, zitadelNamespace, "job-name=cockroachdb-connect-user", 2*time.Minute)
			})
		})
	})
})
