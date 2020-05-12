# Drone plugin for Helm 3

[![Build Status](https://cloud.drone.io/api/badges/pelotech/drone-helm3/status.svg)](https://cloud.drone.io/pelotech/drone-helm3)
[![Go Report](https://goreportcard.com/badge/github.com/pelotech/drone-helm3)](https://goreportcard.com/report/github.com/pelotech/drone-helm3)
[![](https://images.microbadger.com/badges/image/pelotech/drone-helm3.svg)](https://microbadger.com/images/pelotech/drone-helm3 "Get your own image badge on microbadger.com")

This plugin provides an interface between [Drone](https://drone.io/) and [Helm 3](https://github.com/kubernetes/helm):

* Lint your charts
* Deploy your service
* Delete your service
* Publish it to a Helm Registry (OCI-based registries)

The plugin is inpsired by [drone-helm](https://github.com/ipedrazas/drone-helm), which fills the same role for Helm 2. It provides a comparable feature-set and the configuration settings are backward-compatible.

## Example configuration

The examples below give a minimal and sufficient configuration for each use-case. For a full description of each command's settings, see [docs/parameter_reference.md](docs/parameter_reference.md).

### Linting

```yaml
steps:
  - name: lint
    image: pelotech/drone-helm3
    settings:
      mode: lint
      chart: ./
```

### Installation

```yaml
steps:
  - name: deploy
    image: pelotech/drone-helm3
    settings:
      mode: upgrade
      chart: ./
      release: my-project
    environment:
      KUBE_API_SERVER: https://my.kubernetes.installation/clusters/a-1234
      KUBE_TOKEN:
        from_secret: kubernetes_token
```

### Uninstallation

```yaml
steps:
  - name: uninstall
    image: pelotech/drone-helm3
    settings:
      mode: uninstall
      release: my-project
    environment:
      KUBE_API_SERVER: https://my.kubernetes.installation/clusters/a-1234
      KUBE_TOKEN:
        from_secret: kubernetes_token
```

### Publish

```yaml
steps:
  - name: publish
    image: pelotech/drone-helm3
    settings:
      mode: publish
      chart: ./
    environment:
      REGISTRY_URL: 
        from_secret: registry_url
      REGISTRY_LOGIN_USER_ID: 
        from_secret: registry_login_user_id
      REGISTRY_LOGIN_PASSWORD:
        from_secret: registry_login_password
      REGISTRY_REPO_NAME: <helm_repository_name>
      CHART_VERSION: <version>
```

## Upgrading from drone-helm

drone-helm3 is largely backward-compatible with drone-helm. There are some known differences:

* You'll need to migrate the deployments in the cluster [helm-v2-to-helm-v3](https://helm.sh/blog/migrate-from-helm-v2-to-helm-v3/).
* EKS is not supported. See [#5](https://github.com/pelotech/drone-helm3/issues/5) for more information.
* The `prefix` setting is no longer supported. If you were relying on the `prefix` setting with `secrets: [...]`, you'll need to switch to the `from_secret` syntax.
* During uninstallations, the release history is purged by default. Use `keep_history: true` to return to the old behavior.
* Several settings no longer have any effect. The plugin will produce warnings if any of these are present:
    * `purge` -- this is the default behavior in Helm 3
    * `recreate_pods`
    * `tiller_ns`
    * `upgrade`
    * `canary_image`
    * `client_only`
    * `stable_repo_url`
* Several settings have been renamed, to clarify their purpose and provide a more consistent naming scheme. For backward-compatibility, the old names are still available as aliases. If the old and new names are both present, the updated form takes priority. Conflicting settings will make your `.drone.yml` harder to understand, so we recommend updating to the new names:
    * `helm_command` is now `mode`
    ° `helm_repos` is now `add_repos`
    * `api_server` is now `kube_api_server`
    * `service_account` is now `kube_service_account`
    * `kubernetes_token` is now `kube_token`
    * `kubernetes_certificate` is now `kube_certificate`
    * `wait` is now `wait_for_upgrade`
    * `force` is now `force_upgrade`

Since helm 3 does not require Tiller, we also recommend switching to a service account with less-expansive permissions.

### [Contributing](docs/contributing.md)

This repo is setup in a way that if you enable a personal drone server to build your fork it will
 build and publish your image (makes it easier to test PRs and use the image till the contributions get merged)

* Build local ```DRONE_REPO_OWNER=josmo DRONE_REPO_NAME=drone-ecs drone exec```
* on your server (or cloud.drone.io) just make sure you have DOCKER_USERNAME, DOCKER_PASSWORD, and PLUGIN_REPO set as secrets
