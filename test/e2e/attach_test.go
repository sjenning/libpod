package integration

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Podman attach", func() {
	var (
		tempdir    string
		err        error
		podmanTest PodmanTest
	)

	BeforeEach(func() {
		tempdir, err = CreateTempDirInTempDir()
		if err != nil {
			os.Exit(1)
		}
		podmanTest = PodmanCreate(tempdir)
		podmanTest.RestoreAllArtifacts()
	})

	AfterEach(func() {
		podmanTest.Cleanup()

	})

	It("podman attach to bogus container", func() {
		session := podmanTest.Podman([]string{"attach", "foobar"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(125))
	})

	It("podman attach to non-running container", func() {
		session := podmanTest.Podman([]string{"create", "--name", "test1", "-d", "-i", ALPINE, "ls"})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		results := podmanTest.Podman([]string{"attach", "test1"})
		results.WaitWithDefaultTimeout()
		Expect(results.ExitCode()).To(Equal(125))
	})

	It("podman attach to multiple containers", func() {
		session := podmanTest.RunSleepContainer("test1")
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		session = podmanTest.RunSleepContainer("test2")
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(Equal(0))

		results := podmanTest.Podman([]string{"attach", "test1", "test2"})
		results.WaitWithDefaultTimeout()
		Expect(results.ExitCode()).To(Equal(125))
	})
})
