package eureka

import (
    "net/http"
    "github.com/go-resty/resty"
    "time"
    "github.com/kataras/iris/core/errors"
    "fmt"
    "strings"
    "encoding/json"
)

// Refer to: https://github.com/Netflix/eureka/wiki/Eureka-REST-operations
type EurekaServerApi struct {
    BaseUrl string
}

func NewEurekaServerApi(baseUrl string) *EurekaServerApi {
    return &EurekaServerApi{
        BaseUrl: baseUrl,
    }
}

func (t *EurekaServerApi) url(path string) string {
    return strings.TrimRight(t.BaseUrl, "/") + path
}

// http rest request encapsulated simply
// e.g: request(method, url, params, header)
func (t *EurekaServerApi) request(method, url string, args ...interface{}) (*resty.Response, error) {
    var params interface{}
    header := map[string]string{}
    switch len(args) {
    case 1:
        params = args[0]
    case 2:
        header = args[1].(map[string]string)
    }
    header["content-type"] = "application/json"

    var res *resty.Response
    var err error
    req := resty.New().SetTimeout(time.Second * 10).R().SetHeaders(header);
    switch method {
    case http.MethodGet:
        res, err = req.Get(url)
    case http.MethodPost:
        res, err = req.SetBody(params).Post(url)
    case http.MethodPut:
        res, err = req.SetBody(params).Put(url)
    case http.MethodDelete:
        res, err = req.Delete(url)
    default:
        return nil, errors.New("Failed to recognize method: " + method)
    }

    if err != nil {
        return nil, err
    }
    if res.StatusCode() >= 300 {
        return nil, errors.New(fmt.Sprintf("Request failed, Http status code: %d, body: %s", res.StatusCode(), string(res.Body())))
    }

    return res, err
}

// Register new application instance by brief info
func (t *EurekaServerApi) RegisterInstance(appId string, port int) (string, error) {
    vo := DefaultInstanceVo()
    vo.App = appId
    vo.Status = STATUS_STARTING
    vo.Port = positiveInt{Value: port, Enabled: "true"}

    return t.RegisterInstanceWithVo(vo)
}

// Register new application instance
func (t *EurekaServerApi) RegisterInstanceWithVo(vo *InstanceVo) (string, error) {
    if vo.HomePageUrl == "" {
        vo.HomePageUrl = fmt.Sprintf("http://%s:%d/", vo.IppAddr, vo.Port.Value)
    }
    if vo.StatusPageUrl == "" {
        vo.StatusPageUrl = vo.HomePageUrl + "info"
    }

    if vo.HealthCheckUrl == "" {
        vo.HealthCheckUrl = vo.HomePageUrl + "health"
    }

    if vo.InstanceId == "" {
        vo.InstanceId = fmt.Sprintf("%s:%s:%d", vo.Hostname, vo.App, vo.Port.Value)
    }

    body, _ := json.Marshal(map[string]interface{}{"instance": vo})
    _, err := t.request(http.MethodPost, t.url("/apps/"+vo.App), body)
    if err != nil {
        log.Debugf("Failed to register app=%s, err=%s", vo.App, err.Error())
        return "", err
    }

    return vo.InstanceId, nil
}

// De-register application instance
func (t *EurekaServerApi) DeRegisterInstance(appId, instanceId string) error {
    _, err := t.request(http.MethodDelete, t.url(fmt.Sprintf("/apps/%s/%s", appId, instanceId)))
    if err != nil {
        log.Debugf("Failed to De-register application instance, err=%s", err.Error())
        return err
    }

    return nil
}

// Send application instance heartbeat
func (t *EurekaServerApi) SendHeartbeat(appId, instanceId string) error {
    _, err := t.request(http.MethodPut, t.url(fmt.Sprintf("/apps/%s/%s", appId, instanceId)))
    if err != nil {
        log.Debugf("Failed to send instance heartbeat, app-id=%s, instance-id=%s, err=%s", appId, instanceId, err.Error())
        return err
    }

    return nil
}

// Query for all instances
func (t *EurekaServerApi) QueryAllInstances() ([]ApplicationVo, error) {
    res, err := t.request(http.MethodGet, t.url("/apps"))
    if err != nil {
        log.Debugf("Failed to query all instances, err=%s", err.Error())
        return nil, err
    }

    resApps := make(map[string]ApplicationsVo)
    err = json.Unmarshal(res.Body(), &resApps)
    if err != nil {
        log.Debugf("Failed to query all instances, json.Unmarshal err=%s", err.Error())
        return nil, err
    }

    return resApps["applications"].Application, nil
}

// Query for all appId instances
func (t *EurekaServerApi) QueryAllInstanceByAppId(appId string) ([]InstanceVo, error) {
    appId = strings.ToUpper(appId)
    res, err := t.request(http.MethodGet, t.url("/apps/"+appId))
    if err != nil {
        log.Debugf("Failed to query appId instances, err=%s", err.Error())
        return nil, err
    }

    resIns := make(map[string]ApplicationVo)
    err = json.Unmarshal(res.Body(), &resIns)
    if err != nil {
        log.Debugf("Failed to query appId instances, json.Unmarshal err=%s", err.Error())
        return nil, err
    }

    return resIns["application"].Instances, nil
}

// query specific instanceId
func (t *EurekaServerApi) QuerySpecificAppInstance(instanceId string) (*InstanceVo, error) {
    res, err := t.request(http.MethodGet, t.url("/instances/"+instanceId))
    if err != nil {
        log.Debugf("Failed to query specific app instance, err=%s", err.Error())
        return nil, err
    }

    resIns := make(map[string]*InstanceVo)
    err = json.Unmarshal(res.Body(), &resIns)
    if err != nil {
        log.Debugf("Failed to query appId instances, json.Unmarshal err=%s", err.Error())
        return nil, err
    }

    return resIns["instance"], nil
}

// update instance status
func (t *EurekaServerApi) UpdateInstanceStatus(appId, instanceId, status string) error {
    _, err := t.request(http.MethodPut, t.url(fmt.Sprintf("/apps/%s/%s/status?value=%s", appId, instanceId, status)))
    if err != nil {
        log.Debugf("Failed to update instance status, err=%s", err.Error())
        return err
    }

    return nil
}

// Update meta data
func (t *EurekaServerApi) UpdateMeta(appId, instanceId string, meta map[string]string) error {
    queryStr := ""
    for k, v := range meta {
        queryStr += fmt.Sprintf("&%s=%s", k, v)
    }
    queryStr = strings.TrimLeft(queryStr, "&")

    _, err := t.request(http.MethodPut, t.url(fmt.Sprintf("/apps/%s/%s/metadata?%s", appId, instanceId, queryStr)))
    if err != nil {
        log.Debugf("Failed to update instance meta data, err=%s", err.Error())
        return err
    }

    return nil
}

// @TODO Query for all instances under a particular vip address
// Sorry, I don't have the environment for testing
func (t *EurekaServerApi) QueryAllVipInstances() {

}

// @TODO Query for all instances under a particular vip address
// Sorry, I don't have the environment for testing
func (t *EurekaServerApi) QueryAllSVipInstances() {

}
