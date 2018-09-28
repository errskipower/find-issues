package acceptance_test

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Open Issues", func() {
	It("returns the list of open issues for a given repository", func() {
		repo := "ghc-tdd/find-issues"

		output := execute(repo)

		Expect(output).To(ContainSubstring("#1: User can get list of open issues on a given repo"))
	})

	Context("when filtering by help wanted label", func() {
		It("returns the list of open issues that are tagged with help wanted", func() {
			repo := "ghc-tdd/find-issues"

			output := execute(repo, `--label`, `"help wanted"`)

			Expect(output).NotTo(ContainSubstring("#1: User can get list of open issues on a given repo"))
			Expect(output).To(ContainSubstring(`#2: User can see only those issues that have the "help wanted" label`))
		})
	})

	Context("when filtering by github creator", func() {
		It("returns the list of open issues that were created by a provided github username", func() {
			repo := "ghc-tdd/find-issues"

			output := execute(repo, "--creator", "angelachin")

			// TODO: Add expectations here.
			Expect(output).To(ContainSubstring("#9: Mock Issue for Username and Label Filter Testing"))
			Expect(output).To(ContainSubstring("#8: Mock Issue for Username Filter Testing"))
			Expect(output).NotTo(ContainSubstring("#1: User can get list of all open issues on a given repo"))
		})
	})
})

func execute(args ...string) string {
	stdout := bytes.NewBuffer([]byte{})
	stderr := bytes.NewBuffer([]byte{})

	cmd := exec.Command(binaryPath, args...)

	session, err := gexec.Start(cmd, stdout, stderr)
	Expect(err).NotTo(HaveOccurred())
	Eventually(session, 1*time.Minute).Should(gexec.Exit(0))

	return strings.TrimSpace(string(stdout.Bytes()))
}
