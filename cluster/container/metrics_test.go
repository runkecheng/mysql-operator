package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	core "k8s.io/api/core/v1"
)

var (
	metricsMysqlCluster = mysqlv1.Cluster{
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
			MetricsOpts: mysqlv1.MetricsOpts{
				Image: "metrics-image",
			},
		},
	}
	testMetricsCluster = cluster.Cluster{
		Cluster: &metricsMysqlCluster,
	}
	metricsCase = EnsureContainer("metrics", &testMetricsCluster)
)

func TestGetMetricsName(t *testing.T) {
	assert.Equal(t, "metrics", metricsCase.Name)
}
func TestGetMetricsImage(t *testing.T) {
	assert.Equal(t, "metrics-image", metricsCase.Image)
}
func TestGetMetricsCommand(t *testing.T) {
	assert.Nil(t, metricsCase.Command)
}
func TestGetMetricsEnvVar(t *testing.T) {
	{
		optTrue := true
		env := []core.EnvVar{
			{
				Name: "DATA_SOURCE_NAME",
				ValueFrom: &core.EnvVarSource{
					SecretKeyRef: &core.SecretKeySelector{
						LocalObjectReference: core.LocalObjectReference{
							Name: "sample-secret",
						},
						Key:      "data-source",
						Optional: &optTrue,
					},
				},
			},
		}
		assert.Equal(t, env, metricsCase.Env)
	}
}
func TestGetMetricsLifecycle(t *testing.T) {
	assert.Nil(t, metricsCase.Lifecycle)
}
func TestGetMetricsResources(t *testing.T) {
	assert.Equal(t, core.ResourceRequirements{
		Limits:   nil,
		Requests: nil,
	}, metricsCase.Resources)
}
func TestGetMetricsPorts(t *testing.T) {
	port := []core.ContainerPort{
		{
			Name:          "metrics",
			ContainerPort: 9104,
		},
	}
	assert.Equal(t, port, metricsCase.Ports)
}
func TestGetMetricsLivenessProbe(t *testing.T) {
	livenessProbe := &core.Probe{
		Handler: core.Handler{
			HTTPGet: &core.HTTPGetAction{
				Path: "/",
				Port: intstr.IntOrString{
					Type:   0,
					IntVal: int32(9104),
				},
			},
		},
		InitialDelaySeconds: 15,
		TimeoutSeconds:      5,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	assert.Equal(t, livenessProbe, metricsCase.LivenessProbe)
}
func TestGetMetricsReadinessProbe(t *testing.T) {
	readinessProbe := &core.Probe{
		Handler: core.Handler{
			HTTPGet: &core.HTTPGetAction{
				Path: "/",
				Port: intstr.IntOrString{
					Type:   0,
					IntVal: int32(9104),
				},
			},
		},
		InitialDelaySeconds: 5,
		TimeoutSeconds:      1,
		PeriodSeconds:       10,
		SuccessThreshold:    1,
		FailureThreshold:    3,
	}
	assert.Equal(t, readinessProbe, metricsCase.ReadinessProbe)
}
func TestGetMetricsVolumeMounts(t *testing.T) {

	assert.Nil(t, metricsCase.VolumeMounts)
}
