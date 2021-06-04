package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"

	core "k8s.io/api/core/v1"
)

var (
	slowlogMysqlCluster = mysqlv1.Cluster{
		Spec: mysqlv1.ClusterSpec{
			PodSpec: mysqlv1.PodSpec{
				SidecarImage: "sidecar image",
				Resources: core.ResourceRequirements{
					Limits:   nil,
					Requests: nil,
				},
			},
		},
	}
	testSlowlogCluster = cluster.Cluster{
		Cluster: &slowlogMysqlCluster,
	}
	slowlogCase = EnsureContainer("slowlog", &testSlowlogCluster)
)

func TestGetSlowlogName(t *testing.T) {
	assert.Equal(t, "slowlog", slowlogCase.Name)
}
func TestGetSlowlogImage(t *testing.T) {
	assert.Equal(t, "sidecar image", slowlogCase.Image)
}
func TestGetSlowlogCommand(t *testing.T) {
	command := []string{"tail", "-f", "/var/log/mysql" + "/mysql-slow.log"}
	assert.Equal(t, command, slowlogCase.Command)
}
func TestGetSlowlogEnvVar(t *testing.T) {
	assert.Nil(t, slowlogCase.Env)
}
func TestGetSlowlogLifecycle(t *testing.T) {
	assert.Nil(t, slowlogCase.Lifecycle)
}
func TestGetSlowlogResources(t *testing.T) {
	assert.Equal(t, core.ResourceRequirements{
		Limits:   nil,
		Requests: nil,
	}, slowlogCase.Resources)
}
func TestGetSlowlogPorts(t *testing.T) {
	assert.Nil(t, slowlogCase.Ports)
}
func TestGetSlowlogLivenessProbe(t *testing.T) {
	assert.Nil(t, slowlogCase.LivenessProbe)
}
func TestGetSlowlogReadinessProbe(t *testing.T) {
	assert.Nil(t, slowlogCase.ReadinessProbe)
}
func TestGetSlowlogVolumeMounts(t *testing.T) {
	volumeMounts := []core.VolumeMount{
		{
			Name:      "logs",
			MountPath: "/var/log/mysql",
		},
	}
	assert.Equal(t, volumeMounts, slowlogCase.VolumeMounts)
}
