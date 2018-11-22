package main

import (
    "github.com/HikoQiu/go-eureka-client/eureka"
    "log"
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

    c := eureka.DefaultClient.Config(config)
    api, err := c.Api()
    if err != nil {
        log.Fatalln("Failed to pick EurekaServerApi instance, err=", err.Error())
    }
    instances, err := api.QueryAllInstances()
    if err != nil {
        log.Fatalln("Failed to query all instances, err=", err.Error())
    }

    log.Println("all instances: ", instances)
}
