package container

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	"github.com/zhyass/mysql-operator/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	core "k8s.io/api/core/v1"
)

var (
	defeatCount             int32 = 1
	electionTimeout         int32 = 5
	initSidecarMysqlCluster       = mysqlv1.Cluster{
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
	testInitSidecarCluster = cluster.Cluster{
		Cluster: &initSidecarMysqlCluster,
	}
	initSidecarEnvs = []core.EnvVar{
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
			Value: testInitSidecarCluster.Namespace,
		},
		{
			Name:  "SERVICE_NAME",
			Value: "sample-mysql",
		},
		{
			Name:  "ADMIT_DEFEAT_HEARBEAT_COUNT",
			Value: strconv.Itoa(int(*testInitSidecarCluster.Spec.XenonOpts.AdmitDefeatHearbeatCount)),
		},
		{
			Name:  "ELECTION_TIMEOUT",
			Value: strconv.Itoa(int(*testInitSidecarCluster.Spec.XenonOpts.ElectionTimeout)),
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
		{
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
		{
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
	}
	initsidecarVolumeMounts = []core.VolumeMount{
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
	initSidecarCase = EnsureContainer("init-sidecar", &testInitSidecarCluster)
)

func TestGetInitSidecarName(t *testing.T) {
	assert.Equal(t, "init-sidecar", initSidecarCase.Name)
}
func TestGetInitSidecarImage(t *testing.T) {
	assert.Equal(t, "sidecar image", initSidecarCase.Image)
}
func TestGetInitSidecarCommand(t *testing.T) {
	command := []string{"sidecar", "init"}
	assert.Equal(t, command, initSidecarCase.Command)
}
func TestGetInitSidecarEnvVar(t *testing.T) {
	//base env
	{
		assert.Equal(t, initSidecarEnvs, initSidecarCase.Env)
	}
	//initTokuDB
	{
		testToKuDBMysqlCluster := initSidecarMysqlCluster
		testToKuDBMysqlCluster.Spec.MysqlOpts.InitTokuDB = true
		testTokuDBCluster := cluster.Cluster{
			Cluster: &testToKuDBMysqlCluster,
		}
		tokudbCase := EnsureContainer("init-sidecar", &testTokuDBCluster)
		testTokuDBEnv := make([]core.EnvVar, 11, 12)
		copy(testTokuDBEnv, initSidecarEnvs)
		testTokuDBEnv = append(testTokuDBEnv, core.EnvVar{
			Name:  "INIT_TOKUDB",
			Value: "1",
		})
		assert.Equal(t, testTokuDBEnv, tokudbCase.Env)
	}
}
func TestGetInitSidecarLifecycle(t *testing.T) {
	assert.Nil(t, initSidecarCase.Lifecycle)
}
func TestGetInitSidecarResources(t *testing.T) {
	assert.Equal(t, core.ResourceRequirements{
		Limits:   nil,
		Requests: nil,
	}, initSidecarCase.Resources)
}
func TestGetInitSidecarPorts(t *testing.T) {
	assert.Nil(t, initSidecarCase.Ports)
}
func TestGetInitSidecarLivenessProbe(t *testing.T) {
	assert.Nil(t, initSidecarCase.LivenessProbe)
}
func TestGetInitSidecarReadinessProbe(t *testing.T) {
	assert.Nil(t, initSidecarCase.ReadinessProbe)
}
func TestGetInitSidecarVolumeMounts(t *testing.T) {
	//base
	{
		assert.Equal(t, initsidecarVolumeMounts, initSidecarCase.VolumeMounts)
	}
	//init tokudb
	{
		testToKuDBMysqlCluster := initSidecarMysqlCluster
		testToKuDBMysqlCluster.Spec.MysqlOpts.InitTokuDB = true
		testTokuDBCluster := cluster.Cluster{
			Cluster: &testToKuDBMysqlCluster,
		}
		tokudbCase := EnsureContainer("init-sidecar", &testTokuDBCluster)
		tokuDBVolumeMounts := make([]core.VolumeMount, 5, 6)
		copy(tokuDBVolumeMounts, initsidecarVolumeMounts)
		tokuDBVolumeMounts = append(tokuDBVolumeMounts, core.VolumeMount{
			Name:      utils.SysVolumeName,
			MountPath: utils.SysVolumeMountPath,
		})
		assert.Equal(t, tokuDBVolumeMounts, tokudbCase.VolumeMounts)
	}
	//enable persistence
	{
		testPersistenceMysqlCluster := initSidecarMysqlCluster
		testPersistenceMysqlCluster.Spec.Persistence.Enabled = true
		testPersistenceCluster := cluster.Cluster{
			Cluster: &testPersistenceMysqlCluster,
		}
		persistenceCase := EnsureContainer("init-sidecar", &testPersistenceCluster)
		persistenceVolumeMounts := make([]core.VolumeMount, 5, 6)
		copy(persistenceVolumeMounts, initsidecarVolumeMounts)
		persistenceVolumeMounts = append(persistenceVolumeMounts, core.VolumeMount{
			Name:      utils.DataVolumeName,
			MountPath: utils.DataVolumeMountPath,
		})
		assert.Equal(t, persistenceVolumeMounts, persistenceCase.VolumeMounts)

	}
}
