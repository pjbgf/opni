package targets

import (
	"fmt"
	"os/exec"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Dagger mg.Namespace

func daggerRun() []string {
	dagger, err := exec.LookPath("dagger")
	if err != nil {
		return []string{"go", "run"}
	}
	return []string{dagger, "run", "go", "run"}
}

type daggerPackage string

const (
	dagger  daggerPackage = "./dagger"
	daggerx daggerPackage = "./dagger/x"
)

// Invokes 'go run ./dagger' with all arguments
func (ns Dagger) Run(arg0 string) error {
	return ns.run(dagger, takeArgv(arg0)...)
}

func (Dagger) run(pkg daggerPackage, args ...string) error {
	cmds := daggerRun()
	return sh.RunV(cmds[0], append(append(cmds[1:], string(pkg)), args...)...)
}

func (Dagger) do(outputDir string, args ...string) error {
	daggerBinary, err := exec.LookPath("dagger")
	if err != nil {
		return fmt.Errorf("could not find dagger: %w", err)
	}
	return sh.Run(daggerBinary, append([]string{"do", "--output", outputDir, "--project", string(daggerx), "--workdir", string(dagger)}, args...)...)
}

// Invokes 'go run ./dagger --help'
func (ns Dagger) Help() error {
	return sh.RunV(mg.GoCmd(), "run", string(dagger), "--help")
}

// Invokes 'go run ./dagger --setup'
func (Dagger) Setup() error {
	return sh.RunV(mg.GoCmd(), "run", string(dagger), "--setup")
}
