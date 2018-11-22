package eureka

type DiscoveryClient interface {
    GetRegistryApps() map[string]ApplicationVo
}
