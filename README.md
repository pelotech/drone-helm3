# Drone plugin for Helm 3

TODO:

* [x] Make a `.drone.yml` that's sufficient for building an image
* [x] Make a `Dockerfile` that's sufficient for launching the built image
* [x] Make `cmd/drone-helm/main.go` actually invoke `helm`
* [x] Make `golint` part of the build process (and make it pass)
* [x] Implement debug output
* [x] Flesh out `helm upgrade` until it's capable of working
* [x] Implement config settings for `upgrade`
* [ ] Implement `helm lint`
* [ ] Implement `helm delete`
* [ ] Look for command-line flags added in helm3; implement them
* [ ] EKS support
* [ ] Dotenv support
* [ ] Example drone config in this README
* [ ] Change `.drone.yml` to use a real docker registry

Nice-to-haves:

* [ ] Cleanup() method on Steps to close open filehandles, etc.
* [ ] Replace `fmt.Printf` with an actual logger
* [ ] Replace `fmt.Errorf` with `github.com/pkg/errors.Wrap`, since the built-in `Unwrap` doesn't work the way `Cause` does
* [ ] Deprecation warnings if there are environment variables that aren't applicable in helm3
