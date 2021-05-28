package container_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mysqlv1 "github.com/zhyass/mysql-operator/api/v1"
	"github.com/zhyass/mysql-operator/cluster"
	"github.com/zhyass/mysql-operator/cluster/container"

	core "k8s.io/api/core/v1"
)

var _ = Describe("mysql container", func() {
	port := []core.ContainerPort{
		{
			Name:          "mysql",
			ContainerPort: 3306,
		},
	}
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
	reeadinessProbe := &core.Probe{
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
	mysqlCluster := mysqlv1.Cluster{
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
	testCluster := cluster.Cluster{
		Cluster: &mysqlCluster,
	}
	mysqlCase := container.EnsureContainer("mysql", &testCluster)
	Context("initially", func() {
		It("should equal", func() {
			//getName
			Expect("mysql").To(Equal(mysqlCase.Name))
			//getImage
			Expect("percona/percona-server:5.7.33").To(Equal(mysqlCase.Image))
			//getResources
			Expect(core.ResourceRequirements{Limits: nil, Requests: nil}).To(Equal(mysqlCase.Resources))
			//getPorts
			Expect(port).To(Equal(mysqlCase.Ports))
			//getLivenessProbe
			Expect(livenessProbe).To(Equal(mysqlCase.LivenessProbe))
			//getReadinessProbe
			Expect(reeadinessProbe).To(Equal(mysqlCase.ReadinessProbe))
			//getVolumeMounts
			Expect(volumeMounts).To(Equal(mysqlCase.VolumeMounts))
		})
		It("should be Nil", func() {
			//getCommand
			Expect(mysqlCase.Command).To(BeNil())
			//getLifecycle
			Expect(mysqlCase.Lifecycle).To(BeNil())
			//getEnvVars
			Expect(mysqlCase.Env).To(BeNil())
		})
	})
	//getEnvVars
	Context("when initTokuDB", func() {
		initTokuDBCluster := mysqlCluster
		initTokuDBCluster.Spec.MysqlOpts.InitTokuDB = true
		testTokuDBCluster := cluster.Cluster{
			Cluster: &initTokuDBCluster,
		}
		initTokuDBCase := container.EnsureContainer("mysql", &testTokuDBCluster)
		initTokuDBEnv := []core.EnvVar{
			{
				Name:  "INIT_TOKUDB",
				Value: "1",
			},
		}
		It("should equal", func() {
			Expect(initTokuDBEnv).To(Equal(initTokuDBCase.Env))
		})
	})
})
