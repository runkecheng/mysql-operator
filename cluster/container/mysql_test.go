package container

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"

	core "k8s.io/api/core/v1"
)

var (
	mysqlMysqlCluster = mysqlv1.Cluster{
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
	testMysqlCluster = cluster.Cluster{
		Cluster: &mysqlMysqlCluster,
	}
	mysqlCase = EnsureContainer("mysql", &testMysqlCluster)
)

func TestGetMysqlName(t *testing.T) {
	assert.Equal(t, "mysql", mysqlCase.Name)
}
func TestGetMysqlImage(t *testing.T) {
	assert.Equal(t, "percona/percona-server:5.7.33", mysqlCase.Image)
}
func TestGetMysqlCommand(t *testing.T) {
	assert.Nil(t, mysqlCase.Command)
}
func TestGetMysqlEnvVar(t *testing.T) {
	//base env
	{
		assert.Nil(t, mysqlCase.Env)
	}
	//initTokuDB
	{
		volumeMounts := []core.EnvVar{
			{
				Name:  "INIT_TOKUDB",
				Value: "1",
			},
		}
		mysqlCluster := mysqlMysqlCluster
		mysqlCluster.Spec.MysqlOpts.InitTokuDB = true
		testCluster := cluster.Cluster{
			Cluster: &mysqlCluster,
		}
		mysqlCase = EnsureContainer("mysql", &testCluster)
		assert.Equal(t, volumeMounts, mysqlCase.Env)
	}
}
func TestGetMysqlLifecycle(t *testing.T) {
	assert.Nil(t, mysqlCase.Lifecycle)
}
func TestGetMysqlResources(t *testing.T) {
	assert.Equal(t, core.ResourceRequirements{
		Limits:   nil,
		Requests: nil,
	}, mysqlCase.Resources)
}
func TestGetMysqlPorts(t *testing.T) {
	port := []core.ContainerPort{
		{
			Name:          "mysql",
			ContainerPort: 3306,
		},
	}
	assert.Equal(t, port, mysqlCase.Ports)
}
func TestGetMysqlLivenessProbe(t *testing.T) {
	livenessProbe := &core.Probe{
		Handler: core.Handler{
			Exec: &core.ExecAction{
				Command: []string{"sh", "-c", "mysqladmin ping -uroot -p${MYSQL_ROOT_PASSWORD}"},
			},
		},
		InitialDelaySeconds: 30,
		TimeoutSeconds:      5,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	assert.Equal(t, livenessProbe, mysqlCase.LivenessProbe)
}
func TestGetMysqlReadinessProbe(t *testing.T) {
	readinessProbe := &core.Probe{
		Handler: core.Handler{
			Exec: &core.ExecAction{
				Command: []string{"sh", "-c", `mysql -uroot -p${MYSQL_ROOT_PASSWORD} -e "SELECT 1"`},
			},
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      1,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	assert.Equal(t, readinessProbe, mysqlCase.ReadinessProbe)
}
func TestGetMysqlVolumeMounts(t *testing.T) {
	volumeMounts := []core.VolumeMount{
		{
			Name:      "conf",
			MountPath: "/etc/mysql",
		},
		{
			Name:      "data",
			MountPath: "/var/lib/mysql",
		},
		{
			Name:      "logs",
			MountPath: "/var/log/mysql",
		},
	}
	assert.Equal(t, volumeMounts, mysqlCase.VolumeMounts)
}
