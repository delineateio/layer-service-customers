package discovery

import (
	"fmt"
	"os"

	"delineate.io/customers/src/config"
	consul "github.com/hashicorp/consul/api"
)

func Initialize() {
	register()
}

func register() {
	// Get a new client
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}

	check := new(consul.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("%s://%s:%v/healthz",
		getScheme(),
		getHost(),
		getPort())
	check.Interval = getInterval()
	check.Timeout = getTimeout()
	check.DeregisterCriticalServiceAfter = getDeregister()

	info := &consul.AgentServiceRegistration{
		ID:    getID(),
		Name:  getName(),
		Port:  getPort(),
		Check: check,
	}

	err = client.Agent().ServiceRegister(info)
	if err != nil {
		panic(err)
	}
}

func getScheme() string {
	return config.GetString("server.scheme")
}

func getHost() string {
	hostname, _ := os.Hostname()
	return hostname
}

func getID() string {
	prefix := config.GetString("server.prefix")
	return fmt.Sprintf("%s-%s", prefix, getHost())
}

func getName() string {
	return config.GetString("server.name")
}

func getPort() int {
	return config.GetInt("server.port")
}

func getInterval() string {
	return config.GetString("discovery.interval")
}

func getTimeout() string {
	return config.GetString("discovery.timeout")
}

func getDeregister() string {
	return config.GetString("discovery.deregister")
}
