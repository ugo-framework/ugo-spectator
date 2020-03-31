package runner

import "os/exec"

type Runner interface {
	Run() (*exec.Cmd, error)
	Kill() error
}
