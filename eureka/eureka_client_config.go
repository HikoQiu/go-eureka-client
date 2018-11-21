package eureka

import "strings"

const (
    DEFAULT_REGION = "default"
    DEFAULT_PREFIX = "/eureka"
    DEFAULT_ZONE   = "defaultZone"
)

// refer to:
//https://github.com/spring-cloud/spring-cloud-netflix/blob/master/spring-cloud-netflix-eureka-client/src/main/java/org/springframework/cloud/netflix/eureka/EurekaClientConfigBean.java
type EurekaClientConfig struct {
    /**
	 * Indicates how often(in seconds) to fetch the registry information from the eureka
	 * server.
	 */
    RegistryFetchIntervalSeconds int
    /**
	 * Indicates how often(in seconds) to replicate instance changes to be replicated to
	 * the eureka server.
	 */
    //InstanceInfoReplicationIntervalSeconds int

    /**
     * Indicates how long initially (in seconds) to replicate instance info to the eureka
     * server
     */
    //InitialInstanceInfoReplicationIntervalSeconds int

    /**
     * Indicates how often(in seconds) to poll for changes to eureka server information.
     * Eureka servers could be added or removed and this setting controls how soon the
     * eureka clients should know about it.
     */
    //EurekaServiceUrlPollIntervalSeconds int

    /**
     * Gets the proxy port to eureka server if any.
     */
    //ProxyPort string

    /**
     * Gets the proxy host to eureka server if any.
     */
    //ProxyHost string

    /**
     * Gets the proxy user name if any.
     */
    //ProxyUserName string

    /**
     * Gets the proxy password if any.
     */
    //ProxyPassword string

    /**
     * Indicates how long to wait (in seconds) before a read from eureka server needs to
     * timeout.
     */
    //EurekaServerReadTimeoutSeconds int

    /**
     * Indicates how long to wait (in seconds) before a connection to eureka server needs
     * to timeout. Note that the connections in the client are pooled by
     * org.apache.http.client.HttpClient and this setting affects the actual connection
     * creation and also the wait time to get the connection from the pool.
     */
    //EurekaServerConnectTimeoutSeconds int

    /**
     * Gets the name of the implementation which implements BackupRegistry to fetch the
     * registry information as a fall back option for only the first time when the eureka
     * client starts.
     *
     * This may be needed for applications which needs additional resiliency for registry
     * information without which it cannot operate.
     */
    //BackupRegistryImpl string

    /**
     * Gets the total number of connections that is allowed from eureka client to all
     * eureka servers.
     */
    //EurekaServerTotalConnections int

    /**
     * Gets the total number of connections that is allowed from eureka client to a eureka
     * server host.
     */
    //EurekaServerTotalConnectionsPerHost int

    /**
     * Gets the URL context to be used to construct the service url to contact eureka
     * server when the list of eureka servers come from the DNS. This information is not
     * required if the contract returns the service urls from eurekaServerServiceUrls.
     *
     * The DNS mechanism is used when useDnsForFetchingServiceUrls is set to true and the
     * eureka client expects the DNS to configured a certain way so that it can fetch
     * changing eureka servers dynamically. The changes are effective at runtime.
     */
    EurekaServerUrlContext string

    /**
     * Gets the port to be used to construct the service url to contact eureka server when
     * the list of eureka servers come from the DNS.This information is not required if
     * the contract returns the service urls eurekaServerServiceUrls(String).
     *
     * The DNS mechanism is used when useDnsForFetchingServiceUrls is set to true and the
     * eureka client expects the DNS to configured a certain way so that it can fetch
     * changing eureka servers dynamically.
     *
     * The changes are effective at runtime.
     */
    EurekaServerPort string

    /**
     * Gets the DNS name to be queried to get the list of eureka servers.This information
     * is not required if the contract returns the service urls by implementing
     * serviceUrls.
     *
     * The DNS mechanism is used when useDnsForFetchingServiceUrls is set to true and the
     * eureka client expects the DNS to configured a certain way so that it can fetch
     * changing eureka servers dynamically.
     *
     * The changes are effective at runtime.
     */
    EurekaServerDNSName string

    /**
     * Gets the region (used in AWS datacenters) where this instance resides.
     */
    Region string

    /**
     * Indicates how much time (in seconds) that the HTTP connections to eureka server can
     * stay idle before it can be closed.
     *
     * In the AWS environment, it is recommended that the values is 30 seconds or less,
     * since the firewall cleans up the connection information after a few mins leaving
     * the connection hanging in limbo
     */
    //EurekaConnectionIdleTimeoutSeconds int

    /**
     * Indicates whether the client is only interested in the registry information for a
     * single VIP.
     */
    //RegistryRefreshSingleVipAddress string

    /**
     * The thread pool size for the heartbeatExecutor to initialise with
     */
    //HeartbeatExecutorThreadPoolSize int

    /**
     * Heartbeat executor exponential back off related property. It is a maximum
     * multiplier value for retry delay, in case where a sequence of timeouts occurred.
     */
    //HeartbeatExecutorExponentialBackOffBound int

    /**
     * The thread pool size for the cacheRefreshExecutor to initialise with
     */
    //CacheRefreshExecutorThreadPoolSize int

    /**
     * Cache refresh executor exponential back off related property. It is a maximum
     * multiplier value for retry delay, in case where a sequence of timeouts occurred.
     */
    //CacheRefreshExecutorExponentialBackOffBound int

    /**
     * Map of availability zone to list of fully qualified URLs to communicate with eureka
     * server. Each value can be a single URL or a comma separated list of alternative
     * locations.
     *
     * Typically the eureka server URLs carry protocol,host,port,context and version
     * information if any. Example:
     * http://ec2-256-156-243-129.compute-1.amazonaws.com:7001/eureka/
     *
     * The changes are effective at runtime at the next service url refresh cycle as
     * specified by eurekaServiceUrlPollIntervalSeconds.
     */
    ServiceUrl map[string]string

    /**
     * Indicates whether the content fetched from eureka server has to be compressed
     * whenever it is supported by the server. The registry information from the eureka
     * server is compressed for optimum network traffic.
     */
    //GZipContent bool

    /**
     * Indicates whether the eureka client should use the DNS mechanism to fetch a list of
     * eureka servers to talk to. When the DNS name is updated to have additional servers,
     * that information is used immediately after the eureka client polls for that
     * information as specified in eurekaServiceUrlPollIntervalSeconds.
     *
     * Alternatively, the service urls can be returned serviceUrls, but the users should
     * implement their own mechanism to return the updated list in case of changes.
     *
     * The changes are effective at runtime.
     */
    UseDnsForFetchingServiceUrls bool

    /**
     * Indicates whether or not this instance should register its information with eureka
     * server for discovery by others.
     *
     * In some cases, you do not want your instances to be discovered whereas you just
     * want do discover other instances.
     */
    RegisterWithEureka bool

    /**
     * Indicates whether or not this instance should try to use the eureka server in the
     * same zone for latency and/or other reason.
     *
     * Ideally eureka clients are configured to talk to servers in the same zone
     *
     * The changes are effective at runtime at the next registry fetch cycle as specified
     * by registryFetchIntervalSeconds
     */
    PreferSameZoneEureka bool

    /**
     * Indicates whether to log differences between the eureka server and the eureka
     * client in terms of registry information.
     *
     * Eureka client tries to retrieve only delta changes from eureka server to minimize
     * network traffic. After receiving the deltas, eureka client reconciles the
     * information from the server to verify it has not missed out some information.
     * Reconciliation failures could happen when the client has had network issues
     * communicating to server.If the reconciliation fails, eureka client gets the full
     * registry information.
     *
     * While getting the full registry information, the eureka client can log the
     * differences between the client and the server and this setting controls that.
     *
     * The changes are effective at runtime at the next registry fetch cycle as specified
     * by registryFetchIntervalSecondsr
     */
    //LogDeltaDiff bool

    /**
     * Indicates whether the eureka client should disable fetching of delta and should
     * rather resort to getting the full registry information.
     *
     * Note that the delta fetches can reduce the traffic tremendously, because the rate
     * of change with the eureka server is normally much lower than the rate of fetches.
     *
     * The changes are effective at runtime at the next registry fetch cycle as specified
     * by registryFetchIntervalSeconds
     */
    //DisableDelta bool

    /**
     * Comma separated list of regions for which the eureka registry information will be
     * fetched. It is mandatory to define the availability zones for each of these regions
     * as returned by availabilityZones. Failing to do so, will result in failure of
     * discovery client startup.
     *
     */
    FetchRemoteRegionsRegistry string

    /**
     * Gets the list of availability zones (used in AWS data centers) for the region in
     * which this instance resides.
     *
     * The changes are effective at runtime at the next registry fetch cycle as specified
     * by registryFetchIntervalSeconds.
     */
    AvailabilityZones map[string]string

    /**
     * Indicates whether to get the applications after filtering the applications for
     * instances with only InstanceStatus UP states.
     */
    FilterOnlyUpInstances bool

    /**
     * Indicates whether this client should fetch eureka registry information from eureka
     * server.
     */
    FetchRegistry bool

    /**
     * Get a replacement string for Dollar sign <code>$</code> during
     * serializing/deserializing information in eureka server.
     */
    //DollarReplacement string

    /**
     * Get a replacement string for underscore sign <code>_</code> during
     * serializing/deserializing information in eureka server.
     */
    //EscapeCharReplacement string

    /**
     * Indicates whether server can redirect a client request to a backup server/cluster.
     * If set to false, the server will handle the request directly, If set to true, it
     * may send HTTP redirect to the client, with a new server location.
     */
    //AllowRedirects bool

    /**
     * If set to true, local status updates via ApplicationInfoManager will trigger
     * on-demand (but rate limited) register/updates to remote eureka servers
     */
    //OnDemandUpdateStatusChange bool

    /**
     * This is a transient config and once the latest codecs are stable, can be removed
     * (as there will only be one)
     */
    //EncoderName string

    /**
     * This is a transient config and once the latest codecs are stable, can be removed
     * (as there will only be one)
     */
    //decoderName string

    /**
     * Indicates whether the client should explicitly unregister itself from the remote server
     * on client shutdown.
     */
    //ShouldUnregisterOnShutdown bool

    /**
     * Indicates whether the client should enforce registration during initialization. Defaults to false.
     */
    //ShouldEnforceRegistrationAtInit bool

    /**
     * Order of the discovery client used by `CompositeDiscoveryClient` for sorting available clients.
     */
    //Order int

    //
    // extend features
    //
    //
    // (only when UseDnsForFetchingServiceUrls=true effects)
    // Auto lookup dns to update service urls
    AutoUpdateDnsServiceUrls bool

    // when UseDnsForFetchingServiceUrls=true and AutoUpdateDnsServiceUrls=true
    // AutoUpdateDnsServiceUrlsIntervals effects
    // default value: 5*60 seconds
    AutoUpdateDnsServiceUrlsIntervals int

    // eureka client heartbeat intervals
    // Tips:
    // 1. only when RegisterWithEureka=true, HeartbeatIntervals effects
    // 2. HeartbeatIntervals must less than EvictionDurationInSecs(in server_api_vos.go, InstanceVo.LeaseInfo.EvictionDurationInSecs)
    HeartbeatIntervals int
}

