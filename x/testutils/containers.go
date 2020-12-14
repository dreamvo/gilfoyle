package testutils

import (
	"fmt"
	"github.com/orlangure/gnomock"
	"github.com/orlangure/gnomock/preset/rabbitmq"
	"testing"
)

func StopContainer(t *testing.T, c *gnomock.Container) {
	err := gnomock.Stop(c)
	if err != nil {
		// We don't want this error to be fatal, as the container will be stopped anyway
		t.Log(fmt.Errorf("failed to stop container: %e", err))
	}
}

func CreateRabbitMQContainer(t *testing.T, user, password string) *gnomock.Container {
	mq := rabbitmq.Preset(
		rabbitmq.WithUser(user, password),
		rabbitmq.WithVersion("3.8-alpine"),
	)
	container, err := gnomock.Start(mq)
	if err != nil {
		t.Error(err)
	}

	return container
}
