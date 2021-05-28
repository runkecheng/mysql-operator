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

	Describe("auditlog", func() {
		mysqlCluster := mysqlv1.Cluster{
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
		testCluster := cluster.Cluster{
			Cluster: &mysqlCluster,
		}
		auditLogCase := container.EnsureContainer("auditlog", &testCluster)
		command := []string{"tail", "-f", "/var/log/mysql" + "/mysql-audit.log"}
		volumeMounts := []core.VolumeMount{
			{
				Name:      "logs",
				MountPath: "/var/log/mysql",
			},
		}

		It("should equal", func() {
			Expect("auditlog").To(Equal(auditLogCase.Name))
			//getImage
			Expect("busybox").To(Equal(auditLogCase.Image))
			//getCommand
			Expect(command).To(Equal(auditLogCase.Command))
			//getResources
			Expect(core.ResourceRequirements{Limits: nil, Requests: nil}).To(Equal(auditLogCase.Resources))
			//getVolumeMounts
			Expect(volumeMounts).To(Equal(auditLogCase.VolumeMounts))

		})
		It("should be nil", func() {
			//getEnvVars
			Expect(auditLogCase.Env).To(BeNil())
			//getLifecycle
			Expect(auditLogCase.Lifecycle).To(BeNil())
			//getPorts
			Expect(auditLogCase.Ports).To(BeNil())
			//getLivenessProbe
			Expect(auditLogCase.LivenessProbe).To(BeNil())
			//getReadinessProbe
			Expect(auditLogCase.ReadinessProbe).To(BeNil())
		})

	})

})
