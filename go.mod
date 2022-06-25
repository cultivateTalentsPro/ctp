module ctp

go 1.16

replace google.golang.org/grpc => google.golang.org/grpc v1.27.0

require (
	github.com/gin-gonic/gin v1.8.1
	github.com/go-micro/plugins/v4/registry/etcd v1.0.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.2.0 // indirect
	go-micro.dev/v4 v4.6.0
	gopkg.in/ini.v1 v1.62.0
)
