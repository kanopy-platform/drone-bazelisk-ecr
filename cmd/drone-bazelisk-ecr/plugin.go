package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/kelseyhightower/envconfig"
)

// plugin configuraion
type plugin struct {
	Target     string `required:"true"`
	Registry   string `required:"true"`
	Repository string
	Tag        string
	AccessKey  string `split_words:"true"`
	SecretKey  string `split_words:"true"`
}

// plugin constructor
func newPlugin() plugin {
	return plugin{}
}

// process plugin env vars
func (p *plugin) setenv() error {
	err := envconfig.Process("plugin", p)
	if err != nil {
		return err
	}

	// convenience variables to be read by bazel workspace status scripts
	if p.Registry != "" {
		setEnvWithPrefix("REGISTRY", p.Registry)
	}
	if p.Repository != "" {
		setEnvWithPrefix("REPOSITORY", p.Repository)
	}
	if p.Tag != "" {
		setEnvWithPrefix("TAG", p.Tag)
	}

	// setup the credentials used by the amazon-ecr-credential-helper
	if p.AccessKey != "" && p.SecretKey != "" {
		os.Setenv("AWS_ACCESS_KEY_ID", p.AccessKey)
		os.Setenv("AWS_SECRET_ACCESS_KEY", p.SecretKey)
	}

	return nil
}

func (p *plugin) getArgs() []string {
	var args []string

	// append run and target
	args = append(args, "run", p.Target)

	return args
}

// runs the bazel command
func (p *plugin) run() error {
	err := p.setenv()
	if err != nil {
		return err
	}

	// exec bazel
	cmd := exec.Command("bazel", p.getArgs()...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func setEnvWithPrefix(key, val string) {
	os.Setenv(fmt.Sprintf("%s_%s", "DRONE_ECR", key), val)
}
