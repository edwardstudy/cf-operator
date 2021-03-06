package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Deploy", func() {
	Context("when correctly setup", func() {
		podName := "diego-pod"

		AfterEach(func() {
			env.WaitForPodsDelete(env.Namespace)
		})

		It("should deploy a pod", func() {
			tearDown, err := env.CreateConfigMap(env.Namespace, env.DefaultBOSHManifest("manifest"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			_, tearDown, err = env.CreateBOSHDeployment(env.Namespace, env.DefaultBOSHDeployment("test", "manifest"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			// check for pod
			err = env.WaitForPod(env.Namespace, podName)
			Expect(err).NotTo(HaveOccurred(), "error waiting for pod from deployment")
		})

		It("should deploy manifest with multiple ops correctly", func() {
			tearDown, err := env.CreateConfigMap(env.Namespace, env.DefaultBOSHManifest("manifest"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			tearDown, err = env.CreateConfigMap(env.Namespace, env.InterpolateOpsConfigMap("bosh-ops"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			tearDown, err = env.CreateSecret(env.Namespace, env.InterpolateOpsSecret("bosh-ops-secret"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			_, tearDown, err = env.CreateBOSHDeployment(env.Namespace, env.InterpolateBOSHDeployment("test", "manifest", "bosh-ops", "bosh-ops-secret"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			// check for pod
			err = env.WaitForPod(env.Namespace, podName)
			Expect(err).NotTo(HaveOccurred(), "error waiting for pod from deployment")
			labeled, err := env.PodLabeled(env.Namespace, podName, "size", "4")
			Expect(labeled).To(BeTrue())
			Expect(err).NotTo(HaveOccurred(), "error verifying pod label")
		})
	})

	Context("when incorrectly setup", func() {
		It("failed to deploy if ResolveCRD error occurred", func() {
			tearDown, err := env.CreateConfigMap(env.Namespace, env.DefaultBOSHManifest("manifest"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			tearDown, err = env.CreateConfigMap(env.Namespace, env.InterpolateOpsConfigMap("bosh-ops"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			tearDown, err = env.CreateSecret(env.Namespace, env.InterpolateOpsIncorrectSecret("bosh-ops-secret"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			boshDeployment, tearDown, err := env.CreateBOSHDeployment(env.Namespace, env.InterpolateBOSHDeployment("test", "manifest", "bosh-ops", "bosh-ops-secret"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			// check for events
			events, err := env.GetBOSHDeploymentEvents(env.Namespace, boshDeployment.ObjectMeta.Name, string(boshDeployment.ObjectMeta.UID))
			Expect(err).NotTo(HaveOccurred())
			Expect(env.ContainExpectedEvent(events, "ResolveCRD Error", "Failed to interpolate")).To(BeTrue())
		})

		It("failed to deploy a empty manifest", func() {
			tearDown, err := env.CreateConfigMap(env.Namespace, env.DefaultBOSHManifest("manifest"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			_, tearDown, err = env.CreateBOSHDeployment(env.Namespace, env.EmptyBOSHDeployment("test", "manifest"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.manifest.type in body should be one of"))
			Expect(err.Error()).To(ContainSubstring("spec.manifest.ref in body should be at least 1 chars long"))
			defer tearDown()
		})

		It("failed to deploy due to a wrong manifest type", func() {
			tearDown, err := env.CreateConfigMap(env.Namespace, env.DefaultBOSHManifest("manifest"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			_, tearDown, err = env.CreateBOSHDeployment(env.Namespace, env.WrongTypeBOSHDeployment("test", "manifest"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.manifest.type in body should be one of"))
			defer tearDown()
		})

		It("failed to deploy due to a empty manifest ref", func() {
			tearDown, err := env.CreateConfigMap(env.Namespace, env.DefaultBOSHManifest("manifest"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			_, tearDown, err = env.CreateBOSHDeployment(env.Namespace, env.DefaultBOSHDeployment("test", ""))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.manifest.ref in body should be at least 1 chars long"))
			defer tearDown()
		})

		It("failed to deploy due to a wrong ops type", func() {
			tearDown, err := env.CreateConfigMap(env.Namespace, env.DefaultBOSHManifest("manifest"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			_, tearDown, err = env.CreateBOSHDeployment(env.Namespace, env.BOSHDeploymentWithWrongTypeOps("test", "manifest", "ops"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.ops.type in body should be one of"))
			defer tearDown()
		})

		It("failed to deploy due to a empty ops ref", func() {
			tearDown, err := env.CreateConfigMap(env.Namespace, env.DefaultBOSHManifest("manifest"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			_, tearDown, err = env.CreateBOSHDeployment(env.Namespace, env.DefaultBOSHDeploymentWithOps("test", "manifest", ""))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("spec.ops.ref in body should be at least 1 chars long"))
			defer tearDown()
		})
	})
})
