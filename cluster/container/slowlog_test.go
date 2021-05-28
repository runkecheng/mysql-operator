package container_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	"github.com/zhyass/mysql-operator/cluster/container"

	core "k8s.io/api/core/v1"
)

var _ = Describe("Container", func() {

	Describe("slowlog", func() {
		command := []string{"tail", "-f", "/var/log/mysql" + "/mysql-slow.log"}
		volumeMounts := []core.VolumeMount{
			{
				Name:      "logs",
				MountPath: "/var/log/mysql",
			},
		}
		mysqlCluster := mysqlv1.Cluster{
			Spec: mysqlv1.ClusterSpec{
				PodSpec: mysqlv1.PodSpec{
					SidecarImage: "sidecar-image",
					Resources: core.ResourceRequirements{
						Limits:   nil,
						Requests: nil,
					},
				},
			},
		}
		testCluster := cluster.Cluster{
			Cluster: &mysqlCluster,
		}
		slowlogCase := container.EnsureContainer("slowlog", &testCluster)

		It("should equal", func() {
			//getName
			Expect("slowlog").To(Equal(slowlogCase.Name))
			//getImage
			Expect("sidecar-image").To(Equal(slowlogCase.Image))
			//getCommand
			Expect(command).To(Equal(slowlogCase.Command))
			//getResources
			Expect(core.ResourceRequirements{Limits: nil, Requests: nil}).To(Equal(slowlogCase.Resources))
			//getVolumeMounts
			Expect(volumeMounts).To(Equal(slowlogCase.VolumeMounts))

		})
		It("should be nil", func() {
			//getEnvVars
			Expect(slowlogCase.Env).To(BeNil())
			//getLifecycle
			Expect(slowlogCase.Lifecycle).To(BeNil())
			//getPorts
			Expect(slowlogCase.Ports).To(BeNil())
			//getLivenessProbe
			Expect(slowlogCase.LivenessProbe).To(BeNil())
			//getReadinessProbe
			Expect(slowlogCase.ReadinessProbe).To(BeNil())
		})
	})
})
