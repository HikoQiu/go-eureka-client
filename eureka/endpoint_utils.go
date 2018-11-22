package eureka

import (
    "fmt"
    "strings"
    "errors"
)

type EndpointUtils struct {
}

func (t *EndpointUtils) GetDiscoveryServiceUrls(config *EurekaClientConfig, zone string) ([]string, error) {
    if config.UseDnsForFetchingServiceUrls {
        return t.GetServiceUrlsFromDNS(config, zone)
    }

    return t.GetServiceUrlsFromConfig(config, zone)

}

/**
 * Get the zone based CNAMES that are bound to a region.
 *
 * @param region
 *            - The region for which the zone names need to be retrieved
 * @return - The list of CNAMES from which the zone-related information can
 *         be retrieved
 */
func (t *EndpointUtils) getZoneBasedDiscoveryUrlsFromRegion(config *EurekaClientConfig, region string) (map[string][]string, error) {
    discoveryDnsName := fmt.Sprintf("txt.%s.%s", region, config.EurekaServerDNSName)
    zoneCNames, _, err := lookupTXT(discoveryDnsName)
    if err != nil {
        log.Errorf("LookupTXT failed, err=%s", err.Error())
        return nil, err
    }

    zoneCnameSets := map[string][]string{}
    for _, zoneCname := range zoneCNames {
        zone := ""
        cnameTokens := strings.Split(zoneCname, ".")
        zone = cnameTokens[0]

        if _, ok := zoneCnameSets[zone]; !ok {
            zoneCnameSets[zone] = make([]string, 0)
        }

        zoneCnameSets[zone] = append(zoneCnameSets[zone], zoneCname)
    }
    return zoneCnameSets, nil
}

// @TODO GetServiceUrlsMapFromConfig
func (t *EndpointUtils) GetServiceUrlsMapFromConfig(config *EurekaClientConfig, instanceZone string) (map[string][]string, error) {

    return nil, nil
}

/**
 * Get the list of all eureka service urls from DNS for the eureka client to
 * talk to. The client picks up the service url from its zone and then fails over to
 * other zones randomly. If there are multiple servers in the same zone, the client once
 * again picks one randomly. This way the traffic will be distributed in the case of failures.
 *
 * @param clientConfig the clientConfig to use
 * @param instanceZone The zone in which the client resides.
 *
 * @return The list of all eureka service urls for the eureka client to talk to.
 */
func (t *EndpointUtils) GetServiceUrlsFromDNS(config *EurekaClientConfig, instanceZone string) ([]string, error) {
    zoneCnameSets, err := t.getZoneBasedDiscoveryUrlsFromRegion(config, config.GetRegion())
    if err != nil {
        return nil, err
    }

    if len(zoneCnameSets) == 0 {
        log.Errorf("No available zones configured for the instanceZone, instanceZone: %s", instanceZone);
        return nil, err
    }

    // Lookup all zone's service url
    zoneServiceUrls := map[string][]string{}
    for zone, cnames := range zoneCnameSets {
        for _, cname := range cnames {
            dnsName := fmt.Sprintf("txt.%s", cname)
            records, _, err := lookupTXT(dnsName)
            if err != nil {
                log.Errorf("LookupTXT failed, dnsName=%s, err=%s", dnsName, err.Error())
                return nil, err
            }

            if _, ok := zoneServiceUrls[zone]; !ok {
                zoneServiceUrls[zone] = make([]string, 0)
            }
            zoneServiceUrls[zone] = append(zoneServiceUrls[zone], records...)
        }
    }

    // while same zone eureka exist
    if _, ok := zoneServiceUrls[instanceZone]; ok && config.PreferSameZoneEureka {
        return t.formatUrls(config, zoneServiceUrls[instanceZone]), nil
    }

    // return service urls randomly
    for _, urls := range zoneServiceUrls {
        return t.formatUrls(config, urls), nil
    }

    err = errors.New("Fail to match service urls.")
    log.Errorf(err.Error())
    return nil, err
}

func (t *EndpointUtils) formatUrls(config *EurekaClientConfig, urls []string) []string {
    for i, _ := range urls {
        urls[i] = fmt.Sprintf("http://%s:%s/%s", urls[i], config.EurekaServerPort, config.EurekaServerUrlContext)
    }

    return urls
}

/**
 * Get the list of all eureka service urls from properties file for the eureka client to talk to.
 *
 * @param clientConfig the clientConfig to use
 * @param instanceZone The zone in which the client resides
 * @return The list of all eureka service urls for the eureka client to talk to
 */
func (t *EndpointUtils) GetServiceUrlsFromConfig(config *EurekaClientConfig, instanceZone string) ([]string, error) {
    availZones := config.GetAvailabilityZones(config.GetRegion())
    log.Debugf("The availability zone for the given region %s are %v ", config.GetRegion(), availZones)

    urls := make([]string, 0)
    for _, zone := range availZones {
        if _, ok := config.ServiceUrl[zone]; !ok {
            continue
        }

        zoneUrls := strings.Split(config.ServiceUrl[zone], ",")
        urls = append(urls, zoneUrls...)
    }

    return urls, nil
}
