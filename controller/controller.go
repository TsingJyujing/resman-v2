package controller

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	"net/http"
	"os"
	"resman/controller/dao"
)

//go:embed sqlc/schema.sql
var ddl string

var logger = logrus.New()

func GetDDL() string {
	return ddl
}

type Controller struct {
	// Context used by all APIs
	db dao.Queries
}

func New(db dao.Queries) *Controller {
	return &Controller{
		db: db,
	}
}

func HashString(inputString string) string {
	hasher := sha256.New()
	hasher.Write([]byte(inputString))
	hashBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashBytes)

}

func (c *Controller) InitAdminUser(ctx context.Context) error {
	userCount, err := c.db.ListUsers(ctx) // FIXME better to have a count query
	if err != nil {
		return err
	}
	if len(userCount) == 0 {
		logger.Info("No users found, creating default admin user")
		defaultPassword := os.Getenv("DEFAULT_ADMIN_PASSWORD")
		if defaultPassword == "" {
			logger.Warn("No default admin password, using 'resman' as password")
			defaultPassword = "resman"
		}
		_, err := c.db.CreateUser(ctx, dao.CreateUserParams{
			ID:             "admin",
			HashedPassword: HashString(defaultPassword),
		})
		if err != nil {
			return err
		}
	} else {
		logger.Infof("Found %d users, skipping admin user creation", len(userCount))
	}
	return nil
}

// UserLogin handles user login and returns a token upon successful authentication.
func (c *Controller) UserLogin(echoCtx echo.Context) error {
	ctx := echoCtx.Request().Context()
	userId, password, ok := echoCtx.Request().BasicAuth()
	if !ok {
		return echoCtx.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "missing or invalid Authorization header",
		})
	}
	hashedPassword := HashString(password)
	user, err := c.db.GetUser(ctx, userId)
	if err != nil {
		return err
	}
	if user.HashedPassword != hashedPassword {
		return echoCtx.JSON(http.StatusUnauthorized, map[string]interface{}{
			"error": "invalid credentials",
		})
	}
	return echoCtx.JSON(200, map[string]string{"token": "TODO add JWT token"})
}
