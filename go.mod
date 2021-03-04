module github.com/wostzone/logger

go 1.14

require (
	github.com/sirupsen/logrus v1.8.0
	github.com/stretchr/testify v1.3.0
	github.com/wostzone/gateway v0.0.0-20210227062304-ae0ec41fc4b7
)

// Until gateway is stable
replace github.com/wostzone/gateway => ../gateway
