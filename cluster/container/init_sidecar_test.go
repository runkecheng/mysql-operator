package container_test

import (
	"strconv"

	// . "bou.ke/monkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	"github.com/zhyass/mysql-operator/cluster/container"
	"github.com/zhyass/mysql-operator/utils"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Container", func() {

	Describe("init-sidecar", func() {
		command := []string{"sidecar", "init"}
		var defeatCount int32 = 1
		var electionTimeout int32 = 5
		mysqlCluster := mysqlv1.Cluster{
			ObjectMeta: v1.ObjectMeta{
				Name: "sample",
			},
			Spec: mysqlv1.ClusterSpec{
				PodSpec: mysqlv1.PodSpec{
					SidecarImage: "sidecar image",
					Resources: core.ResourceRequirements{
						Limits:   nil,
						Requests: nil,
					},
				},
				XenonOpts: mysqlv1.XenonOpts{
					AdmitDefeatHearbeatCount: &defeatCount,
					ElectionTimeout:          &electionTimeout,
				},
				MetricsOpts: mysqlv1.MetricsOpts{
					Enabled: false,
				},
				MysqlOpts: mysqlv1.MysqlOpts{
					InitTokuDB: false,
				},
				Persistence: mysqlv1.Persistence{
					Enabled: false,
				},
			},
		}
		initSidecarCluster := cluster.Cluster{
			Cluster: &mysqlCluster,
		}
		testInitSidecar := container.EnsureContainer("init-sidecar", &initSidecarCluster)

		It("should equal init-sidecar", func() {
			//getName
			Expect("init-sidecar").To(Equal(testInitSidecar.Name))
			//getImage
			Expect("sidecar image").To(Equal(testInitSidecar.Image))
			//getCommand
			Expect(command).To(Equal(testInitSidecar.Command))
			//getLifecycle
			Expect(testInitSidecar.Lifecycle).To(BeNil())
			//getResources
			Expect(testInitSidecar.Resources).To(Equal(core.ResourceRequirements{Limits: nil, Requests: nil}))
			//getPorts
			Expect(testInitSidecar.Ports).To(BeNil())
			//getLivenessProbe
			Expect(testInitSidecar.LivenessProbe).To(BeNil())
			//getReadinessProbe
			Expect(testInitSidecar.ReadinessProbe).To(BeNil())
		})
		//getEnvVars
		Describe("getEnvVars", func() {

			sctName := "sample-secret"
			optFalse := false
			optTrue := true
			envs := []core.EnvVar{
				{
					Name: "POD_HOSTNAME",
					ValueFrom: &core.EnvVarSource{
						FieldRef: &core.ObjectFieldSelector{
							APIVersion: "v1",
							FieldPath:  "metadata.name",
						},
					},
				},
				{
					Name:  "NAMESPACE",
					Value: initSidecarCluster.Namespace,
				},
				{
					Name:  "SERVICE_NAME",
					Value: "sample-mysql",
				},
				{
					Name:  "ADMIT_DEFEAT_HEARBEAT_COUNT",
					Value: strconv.Itoa(int(*initSidecarCluster.Spec.XenonOpts.AdmitDefeatHearbeatCount)),
				},
				{
					Name:  "ELECTION_TIMEOUT",
					Value: strconv.Itoa(int(*initSidecarCluster.Spec.XenonOpts.ElectionTimeout)),
				},
				{
					Name:  "MY_MYSQL_VERSION",
					Value: "5.7.33",
				},
				{
					Name: "MYSQL_ROOT_PASSWORD",
					ValueFrom: &core.EnvVarSource{
						SecretKeyRef: &core.SecretKeySelector{
							LocalObjectReference: core.LocalObjectReference{
								Name: sctName,
							},
							Key:      "root-password",
							Optional: &optFalse,
						},
					},
				},
				{
					Name: "MYSQL_REPL_USER",
					ValueFrom: &core.EnvVarSource{
						SecretKeyRef: &core.SecretKeySelector{
							LocalObjectReference: core.LocalObjectReference{
								Name: sctName,
							},
							Key:      "replication-user",
							Optional: &optTrue,
						},
					},
				},
				{
					Name: "MYSQL_REPL_PASSWORD",
					ValueFrom: &core.EnvVarSource{
						SecretKeyRef: &core.SecretKeySelector{
							LocalObjectReference: core.LocalObjectReference{
								Name: sctName,
							},
							Key:      "replication-password",
							Optional: &optTrue,
						},
					},
				},
			}
			It("should equal", func() {
				Expect(envs).To(Equal(testInitSidecar.Env))
			})

			Context("when init tukudb", func() {
				envTukudbCluster := mysqlCluster
				envTukudbCluster.Spec.MysqlOpts.InitTokuDB = true
				testCluster := cluster.Cluster{
					Cluster: &envTukudbCluster,
				}
				tukudbEnv := make([]core.EnvVar, 9, 10)
				copy(tukudbEnv, envs)
				tukudbEnv = append(tukudbEnv, core.EnvVar{
					Name:  "INIT_TOKUDB",
					Value: "1",
				})
				tukudbCase := container.EnsureContainer("init-sidecar", &testCluster)

				It("should equal", func() {
					Expect(tukudbEnv).To(Equal(tukudbCase.Env))
				})
			})
			Context("when enable metrics", func() {
				envMetricsCluster := mysqlCluster
				envMetricsCluster.Spec.MetricsOpts.Enabled = true
				testCluster := cluster.Cluster{
					Cluster: &envMetricsCluster,
				}
				metricsEnv := make([]core.EnvVar, 9, 11)
				copy(metricsEnv, envs)
				metricsEnv = append(metricsEnv,
					core.EnvVar{
						Name: "METRICS_USER",
						ValueFrom: &core.EnvVarSource{
							SecretKeyRef: &core.SecretKeySelector{
								LocalObjectReference: core.LocalObjectReference{
									Name: sctName,
								},
								Key:      "metrics-user",
								Optional: &optTrue,
							},
						},
					},
					core.EnvVar{
						Name: "METRICS_PASSWORD",
						ValueFrom: &core.EnvVarSource{
							SecretKeyRef: &core.SecretKeySelector{
								LocalObjectReference: core.LocalObjectReference{
									Name: sctName,
								},
								Key:      "metrics-password",
								Optional: &optTrue,
							},
						},
					},
				)
				metricsCase := container.EnsureContainer("init-sidecar", &testCluster)
				It("should equal", func() {
					Expect(metricsEnv).To(Equal(metricsCase.Env))
				})
			})
		})

		//getVolumeMounts
		Describe("getVolumeMounts", func() {
			volumeMounts := []core.VolumeMount{
				{
					Name:      utils.ConfVolumeName,
					MountPath: utils.ConfVolumeMountPath,
				},
				{
					Name:      utils.ConfMapVolumeName,
					MountPath: utils.ConfMapVolumeMountPath,
				},
				{
					Name:      utils.ScriptsVolumeName,
					MountPath: utils.ScriptsVolumeMountPath,
				},
				{
					Name:      utils.XenonVolumeName,
					MountPath: utils.XenonVolumeMountPath,
				},
				{
					Name:      utils.InitFileVolumeName,
					MountPath: utils.InitFileVolumeMountPath,
				},
			}
			tokuDBVolumeMounts := make([]core.VolumeMount, 5, 6)
			persistenceVolumeMounts := make([]core.VolumeMount, 5, 6)
			copy(tokuDBVolumeMounts, volumeMounts)
			copy(persistenceVolumeMounts, volumeMounts)
			tokuDBVolumeMounts = append(tokuDBVolumeMounts, core.VolumeMount{
				Name:      utils.SysVolumeName,
				MountPath: utils.SysVolumeMountPath,
			})
			persistenceVolumeMounts = append(persistenceVolumeMounts, core.VolumeMount{
				Name:      utils.DataVolumeName,
				MountPath: utils.DataVolumeMountPath,
			})

			It("should equal", func() {
				Expect(volumeMounts).To(Equal(testInitSidecar.VolumeMounts))
			})
			Context("when init tokuDB", func() {
				tokudbCluster := mysqlCluster
				tokudbCluster.Spec.MysqlOpts.InitTokuDB = true
				testCluster := cluster.Cluster{
					Cluster: &tokudbCluster,
				}
				initTokuDBCase := container.EnsureContainer("init-sidecar", &testCluster)
				It("should equal", func() {
					Expect(tokuDBVolumeMounts).To(Equal(initTokuDBCase.VolumeMounts))
				})
			})
			Context("when enable persistence", func() {
				volumePersistenceCluster := mysqlCluster
				volumePersistenceCluster.Spec.Persistence.Enabled = true
				testCluster := cluster.Cluster{
					Cluster: &volumePersistenceCluster,
				}
				persistenceCase := container.EnsureContainer("init-sidecar", &testCluster)
				It("should equal", func() {
					Expect(persistenceVolumeMounts).To(Equal(persistenceCase.VolumeMounts))
				})

			})
		})
	})

})
