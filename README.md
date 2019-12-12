# Drone plugin for Helm 3

TODO:

* [x] Make a `.drone.yml` that's sufficient for building an image
* [x] Make a `Dockerfile` that's sufficient for launching the built image
* [x] Make `cmd/drone-helm/main.go` actually invoke `helm`
* [x] Make `golint` part of the build process (and make it pass)
* [x] Implement debug output
* [x] Flesh out `helm upgrade` until it's capable of working
* [ ] Implement config settings for `upgrade`
* [ ] Implement `helm lint`
* [ ] Implement `helm delete`
* [ ] EKS support
* [ ] Change `.drone.yml` to use a real docker registry
