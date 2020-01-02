# Parameter reference

## Global
| Param name          | Type            | Purpose |
|---------------------|-----------------|---------|
| helm_command        | string          | Indicates the operation to perform. Recommended, but not required. Valid options are `upgrade`, `uninstall`, `lint`, and `help`. |
| update_dependencies | boolean         | Calls `helm dependency update` before running the main command.|
| helm_repos          | list\<string\>  | Calls `helm repo add $repo` before running the main command. Each string should be formatted as `repo_name=https://repo.url/`. |
| namespace           | string          | Kubernetes namespace to use for this operation. |
| debug               | boolean         | Generate debug output within drone-helm3 and pass `--debug` to all helm commands. Use with care, since the debug output may include secrets. |

## Linting

Linting is only triggered when the `helm_command` setting is "lint".

| Param name    | Type           | Required | Purpose |
|---------------|----------------|----------|---------|
| chart         | string         | yes      | The chart to be linted. Must be a local path. |
| values        | list\<string\> |          | Chart values to use as the `--set` argument to `helm lint`. |
| string_values | list\<string\> |          | Chart values to use as the `--set-string` argument to `helm lint`. |
| values_files  | list\<string\> |          | Values to use as `--values` arguments to `helm lint`. |
| lint_strictly | boolean        |          | Pass `--strict` to `helm lint`, to turn warnings into errors. |

## Installation

Installations are triggered when the `helm_command` setting is "upgrade." They can also be triggered when the build was triggered by a `push`, `tag`, `deployment`, `pull_request`, `promote`, or `rollback` Drone event.

| Param name             | Type           | Required | Purpose |
|------------------------|----------------|----------|---------|
| chart                  | string         | yes      | The chart to use for this installation. |
| release                | string         | yes      | The release name for helm to use. |
| api_server             | string         | yes      | API endpoint for the Kubernetes cluster. |
| kubernetes_token       | string         | yes      | Token for authenticating to Kubernetes. |
| service_account        | string         |          | Service account for authenticating to Kubernetes. Default is `helm`. |
| kubernetes_certificate | string         |          | Base64 encoded TLS certificate used by the Kubernetes cluster's certificate authority. |
| chart_version          | string         |          | Specific chart version to install. |
| dry_run                | boolean        |          | Pass `--dry-run` to `helm upgrade`. |
| wait                   | boolean        |          | Wait until kubernetes resources are in a ready state before marking the installation successful. |
| timeout                | duration       |          | Timeout for any *individual* Kubernetes operation. The installation's full runtime may exceed this duration. |
| force                  | boolean        |          | Pass `--force` to `helm upgrade`. |
| values                 | list\<string\> |          | Chart values to use as the `--set` argument to `helm upgrade`. |
| string_values          | list\<string\> |          | Chart values to use as the `--set-string` argument to `helm upgrade`. |
| values_files           | list\<string\> |          | Values to use as `--values` arguments to `helm upgrade`. |
| reuse_values           | boolean        |          | Reuse the values from a previous release. |
| skip_tls_verify        | boolean        |          | Connect to the Kubernetes cluster without checking for a valid TLS certificate. Not recommended in production. |

## Uninstallation

Uninstallations are triggered when the `helm_command` setting is "uninstall" or "delete." They can also be triggered when the build was triggered by a `delete` Drone event.

| Param name             | Type     | Required | Purpose |
|------------------------|----------|----------|---------|
| release                | string   | yes      | The release name for helm to use. |
| api_server             | string   | yes      | API endpoint for the Kubernetes cluster. |
| kubernetes_token       | string   | yes      | Token for authenticating to Kubernetes. |
| service_account        | string   |          | Service account for authenticating to Kubernetes. Default is `helm`. |
| kubernetes_certificate | string   |          | Base64 encoded TLS certificate used by the Kubernetes cluster's certificate authority. |
| dry_run                | boolean  |          | Pass `--dry-run` to `helm uninstall`. |
| timeout                | duration |          | Timeout for any *individual* Kubernetes operation. The uninstallation's full runtime may exceed this duration. |
| skip_tls_verify        | boolean  |          | Connect to the Kubernetes cluster without checking for a valid TLS certificate. Not recommended in production. |
| chart                  | string   |          | Required when the global `update_dependencies` parameter is true. No effect otherwise. |

### Where to put settings

Any setting can go in either the `settings` or `environment` section. If a setting exists in _both_ sections, the version in `environment` will override the version in `settings`.

We recommend putting all drone-helm3 configuration in the `settings` block and limiting the `environment` block to variables that are used when building your charts.

### Formatting non-string values

* Booleans can be yaml's `true` and `false` literals or the strings `"true"` and `"false"`.
* Durations are strings formatted with the syntax accepted by [golang's ParseDuration function](https://golang.org/pkg/time/#ParseDuration) (e.g. 5m30s)
  * For backward-compatibility with drone-helm, a duration can also be an integer, in which case it will be interpreted to mean seconds.
* List\<string\>s can be a yaml sequence or a comma-separated string.

All of the following are equivalent:

```yaml
values: "foo=1,bar=2"
values: ["foo=1", "bar=2"]
values:
  - foo=1
  - bar=2
```

Note that **list members must not contain commas**. Both of the following are equivalent:

```yaml
values_files: [ "./over_9,000.yml" ]
values_files: [ "./over_9", "000.yml" ]
```