// get default config
func GetDefaultEurekaClientConfig() *EurekaClientConfig {
    return &EurekaClientConfig{
        AvailabilityZones: map[string]string{},
        ServiceUrl: map[string]string{
            DEFAULT_ZONE: "http://localhost:8761" + DEFAULT_PREFIX,
        },
        UseDnsForFetchingServiceUrls: false,
        RegisterWithEureka:           true,
        PreferSameZoneEureka:         true,
        FilterOnlyUpInstances:        true,
        RegistryFetchIntervalSeconds: 30,
        FetchRegistry:                true,
        EurekaServerPort:             "8761",
        EurekaServerUrlContext:       "eureka",

        // extend features
        AutoUpdateDnsServiceUrls:          true,
        AutoUpdateDnsServiceUrlsIntervals: 5 * 60,
        HeartbeatIntervals:                30,

        // @TODO Features not implement
        //InstanceInfoReplicationIntervalSeconds:        30,
        //InitialInstanceInfoReplicationIntervalSeconds: 40,
        //EurekaServiceUrlPollIntervalSeconds:           5 * 60,
        //EurekaServerReadTimeoutSeconds:                8,
        //EurekaServerConnectTimeoutSeconds:             5,
        //EurekaServerTotalConnections:                  200,
        //EurekaServerTotalConnectionsPerHost:           50,
        //EurekaConnectionIdleTimeoutSeconds:            30,
        //HeartbeatExecutorThreadPoolSize:               2,
        //HeartbeatExecutorExponentialBackOffBound:      10,
        //CacheRefreshExecutorThreadPoolSize:            2,
        //CacheRefreshExecutorExponentialBackOffBound:   10,
        //GZipContent:                     true,
        //DollarReplacement:               "_-",
        //EscapeCharReplacement:           "__",
        //AllowRedirects:                  false,
        //OnDemandUpdateStatusChange:      true,
        //ShouldUnregisterOnShutdown:      true,
        //ShouldEnforceRegistrationAtInit: false,
    }
}

func (t *EurekaClientConfig) GetRegion() string {
    if t.Region == "" {
        return DEFAULT_REGION
    }

    return strings.ToLower(t.Region)
}

func (t *EurekaClientConfig) GetAvailabilityZones(region string) []string {
    if _, ok := t.AvailabilityZones[region]; ok {
        return strings.Split(t.AvailabilityZones[region], ",")
    }

    return []string{DEFAULT_ZONE}
}
