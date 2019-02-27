package eureka

const (
    STATUS_UP             = "UP"
    STATUS_DOWN           = "DOWN"
    STATUS_STARTING       = "STARTING"
    STATUS_OUT_OF_SERVICE = "OUT_OF_SERVICE"
    STATUS_UNKNOWN        = "UNKNOWN"

    DC_NAME_TYPE_MY_OWN = "MyOwn"
    DC_NAME_TYPE_AMAZON = "Amazon"
)

type (
    InstanceVo struct {
        // Register application instance needed -- BEGIN
        Hostname         string         `json:"hostName"`
        App              string         `json:"app"`
        IppAddr          string         `json:"ippAddr"`
        VipAddress       string         `json:"vipAddress"`
        SecureVipAddress string         `json:"secureVipAddress"`
        Status           string         `json:"status"`
        Port             positiveInt    `json:"port"`
        SecurePort       positiveInt    `json:"securePort"`
        HomePageUrl      string         `json:"homePageUrl"`
        StatusPageUrl    string         `json:"statusPageUrl"`
        HealthCheckUrl   string         `json:"healthCheckUrl"`
        DataCenterInfo   DataCenterInfo `json:"dataCenterInfo"`
        LeaseInfo        LeaseInfo      `json:"leaseInfo"`
        // Register application instance needed -- END

        InstanceId           string `json:"instanceId,omitempty"`
        OverriddenStatus     string `json:"overriddenstatus,omitempty"`
        LastUpdatedTimestamp int    `json:"lastUpdatedTimestamp,omitempty"`
        LastDirtyTimestamp   int    `json:"lastUpdatedTimestamp,omitempty"`
        ActionType           string `json:"actionType,omitempty"`
    }

    positiveInt struct {
        Value   int    `json:"$"`
        Enabled string `json:"@enabled"` // true|false
    }

    DataCenterInfo struct {
        // MyOwn | Amazon
        Name string `json:"name"`
        // metadata is only required if name is Amazon
        Metadata string `json:"metadata,omitempty"`
        Class    string `json:"@class"`
    }

    LeaseInfo struct {
        // (optional) if you want to change the length of lease - default if 90 seconds
        EvictionDurationInSecs int `json:"eviction_duration_in_secs,omitempty"`
    }

    // application
    ApplicationVo struct {
        Name      string       `json:"name"`
        Instances []InstanceVo `json:"instance"`
    }

    ApplicationsVo struct {
        VersionDelta string          `json:"version__delta"`
        AppsHashCode string          `json:"apps_hash__code"`
        Application  []ApplicationVo `json:"application"`
    }
)

func DefaultInstanceVo() *InstanceVo {
    ip := getLocalIp()
    //hostname, err := os.Hostname()
    //if err != nil {
    //    log.Errorf("Failed to get hostname, err=%s, user ip as hostname, ip=%s", err.Error(), ip)
    //    hostname = ip
    //}
    return &InstanceVo{
        //Hostname:         hostname,
        Hostname:         ip,
        App:              "",
        IppAddr:          ip,
        VipAddress:       ip,
        SecureVipAddress: ip,
        Status:           STATUS_STARTING,
        Port:             positiveInt{Value: 8080, Enabled: "true"},
        SecurePort:       positiveInt{Value: 443, Enabled: "false"},
        HomePageUrl:      "",
        StatusPageUrl:    "",
        HealthCheckUrl:   "",
        DataCenterInfo: DataCenterInfo{
            Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
            Name:  DC_NAME_TYPE_MY_OWN,
        },
        LeaseInfo: LeaseInfo{
            EvictionDurationInSecs: 30,
        },
    }
}
