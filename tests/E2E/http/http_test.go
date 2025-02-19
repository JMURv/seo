package http

import (
	"database/sql"
	"fmt"
	"github.com/JMURv/seo/internal/cache/redis"
	"github.com/JMURv/seo/internal/config"
	"github.com/JMURv/seo/internal/ctrl"
	hdl "github.com/JMURv/seo/internal/hdl/http"
	"github.com/JMURv/seo/internal/repo/db"
	"github.com/JMURv/seo/tests/mocks"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const configPath = "../../../configs/test.config.yaml"
const getTables = `
SELECT tablename 
FROM pg_tables 
WHERE schemaname = 'public';
`

func setupTestServer(t *testing.T) (*httptest.Server, *mocks.MockSSOSvc, func()) {
	ctrlMock := gomock.NewController(t)
	defer ctrlMock.Finish()
	ssoCtrl := mocks.NewMockSSOSvc(ctrlMock)

	zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	conf := config.MustLoad(configPath)

	repo := db.New(conf.DB)
	cache := redis.New(conf.Redis)
	svc := ctrl.New(repo, cache)
	//h := hdl.New(svc, sso.New(conf.Services))
	h := hdl.New(svc, ssoCtrl)

	mux := http.NewServeMux()
	hdl.RegisterSEORoutes(mux, h)
	hdl.RegisterPageRoutes(mux, h)

	cleanupFunc := func() {
		conn, err := sql.Open(
			"postgres", fmt.Sprintf(
				"postgres://%s:%s@%s:%d/%s?sslmode=disable",
				conf.DB.User,
				conf.DB.Password,
				conf.DB.Host,
				conf.DB.Port,
				conf.DB.Database,
			),
		)
		if err != nil {
			zap.L().Fatal("Failed to connect to the database", zap.Error(err))
		}

		if err = conn.Ping(); err != nil {
			zap.L().Fatal("Failed to ping the database", zap.Error(err))
		}

		rows, err := conn.Query(getTables)
		if err != nil {
			zap.L().Fatal("Failed to fetch table names", zap.Error(err))
		}
		defer func(rows *sql.Rows) {
			if err := rows.Close(); err != nil {
				zap.L().Debug("Error while closing rows", zap.Error(err))
			}
		}(rows)

		var tables []string
		for rows.Next() {
			var name string
			if err := rows.Scan(&name); err != nil {
				zap.L().Fatal("Failed to scan table name", zap.Error(err))
			}
			tables = append(tables, name)
		}

		if len(tables) == 0 {
			return
		}

		_, err = conn.Exec(fmt.Sprintf("TRUNCATE TABLE %v RESTART IDENTITY CASCADE;", strings.Join(tables, ", ")))
		if err != nil {
			zap.L().Fatal("Failed to truncate tables", zap.Error(err))
		}
	}

	return httptest.NewServer(mux), ssoCtrl, cleanupFunc
}
