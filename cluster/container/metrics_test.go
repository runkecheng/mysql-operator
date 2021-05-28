package container_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/zhyass/mysql-operator/utils"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	"github.com/zhyass/mysql-operator/cluster/container"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	core "k8s.io/api/core/v1"
)

var _ = Describe("Container", func() {

	Describe("metrics", func() {
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
		livenessProbe := core.Probe{
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
		readinessProbe := core.Probe{
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
				MetricsOpts: mysqlv1.MetricsOpts{
					Image: "metrics-image",
				},
			},
		}
		testCluster := cluster.Cluster{
			Cluster: &mysqlCluster,
		}
		testMetrics := container.EnsureContainer("metrics", &testCluster)

		It("should equal", func() {
			//getName
			Expect("metrics").To(Equal(testMetrics.Name))
			//getImage
			Expect("metrics-image").To(Equal(testMetrics.Image))
			//getEnvVars
			Expect(env).To(Equal(testMetrics.Env))
			//getResources
			Expect(core.ResourceRequirements{Limits: nil, Requests: nil}).To(Equal(testMetrics.Resources))
			//getPorts
			Expect([]core.ContainerPort{
				{
					Name:          utils.MetricsPortName,
					ContainerPort: utils.MetricsPort,
				},
			}).To(Equal(testMetrics.Ports))
			//getLivenessProbe
			Expect(&livenessProbe).To(Equal(testMetrics.LivenessProbe))
			//getReadinessProbe
			Expect(&readinessProbe).To(Equal(testMetrics.ReadinessProbe))
		})
		It("should be nil", func() {
			//getCommand
			Expect(testMetrics.Command).To(BeNil())
			//getLifecycle
			Expect(testMetrics.Lifecycle).To(BeNil())
			//getVolumeMounts
			Expect(testMetrics.VolumeMounts).To(BeNil())
		})
	})
})
