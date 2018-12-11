package integration_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ExtendedJob", func() {
	Context("when correctly setup", func() {
		AfterEach(func() {
			err := env.WaitForPodsDelete(env.Namespace)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should start a job", func() {
			_, tearDown, err := env.CreateExtendedJob(env.Namespace, *env.DefaultExtendedJob("extendedjob"))
			Expect(err).NotTo(HaveOccurred())
			defer tearDown()

			// check for job
			//err = env.WaitForJob(env.Namespace, "defaultJob")
			//Expect(err).NotTo(HaveOccurred(), "error waiting for job from extendedjob")
		})
	})
})
