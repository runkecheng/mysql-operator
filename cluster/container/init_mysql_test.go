package container_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	"github.com/zhyass/mysql-operator/cluster/container"
	"github.com/zhyass/mysql-operator/utils"

	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Container", func() {
	Describe("init-mysql", func() {
		volumeMounts := []core.VolumeMount{
			{
				Name:      utils.ConfVolumeName,
				MountPath: utils.ConfVolumeMountPath,
			},
			{
				Name:      utils.DataVolumeName,
				MountPath: utils.DataVolumeMountPath,
			},
			{
				Name:      utils.LogsVolumeName,
				MountPath: utils.LogsVolumeMountPath,
			},
			{
				Name:      utils.InitFileVolumeName,
				MountPath: utils.InitFileVolumeMountPath,
			},
		}

		mysqlCluster := mysqlv1.Cluster{
			ObjectMeta: v1.ObjectMeta{
				Name: "sample",
			},
			Spec: mysqlv1.ClusterSpec{
				PodSpec: mysqlv1.PodSpec{
					Resources: core.ResourceRequirements{
						Limits:   nil,
						Requests: nil,
					},
				},
				MysqlVersion: "5.7",
				MysqlOpts: mysqlv1.MysqlOpts{
					InitTokuDB: false,
				},
			},
		}
		initMysqlCluster := cluster.Cluster{
			Cluster: &mysqlCluster,
		}
		InitMysqlCase := container.EnsureContainer("init-mysql", &initMysqlCluster)
		It("should equal", func() {
			//getName
			Expect("init-mysql").To(Equal(InitMysqlCase.Name))
			//getImage
			Expect("percona/percona-server:5.7.33").To(Equal(InitMysqlCase.Image))
			//getResources
			Expect(core.ResourceRequirements{
				Limits:   nil,
				Requests: nil,
			}).To(Equal(InitMysqlCase.Resources))
			//getVolumeMounts
			Expect(volumeMounts).To(Equal(InitMysqlCase.VolumeMounts))
		})
		It("should be nil", func() {
			//getCommand
			Expect(InitMysqlCase.Command).To(BeNil())
			//getLifecycle
			Expect(InitMysqlCase.Lifecycle).To(BeNil())
			//getPorts
			Expect(InitMysqlCase.Ports).To(BeNil())
			//getLivenessProbe
			Expect(InitMysqlCase.LivenessProbe).To(BeNil())
			//getReadinessProbe
			Expect(InitMysqlCase.ReadinessProbe).To(BeNil())
		})

		Describe("getEnvVars", func() {
			optFalse := false
			optTrue := true
			sctName := "sample-secret"
			envs := []core.EnvVar{
				{
					Name:  "MYSQL_ALLOW_EMPTY_PASSWORD",
					Value: "yes",
				},
				{
					Name:  "MYSQL_ROOT_HOST",
					Value: "127.0.0.1",
				},
				{
					Name:  "MYSQL_INIT_ONLY",
					Value: "1",
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
					Name: "MYSQL_DATABASE",
					ValueFrom: &core.EnvVarSource{
						SecretKeyRef: &core.SecretKeySelector{
							LocalObjectReference: core.LocalObjectReference{
								Name: sctName,
							},
							Key:      "mysql-database",
							Optional: &optTrue,
						},
					},
				},
				{
					Name: "MYSQL_USER",
					ValueFrom: &core.EnvVarSource{
						SecretKeyRef: &core.SecretKeySelector{
							LocalObjectReference: core.LocalObjectReference{
								Name: sctName,
							},
							Key:      "mysql-user",
							Optional: &optTrue,
						},
					},
				},
				{
					Name: "MYSQL_PASSWORD",
					ValueFrom: &core.EnvVarSource{
						SecretKeyRef: &core.SecretKeySelector{
							LocalObjectReference: core.LocalObjectReference{
								Name: sctName,
							},
							Key:      "mysql-password",
							Optional: &optTrue,
						},
					},
				},
			}
			It("should equal", func() {
				Expect(envs).To(Equal(InitMysqlCase.Env))
			})

			Context("when init Tokudb", func() {
				testToKuDBMysqlCluster := mysqlCluster
				testToKuDBMysqlCluster.Spec.MysqlOpts.InitTokuDB = true
				testTokuDBCluster := cluster.Cluster{
					Cluster: &testToKuDBMysqlCluster,
				}
				tokudbCase := container.EnsureContainer("init-mysql", &testTokuDBCluster)
				testEnv := envs
				testEnv = append(envs, core.EnvVar{
					Name:  "INIT_TOKUDB",
					Value: "1",
				})
				It("should equal", func() {
					Expect(testEnv).To(Equal(tokudbCase.Env))
				})
			})
		})
	})
})
