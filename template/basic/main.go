package main

import (
	"errors"
	"flag"
	"os"

	"github.com/opensourceways/server-common-lib/config"
	"github.com/opensourceways/server-common-lib/logrusutil"
	liboptions "github.com/opensourceways/server-common-lib/options"
	"github.com/opensourceways/server-common-lib/secret"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/robot-gitlab-lib/client"
	"github.com/opensourceways/robot-gitlab-lib/framework"
)

type options struct {
	service  liboptions.ServiceOptions
	token    liboptions.GitLabOptions
	endpoint string
}

func (o *options) Validate() error {
	if err := o.service.Validate(); err != nil {
		return err
	}

	return o.token.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.token.AddFlags(fs)
	o.service.AddFlags(fs)
	fs.StringVar(&o.endpoint, "gitlab-endpoint", "", "the endpoint of gitlab.")

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit(botName)

	o := gatherOptions(flag.NewFlagSet(os.Args[0], flag.ExitOnError), os.Args[1:]...)
	if err := o.Validate(); err != nil {
		logrus.WithError(err).Fatal("Invalid options")
	}

	secretAgent := new(secret.Agent)
	if err := secretAgent.Start([]string{o.token.TokenPath}); err != nil {
		logrus.WithError(err).Fatal("Error starting secret agent.")
	}

	defer secretAgent.Stop()

	agent := config.NewConfigAgent(func() config.Config {
		return &configuration{}
	})

	if err := agent.Start(o.service.ConfigFile); err != nil {
		logrus.WithError(err).Errorf("start config:%s", o.service.ConfigFile)
		return
	}

	defer agent.Stop()

	c := client.NewGitlabClient(
		secretAgent.GetTokenGenerator(o.token.TokenPath),
		o.endpoint,
	)

	r := newRobot(c, func() (*configuration, error) {
		_, cfg := agent.GetConfig()
		if c, ok := cfg.(*configuration); ok {
			return c, nil
		}

		return nil, errors.New("can't convert to configuration")
	})

	framework.Run(r, o.service.Port, o.service.GracePeriod)
}
