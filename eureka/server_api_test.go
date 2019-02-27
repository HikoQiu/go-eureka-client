package eureka

import (
    "encoding/json"
    "testing"
)

const (
    test_app_name      = "test-go-eureka-client"
    test_instance_port = 9000
)

func Test_RegisterInstance(t *testing.T) {
    config := getTestDnsEurekaConfig()
    endpointUtils := new(EndpointUtils)
    urls, err := endpointUtils.GetDiscoveryServiceUrls(config, "zone-cn-hz-1")
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Eureka server url: ", urls[0])
    instanceId, err := NewEurekaServerApi(urls[0]).RegisterInstance(test_app_name, test_instance_port)
    if err != nil {
        t.Error("Failed to register app: ", err.Error())
    }

    t.Log("Success to register app, instance-id: ", instanceId)
}

func Test_QueryAllInstances(t *testing.T) {
    config := getTestDnsEurekaConfig()
    endpointUtils := new(EndpointUtils)
    urls, err := endpointUtils.GetDiscoveryServiceUrls(config, "zone-cn-hz-1")
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Eureka server url: ", urls[0])
    applications, err := NewEurekaServerApi(urls[0]).QueryAllInstances()
    if err != nil {
        t.Fatal("Failed to query all intances: ", err.Error())
    }
    str, _ := json.Marshal(applications)
    t.Log("Success to query all instances, applications: ", string(str))
}

func Test_SendHeartbeat(t *testing.T) {
    config := getTestDnsEurekaConfig()
    endpointUtils := new(EndpointUtils)
    urls, err := endpointUtils.GetDiscoveryServiceUrls(config, "zone-cn-hz-1")
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Eureka server url: ", urls[0])

    // TEST NOT-FOUND
    err = NewEurekaServerApi(urls[0]).SendHeartbeat("NOT-FOUND", "NOT-FOUND")
    if err != nil {
        t.Log("NOT-FOUND-TEST: ", err.Error())
    }

    // TEST APP
    err = NewEurekaServerApi(urls[0]).SendHeartbeat(test_app_name, genDefaultInstanceId(test_app_name, test_instance_port))
    if err != nil {
        t.Log("TEST heartbeat, app=", test_app_name, ", err=", err.Error())
    }

    t.Log("SUCCESS")
}

func Test_UpdateInstanceStatus(t *testing.T) {
    config := getTestDnsEurekaConfig()
    endpointUtils := new(EndpointUtils)
    urls, err := endpointUtils.GetDiscoveryServiceUrls(config, "zone-cn-hz-1")
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Eureka server url: ", urls[0])

    // 1.1 register instance
    instanceId, err := NewEurekaServerApi(urls[0]).RegisterInstance(test_app_name, test_instance_port)
    if err != nil {
        t.Error("Failed to register app: ", err.Error())
    }

    // 1.2 update instance's status
    err = NewEurekaServerApi(urls[0]).UpdateInstanceStatus(test_app_name, instanceId, STATUS_UP)
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Success to update instance status")
}

func Test_QueryAllInstanceByAppId(t *testing.T) {
    config := getTestDnsEurekaConfig()
    endpointUtils := new(EndpointUtils)
    urls, err := endpointUtils.GetDiscoveryServiceUrls(config, "zone-cn-hz-1")
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Eureka server url: ", urls[0])

    instances, err := NewEurekaServerApi(urls[0]).QueryAllInstanceByAppId(test_app_name)
    if err != nil {
        t.Fatal("Failed to register app: ", err.Error())
    }

    t.Log(test_app_name+" instances: ", instances)
}

func Test_QuerySpecificAppInstance(t *testing.T) {
    config := getTestDnsEurekaConfig()
    endpointUtils := new(EndpointUtils)
    urls, err := endpointUtils.GetDiscoveryServiceUrls(config, "zone-cn-hz-1")
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Eureka server url: ", urls[0])

    instance, err := NewEurekaServerApi(urls[0]).QuerySpecificAppInstance(genDefaultInstanceId(test_app_name, test_instance_port))
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("instance: ", instance)
}

func Test_UpdateMeta(t *testing.T) {
    config := getTestDnsEurekaConfig()
    endpointUtils := new(EndpointUtils)
    urls, err := endpointUtils.GetDiscoveryServiceUrls(config, "zone-cn-hz-1")
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Eureka server url: ", urls[0])

    err = NewEurekaServerApi(urls[0]).UpdateMeta(test_app_name,
        genDefaultInstanceId(test_app_name, test_instance_port),
        map[string]string{
            "key": "value",
        })
    if err != nil {
        t.Fatal(err.Error())
    }

    t.Log("Succeed to update meta.")
}
