# Drone plugin for Helm 3

This plugin provides an interface between [Drone](https://drone.io/) and [Helm 3](https://github.com/kubernetes/helm):

* Lint your charts
* Deploy your service
* Delete your service

The plugin is inpsired by [drone-helm](https://github.com/ipedrazas/drone-helm), which fills the same role for Helm 2. It provides a comparable feature-set and the configuration settings are backwards-compatible.

## Example configuration

These examples give a minimal and sufficient configuration for each use-case. For a full description of each case, see [docs/lint_example.yml](docs/lint_example.yml), [docs/upgrade_example.yml](docs/upgrade_example.yml), and [docs/delete_example.yml](docs/delete_example.yml).

### Linting

```yaml
steps:
  - name: lint
    image: pelotech/drone-helm3
    settings:
      helm_command: lint
      chart: ./
```

### Deployment

```yaml
steps:
  - name: deploy
    image: pelotech/drone-helm3
    settings:
      helm_command: upgrade
      chart: ./
      release: my-project
    environment:
      API_SERVER: https://my.kubernetes.installation/clusters/a-1234
      KUBERNETES_TOKEN:
        from_secret: kubernetes_token
```

### Deletion

```yaml
steps:
  - name: delete
    image: pelotech/drone-helm3
    settings:
      helm_command: delete
      release: my-project
    environment:
      API_SERVER: https://my.kubernetes.installation/clusters/a-1234
      KUBERNETES_TOKEN:
        from_secret: kubernetes_token
```

## Upgrading from drone-helm

The setting names for drone-helm3 are backwards-compatible with those for drone-helm, so the only mandatory step is to update the `image` clause so that drone uses the new plugin.

There are several settings that no longer have any effect:

* `purge` -- this is the default behavior in Helm 3
* `recreate_pods`
* `tiller_ns`
* `upgrade`
* `canary_image`
* `client_only`
* `stable_repo_url`

If your `.drone.yml` contains those settings, we recommend removing them.
