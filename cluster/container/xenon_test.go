package container_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	"github.com/zhyass/mysql-operator/cluster/container"

	core "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("xenon container", func() {
	var replicas int32 = 1
	lifecycle := &core.Lifecycle{
		PostStart: &core.Handler{
			Exec: &core.ExecAction{
				Command: []string{"sh", "-c",
					"until (xenoncli xenon ping && xenoncli cluster add sample-mysql-0.sample-mysql.default:8801) > /dev/null 2>&1; do sleep 2; done",
				},
			},
		},
	}
	port := []core.ContainerPort{
		{
			Name:          "xenon",
			ContainerPort: 8801,
		},
	}
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
	mysqlCluster := mysqlv1.Cluster{
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
			Replicas: &replicas,
		},
	}
	testCluster := cluster.Cluster{
		Cluster: &mysqlCluster,
	}
	xenonCase := container.EnsureContainer("xenon", &testCluster)

	Context("initially", func() {
		It("should equal", func() {
			//getName
			Expect("xenon").To(Equal(xenonCase.Name))
			//getImage
			Expect("xenon image").To(Equal(xenonCase.Image))
			//getLifecycle
			Expect(lifecycle).To(Equal(xenonCase.Lifecycle))
			//getResources
			Expect(core.ResourceRequirements{Limits: nil, Requests: nil}).To(Equal(xenonCase.Resources))
			//getPorts
			Expect(port).To(Equal(xenonCase.Ports))
			//getLivenessProbe
			Expect(livenessProbe).To(Equal(xenonCase.LivenessProbe))
			//getReadinessProbe
			Expect(readinessProbe).To(Equal(xenonCase.ReadinessProbe))
			//getVolumeMounts
			Expect(volumeMounts).To(Equal(xenonCase.VolumeMounts))
		})
		It("should be Nil", func() {
			//getCommand
			Expect(xenonCase.Command).To(BeNil())
			//getEnv
			Expect(xenonCase.Env).To(BeNil())
			//getEnvVars
			Expect(xenonCase.Env).To(BeNil())
		})
	})
	//getLifecycle
	Context("when replicas greater than one", func() {
		testLifecycle := &core.Lifecycle{
			PostStart: &core.Handler{
				Exec: &core.ExecAction{
					Command: []string{"sh", "-c",
						"until (xenoncli xenon ping && xenoncli cluster add sample-mysql-0.sample-mysql.default:8801,sample-mysql-1.sample-mysql.default:8801) > /dev/null 2>&1; do sleep 2; done",
					},
				},
			},
		}
		It("should equal", func() {
			replicas = 2
			replicasCase := container.EnsureContainer("xenon", &testCluster)
			Expect(testLifecycle).To(Equal(replicasCase.Lifecycle))
		})
	})

})
