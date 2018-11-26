package eureka

import (
    "time"
    "sync"
    "syscall"
    "os"
    "os/signal"
    "math/rand"
    "errors"
)

const (
    DEFAULT_SLEEP_INTERVALS = 3
)

var DefaultClient = new(Client)

// eureka client
type Client struct {
    // eureka client config
    config *EurekaClientConfig

    // current client (instance) config
    instance *InstanceVo

    // eureka server base url list
    serviceUrls []string

    // applications int registry
    // key: appId
    // value: ApplicationVo
    registryApps map[string]ApplicationVo

    // for monitor system signal
    signalChan chan os.Signal

    mu sync.RWMutex
}

func (t *Client) Config(config *EurekaClientConfig) *Client {
    t.config = config
    return t
}

// user brief parameters to register instance
func (t *Client) Register(appId string, port int) *Client {
    vo := DefaultInstanceVo()
    vo.App = appId
    vo.Status = STATUS_STARTING
    vo.Port = positiveInt{Value: port, Enabled: "true"}
    t.instance = vo
    return t
}

// user raw instanceVo to register instance
func (t *Client) RegisterVo(vo *InstanceVo) *Client {
    t.instance = vo
    return t
}

// Api for sending rest http to eureka server
func (t *Client) Api() (*EurekaServerApi, error) {
    api, err := t.pickEurekaServerApi()
    if err != nil {
        return nil, err
    }

    return api, nil
}

func (t *Client) GetInstance() *InstanceVo {
    return t.instance
}

func (t *Client) GetRegistryApps() map[string]ApplicationVo {
    t.mu.RLock()
    defer t.mu.RUnlock()

    return t.registryApps
}

// start eureka client
// 1. parse/get service urls
// 2. register client to eureka server and send heartbeat
func (t *Client) Run() {
    err := t.refreshServiceUrls()
    if err != nil {
        log.Errorf("Failed to refresh service urls, err=%s", err.Error())
        return
    }

    // handle exit signal to de-register instance
    go t.handleSignal()

    // (if FetchRegistry is true), fetch registry apps periodically
    // and update to t.registryApps
    go t.refreshRegistry()

    t.registerWithEureka()
}

func (t *Client) refreshServiceUrls() error {
    err := t.getServiceUrlsWithZones()
    if err != nil {
        log.Errorf("Failed to init service urls, err=%s", err.Error())
        return err
    }

    // auto update service urls
    // (only) while userDnsForFetchingServiceUrls=true and AutoUpdateDnsServiceUrls=true
    go func() {
        if !t.config.UseDnsForFetchingServiceUrls && t.config.AutoUpdateDnsServiceUrls {
            return
        }

        for {
            t.getServiceUrlsWithZones()

            time.Sleep(time.Duration(t.config.AutoUpdateDnsServiceUrlsIntervals) * time.Second)
            log.Errorf("AutoUpdateDnsServiceUrls... ok")
        }
    }()

    return nil
}

func (t *Client) getServiceUrlsWithZones() error {
    availZones := t.config.GetAvailabilityZones(t.config.Region)
    endpointUtils := new(EndpointUtils)

    // loop to get zone's service urls
    var err error
    var urls []string
    for _, zone := range availZones {
        urls, err = endpointUtils.GetDiscoveryServiceUrls(t.config, zone)
        if err != nil {
            log.Errorf("Failed to boot eureka client, zone=%s, err=%s", zone, err.Error())
            continue
        }

        t.mu.Lock()
        t.mu.Unlock()
        t.serviceUrls = urls
        break
    }

    return err
}

// rand to pick service url
func (t *Client) pickServiceUrl() (string, bool) {
    if len(t.serviceUrls) == 0 {
        // if serviceUrls not init, try to fetch service urls one time
        err := t.getServiceUrlsWithZones()
        if err != nil {
            return "", false
        }
    }

    t.mu.RLock()
    defer t.mu.RUnlock()

    index := 0
    if len(t.serviceUrls) > 1 {
        index = rand.Intn(len(t.serviceUrls) - 1)
    }
    return t.serviceUrls[index], true
}

