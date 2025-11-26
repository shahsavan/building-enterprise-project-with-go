package test_containers

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

var (
	mysqlOnce sync.Once
	mysqlHost string
	mysqlPort string
)

func GetMySqlContainer(db, user, pass string, port *int) (string, string, error) {
	mysqlOnce.Do(func() {
		c, err := testContainerRunner{
			servicePort:  3306,
			name:         "mysql",
			image:        "mysql:8.0.36", // pin to a stable version
			exposedPorts: []string{"3306/tcp"},
			env: map[string]string{
				"MYSQL_ROOT_PASSWORD": pass, // root password
				"MYSQL_DATABASE":      db,   // database name
				"MYSQL_USER":          user, // user
				"MYSQL_PASSWORD":      pass, // user password
			},
			hostConfigModifier: mysqlHostConfigModifier(port),
		}.Run()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to run Test Container")
		}
		port, err := c.MappedPort(context.Background(), "5432")
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get port")
		}
		h, err := c.Host(context.Background())
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to get host")
		}
		mysqlHost = h
		mysqlPort = port.Port()
		err = isPortAccessible(mysqlHost, mysqlPort)
	})
	return mysqlHost, mysqlPort, nil
}
func mysqlHostConfigModifier(port *int) func(hostConfig *container.HostConfig) {
	return func(hostConfig *container.HostConfig) {
		hostConfig.AutoRemove = true
		if port != nil {
			hostConfig.PortBindings = nat.PortMap{"5432/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: fmt.Sprintf("%d", *port),
				},
			}}
		}
	}
}
func isPortAccessible(host string, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	timeout := 2 * time.Second
	var err error

	for i := 1; i <= 10; i++ {
		var conn net.Conn
		conn, err = net.DialTimeout("tcp", address, timeout)
		if err == nil {
			log.Info().Msgf("Successfully connected to the MySQL on Host %s and Port %s on attempt %d.\n", host, port, i)
			break
		}
		defer conn.Close()
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("[isPortAccessible] %w", err)
	}
	time.Sleep(10 * time.Second)
	return nil
}
