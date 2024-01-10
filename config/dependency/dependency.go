package dependency

import (
	"mydream_project/config"
	"mydream_project/pkg"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type Depend struct {
	dig.In
	Db      	*gorm.DB 
	Config      *config.Config
	Echo    	*echo.Echo 
	Log     	*logrus.Logger 
	Gcp         *pkg.StorageGCP
    Rds     	*redis.Client 
	Mds     	*pkg.Midtrans
	Nsq     	*pkg.NSQProducer
	Validation  *pkg.Validation
	PromErr map[string]string 
}