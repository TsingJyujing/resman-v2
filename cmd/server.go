package cmd

import (
	"database/sql"
	"net/http"
	"resman/controller"
	"resman/controller/dao"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	_ "modernc.org/sqlite"
)

var logger = logrus.New()

var serverCommand = &cobra.Command{
	Use:   "server",
	Short: "Starting server",
	Run: func(cmd *cobra.Command, args []string) {
		e := echo.New()
		goCtx := cmd.Context()
		db, err := sql.Open("sqlite", ":memory:")
		if err != nil {
			logger.WithError(err).Fatal("Failed to open database")
		}
		// create tables
		if _, err := db.ExecContext(goCtx, controller.GetDDL()); err != nil {
			logger.WithError(err).Fatal("Failed to create tables")
		}
		c := controller.New(*dao.New(db))
		initErr := c.InitAdminUser(cmd.Context())
		if initErr != nil {
			logger.WithError(initErr).Fatal("Failed to initialize admin user")
		}
		e.GET("/health", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
		})
		e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
		// User Login and authentication
		e.GET("/login", c.UserLogin)
		e.Logger.Fatal(e.Start(":1323"))
	},
}
