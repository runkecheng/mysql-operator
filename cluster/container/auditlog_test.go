package container

import (
	"testing"

	"github.com/stretchr/testify/assert"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"

	core "k8s.io/api/core/v1"
)

var (
	auditlogMysqlCluster = mysqlv1.Cluster{
		Spec: mysqlv1.ClusterSpec{
			PodSpec: mysqlv1.PodSpec{
				BusyboxImage: "busybox",
				Resources: core.ResourceRequirements{
					Limits:   nil,
					Requests: nil,
				},
			},
		},
	}
	testAuditlogCluster = cluster.Cluster{
		Cluster: &auditlogMysqlCluster,
	}
	auditLogCase         = EnsureContainer("auditlog", &testAuditlogCluster)
	auditlogCommand      = []string{"tail", "-f", "/var/log/mysql" + "/mysql-audit.log"}
	auditlogVolumeMounts = []core.VolumeMount{
		{
			Name:      "logs",
			MountPath: "/var/log/mysql",
		},
	}
)

func TestGetAuditlogName(t *testing.T) {
	assert.Equal(t, "auditlog", auditLogCase.Name)
}
func TestGetAuditlogImage(t *testing.T) {
	assert.Equal(t, "busybox", auditLogCase.Image)
}
func TestGetAuditlogCommand(t *testing.T) {
	assert.Equal(t, auditlogCommand, auditLogCase.Command)
}
func TestGetAuditlogEnvVar(t *testing.T) {
	assert.Nil(t, auditLogCase.Env)
}
func TestGetAuditlogLifecycle(t *testing.T) {
	assert.Nil(t, auditLogCase.Lifecycle)
}
func TestGetAuditlogResources(t *testing.T) {
	assert.Equal(t, core.ResourceRequirements{
		Limits:   nil,
		Requests: nil,
	}, auditLogCase.Resources)
}
func TestGetAuditlogPorts(t *testing.T) {
	assert.Nil(t, auditLogCase.Ports)
}
func TestGetAuditlogLivenessProbe(t *testing.T) {
	assert.Nil(t, auditLogCase.LivenessProbe)
}
func TestGetAuditlogReadinessProbe(t *testing.T) {
	assert.Nil(t, auditLogCase.ReadinessProbe)
}
func TestGetAuditlogVolumeMounts(t *testing.T) {
	assert.Equal(t, auditlogVolumeMounts, auditLogCase.VolumeMounts)
}
