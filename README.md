# Kubernetes 1.16 API Deprecation Sonobuoy Plugin

This is pretty ugly right now, basically just a POC. It's actually much harder than you'd think to figure out which API version an object was applied with. This plugin looks at objects that have a `kubectl.kubernetes.io/last-applied-configuration` annotation to determine the API version at application time.
