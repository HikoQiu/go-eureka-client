package eureka

import (
    "testing"
    "strings"
)

func Test_LookupTXT(t *testing.T) {
    records, duration, err := lookupTXT("txt.zone-cn-hz-1.dev.ms-registry.xf.io")
    if err != nil {
        t.Fatal(err.Error())
    }
    t.Log("Records: ", records, duration)
}

func getTestDnsEurekaConfig() *EurekaClientConfig {
    config := GetDefaultEurekaClientConfig()
    config.UseDnsForFetchingServiceUrls = true
    config.Region = "region-cn-hd-1"
    config.AvailabilityZones = map[string]string{
        "region-cn-hd-1": "zone-cn-hz-1",
    }
    config.EurekaServerDNSName = "dev.ms-registry.xf.io"
    config.EurekaServerUrlContext = "eureka"
    config.EurekaServerPort = "9001"
    return config
}

func getTestConfiguredEurekaConfig() *EurekaClientConfig {
    config := GetDefaultEurekaClientConfig()
    config.UseDnsForFetchingServiceUrls = false
    config.Region = "region-cn-hd-1"
    config.AvailabilityZones = map[string]string{
        "region-cn-hd-1": "zone-cn-hz-1",
    }
    config.ServiceUrl = map[string]string {
       "zone-cn-hz-1": "http://192.168.20.236:9001/eureka,http://192.168.20.237:9001/eureka" ,
    }

    return config
}

// Get service urls by dns
func Test_GetServiceUrlsByDns(t *testing.T) {
    config := getTestDnsEurekaConfig()
    endpointUtils := new(EndpointUtils)
    urls, err := endpointUtils.GetDiscoveryServiceUrls(config, "zone-cn-hz-1")
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Eureka server urls: ", strings.Join(urls, ","))
}

// Get service urls by config
func Test_GetServiceUrlsByConfig(t *testing.T) {
    config := getTestConfiguredEurekaConfig()
    endpointUtils := new(EndpointUtils)
    urls, err := endpointUtils.GetDiscoveryServiceUrls(config, "zone-cn-hz-1")
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Eureka server urls: ", strings.Join(urls, ","))
}
