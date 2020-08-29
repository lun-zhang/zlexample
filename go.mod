module zlexample

go 1.12

require (
	github.com/aws/aws-xray-sdk-go v1.0.0-rc.5.0.20180720202646-037b81b2bf76
	github.com/fvbock/endless v0.0.0-20170109170031-447134032cb6
	github.com/gin-gonic/gin v1.6.2
	github.com/go-playground/universal-translator v0.17.0 // indirect
	github.com/leodido/go-urn v1.2.0 // indirect
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/lun-zhang/gorm v1.9.14-beta.1.14.0
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
	gopkg.in/go-playground/validator.v8 v8.18.2 // indirect
	zlutils v0.0.0
)

replace zlutils v0.0.0 => github.com/lun-zhang/zlutils/v7 v7.26.0
