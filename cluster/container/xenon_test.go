package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	core "k8s.io/api/core/v1"
)

var (
	xenonReplicas     int32 = 1
	xenonMysqlCluster       = mysqlv1.Cluster{
		ObjectMeta: v1.ObjectMeta{
			Name:      "sample",
			Namespace: "default",
		},
		Spec: mysqlv1.ClusterSpec{
			PodSpec: mysqlv1.PodSpec{
				Resources: core.ResourceRequirements{
					Limits:   nil,
					Requests: nil,
				},
			},
			XenonOpts: mysqlv1.XenonOpts{
				Image: "xenon image",
			},
			Replicas: &xenonReplicas,
		},
	}
	testXenonCluster = cluster.Cluster{
		Cluster: &xenonMysqlCluster,
	}
	xenonCase = EnsureContainer("xenon", &testXenonCluster)
)

func TestGetXenonName(t *testing.T) {
	assert.Equal(t, "xenon", xenonCase.Name)
}
func TestGetXenonImage(t *testing.T) {
	assert.Equal(t, "xenon image", xenonCase.Image)
}
func TestGetXenonCommand(t *testing.T) {
	assert.Nil(t, xenonCase.Command)
}
func TestGetXenonEnvVar(t *testing.T) {
	assert.Nil(t, xenonCase.Env)
}
func TestGetXenonLifecycle(t *testing.T) {
	lifecycle := &core.Lifecycle{
		PostStart: &core.Handler{
			Exec: &core.ExecAction{
				Command: []string{"sh", "-c",
					"until (xenoncli xenon ping && xenoncli cluster add sample-mysql-0.sample-mysql.default:8801) > /dev/null 2>&1; do sleep 2; done",
				},
			},
		},
	}
	assert.Equal(t, lifecycle, xenonCase.Lifecycle)
}
func TestGetXenonResources(t *testing.T) {
	assert.Equal(t, core.ResourceRequirements{
		Limits:   nil,
		Requests: nil,
	}, xenonCase.Resources)
}
func TestGetXenonPorts(t *testing.T) {
	port := []core.ContainerPort{
		{
			Name:          "xenon",
			ContainerPort: 8801,
		},
	}
	assert.Equal(t, port, xenonCase.Ports)
}
func TestGetXenonLivenessProbe(t *testing.T) {
	livenessProbe := &core.Probe{
		Handler: core.Handler{
			Exec: &core.ExecAction{
				Command: []string{"pgrep", "xenon"},
			},
		},
		InitialDelaySeconds: 30,
		TimeoutSeconds:      5,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	assert.Equal(t, livenessProbe, xenonCase.LivenessProbe)
}
func TestGetXenonReadinessProbe(t *testing.T) {
	readinessProbe := &core.Probe{
		Handler: core.Handler{
			Exec: &core.ExecAction{
				Command: []string{"sh", "-c", "xenoncli xenon ping"},
			},
		},
		InitialDelaySeconds: 10,
		TimeoutSeconds:      1,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	assert.Equal(t, readinessProbe, xenonCase.ReadinessProbe)
}
func TestGetXenonVolumeMounts(t *testing.T) {
	volumeMounts := []core.VolumeMount{
		{
			Name:      "scripts",
			MountPath: "/scripts",
		},
		{
			Name:      "xenon",
			MountPath: "/etc/xenon",
		},
	}
	assert.Equal(t, volumeMounts, xenonCase.VolumeMounts)
}
