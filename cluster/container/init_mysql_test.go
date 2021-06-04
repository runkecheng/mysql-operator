package container

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	"github.com/zhyass/mysql-operator/utils"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	core "k8s.io/api/core/v1"
)

var (
	initMysqlMysqlCluster = mysqlv1.Cluster{
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
	testInitMysqlCluster = cluster.Cluster{
		Cluster: &initMysqlMysqlCluster,
	}
	initMysqlVolumeMounts = []core.VolumeMount{
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
	optFalse      = false
	optTrue       = true
	sctName       = "sample-secret"
	initMysqlEnvs = []core.EnvVar{
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
	initMysqlCase = EnsureContainer("init-mysql", &testInitMysqlCluster)
)

func TestGetInitMysqlName(t *testing.T) {
	assert.Equal(t, "init-mysql", initMysqlCase.Name)
}
func TestGetInitMysqlImage(t *testing.T) {
	assert.Equal(t, "percona/percona-server:5.7.33", initMysqlCase.Image)
}
func TestGetInitMysqlCommand(t *testing.T) {
	assert.Nil(t, initMysqlCase.Command)
}
func TestGetInitMysqlEnvVar(t *testing.T) {
	//base env
	{
		assert.Equal(t, initMysqlEnvs, initMysqlCase.Env)
	}
	//initTokuDB
	{
		testToKuDBMysqlCluster := initMysqlMysqlCluster
		testToKuDBMysqlCluster.Spec.MysqlOpts.InitTokuDB = true
		testTokuDBCluster := cluster.Cluster{
			Cluster: &testToKuDBMysqlCluster,
		}
		tokudbCase := EnsureContainer("init-mysql", &testTokuDBCluster)
		testEnv := append(initMysqlEnvs, core.EnvVar{
			Name:  "INIT_TOKUDB",
			Value: "1",
		})
		assert.Equal(t, testEnv, tokudbCase.Env)
	}
}
func TestGetInitMysqlLifecycle(t *testing.T) {
	assert.Nil(t, initMysqlCase.Lifecycle)
}
func TestGetInitMysqlResources(t *testing.T) {
	assert.Equal(t, core.ResourceRequirements{
		Limits:   nil,
		Requests: nil,
	}, initMysqlCase.Resources)
}
func TestGetInitMysqlPorts(t *testing.T) {
	assert.Nil(t, initMysqlCase.Ports)
}
func TestGetInitMysqlLivenessProbe(t *testing.T) {
	assert.Nil(t, initMysqlCase.LivenessProbe)
}
func TestGetInitMysqlReadinessProbe(t *testing.T) {
	assert.Nil(t, initMysqlCase.ReadinessProbe)
}
func TestGetInitMysqlVolumeMounts(t *testing.T) {
	assert.Equal(t, initMysqlVolumeMounts, initMysqlCase.VolumeMounts)
}
