package main

import (
    "github.com/HikoQiu/go-eureka-client/eureka"
)

func main() {
    config := eureka.GetDefaultEurekaClientConfig()
    config.UseDnsForFetchingServiceUrls = false
    config.Region = "region-cn-hd-1"
    config.AvailabilityZones = map[string]string{
        "region-cn-hd-1": "zone-cn-hz-1",
    }
    config.ServiceUrl = map[string]string{
        "zone-cn-hz-1": "http://192.168.20.236:9001/eureka,http://192.168.20.237:9001/eureka",
    }

    // custom logger
    //eureka.SetLogger(func(format string, a ...interface{}) {
    //   fmt.Println("[custom logger] " + format, a)
    //})

    // run eureka client async
    eureka.DefaultClient.Config(config).
        Register("APP_ID_CLIENT_FROM_CONFIG", 9000).
        Run()

    select {}
}