// rand to pick service url and new EurekaServerApi instance
func (t *Client) pickEurekaServerApi() (*EurekaServerApi, error) {
    url, ok := t.pickServiceUrl()
    if !ok {
        log.Errorf("No service url is available to pick.")
        return nil, errors.New("No service url is available to pick.")
    }

    return NewEurekaServerApi(url), nil
}

// register instance (default current status is STARTING)
// and update instance status to UP
func (t *Client) registerWithEureka() {
    if !t.config.RegisterWithEureka {
        return
    }

    // ensure client succeed to register to eureka server
    for {
        if t.instance == nil {
            log.Errorf("Eureka instance can't be nil")
            return
        }

        api, err := t.Api()
        if err != nil {
            time.Sleep(time.Second * DEFAULT_SLEEP_INTERVALS)
            continue
        }

        instanceId, err := api.RegisterInstanceWithVo(t.instance)
        if err != nil {
            log.Errorf("Client register failed, err=%s", err.Error())
            time.Sleep(time.Second * DEFAULT_SLEEP_INTERVALS)
            continue
        }
        t.instance.InstanceId = instanceId

        err = api.UpdateInstanceStatus(t.instance.App, t.instance.InstanceId, STATUS_UP)
        if err != nil {
            log.Errorf("Client UP failed, err=%s", err.Error())
            time.Sleep(time.Second * DEFAULT_SLEEP_INTERVALS)
            continue
        }

        // if success to register to eureka and update status tu UP
        // then break loop
        break;
    }

    // send heartbeat
    t.heartbeat()
}

// eureka client heartbeat
func (t *Client) heartbeat() {
    go func() {
        for {
            api, err := t.Api()
            if err != nil {
                time.Sleep(time.Second * DEFAULT_SLEEP_INTERVALS)
                continue
            }

            err = api.SendHeartbeat(t.instance.App, t.instance.InstanceId)
            if err != nil {
                log.Errorf("Failed to send heartbeat, err=%s", err.Error())
                time.Sleep(time.Second * DEFAULT_SLEEP_INTERVALS)
                continue
            }

            log.Debugf("Heartbeat app=%s, instanceId=%s", t.instance.App, t.instance.InstanceId)
            time.Sleep(time.Duration(t.config.HeartbeatIntervals) * time.Second)
        }
    }()
}

func (t *Client) refreshRegistry() {
    if !t.config.FetchRegistry {
        return
    }

    for {
        t.fetchRegistry()
        time.Sleep(time.Second * time.Duration(t.config.RegistryFetchIntervalSeconds))
    }
}

func (t *Client) fetchRegistry() (map[string]ApplicationVo, error) {
    api, err := t.Api()
    if err != nil {
        log.Errorf("Failed to QueryAllInstances, err=%s", err.Error())
        return nil, err
    }

    apps, err := api.QueryAllInstances()
    if err != nil {
        log.Errorf("Failed to QueryAllInstances, err=%s", err.Error())
        return nil, err
    }

    t.mu.Lock()
    defer t.mu.Unlock()

    // @TODO  FilterOnlyUpInstances  true,

    t.registryApps = make(map[string]ApplicationVo)
    for _, app := range apps {
        t.registryApps[app.Name] = app
    }

    return t.registryApps, nil
}

// for graceful kill. Here handle SIGTERM signal to do sth
// e.g: kill -TERM $pid
//      or "ctrl + c" to exit
func (t *Client) handleSignal() {
    if t.signalChan == nil {
        t.signalChan = make(chan os.Signal)
    }

    signal.Notify(t.signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

    for {
        switch <-t.signalChan {
        case syscall.SIGINT:
            fallthrough
        case syscall.SIGKILL:
            fallthrough
        case syscall.SIGTERM:
            log.Infof("Receive exit signal, client instance going to de-register, instanceId=%s.", t.instance.InstanceId)

            // de-register instance
            api, err := t.Api()
            if err != nil {
                log.Errorf("Failed to get EurekaServerApi instance, de-register %s failed, err=%s", t.instance.InstanceId, err.Error())
                return
            }

            err = api.DeRegisterInstance(t.instance.App, t.instance.InstanceId)
            if err != nil {
                log.Errorf("Failed to de-register %s, err=%s", t.instance.InstanceId, err.Error())
                return
            }

            log.Infof("de-register %s success.", t.instance.InstanceId)
            os.Exit(0)
        }
    }
}
