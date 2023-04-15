package internal_test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	tc "github.com/testcontainers/testcontainers-go/modules/compose"
	"github.com/testcontainers/testcontainers-go/wait"
	"test-containers/internal"
	"testing"
	"time"
)

type DBSuite struct {
	suite.Suite
	comp tc.ComposeStack
	db   *internal.PostgresPool
}

func (s *DBSuite) SetupSuite() {
	const (
		host = "localhost"
		port = "5432"
	)

	var dsnBuilder = func(host string, port nat.Port) string {
		return "postgres://" + host + ":" + port.Port() + "/postgres?user=ost&password=pass"
	}

	comp, err := s.upPostgresContainer(port, dsnBuilder)
	require.NoError(s.T(), err, "Container up failed")
	s.comp = comp

	db, err := internal.NewPostgres(dsnBuilder(host, port))
	if err != nil {
		s.T().Fatal("Failed to make connection pool", err)
	}
	s.db = db
}

func (s *DBSuite) TearDownSuite() {
	s.db.Close()
	assert.NoError(s.T(), s.comp.Down(context.TODO(), tc.RemoveOrphans(true), tc.RemoveImagesLocal), "compose.Down()")
}

func TestDBSuite_Run(t *testing.T) {
	suite.Run(t, &DBSuite{})
}

func (s *DBSuite) TestGreeting() {
	greeting, err := s.db.ReadHelloWorld()
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "Hello, world!", greeting)
}

func (s *DBSuite) upPostgresContainer(port nat.Port, dsnBuilder func(host string, port nat.Port) string) (tc.ComposeStack, error) {
	identifier := tc.StackIdentifier("db_tests")

	comp, err := tc.NewDockerComposeWith(tc.WithStackFiles("testdata/docker-compose.yaml"), identifier)
	if err != nil {
		return nil, fmt.Errorf("docker compose init:%v", err)
	}

	waitForSql := wait.
		ForSQL(port, "pgx", dsnBuilder).
		WithQuery("SELECT 1").
		WithPollInterval(1 * time.Second)
	comp.WaitForService("db", waitForSql)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := comp.Up(ctx, tc.Wait(true)); err != nil {
		return nil, fmt.Errorf("docker compose up:%v", err)
	}

	return comp, nil
}
