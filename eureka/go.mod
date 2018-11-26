module github.com/HikoQiu/go-eureka-client/eureka

replace (
	golang.org/x/crypto v0.0.0-20180820150726-614d502a4dac => github.com/golang/crypto v0.0.0-20180820150726-614d502a4dac
	golang.org/x/net v0.0.0-20181011144130-49bb7cea24b1 => github.com/golang/net v0.0.0-20180826012351-8a410e7b638d
	golang.org/x/text v0.3.0 => github.com/golang/text v0.3.0
)

require (
	github.com/miekg/dns v1.0.15
	gopkg.in/resty.v1 v1.10.2
)
