module ctp

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.27.0

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gin-gonic/gin v1.8.1
	github.com/go-micro/plugins/v4/registry/etcd v1.0.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jinzhu/gorm v1.9.16
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/lib/pq v1.2.0 // indirect
	github.com/sirupsen/logrus v1.7.0
	go-micro.dev/v4 v4.6.0
	golang.org/x/crypto v0.0.0-20210711020723-a769d52b0f97
	gopkg.in/ini.v1 v1.62.0
)
