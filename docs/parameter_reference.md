# Parameter reference

## Global
| Param name          | Type            | Alias        | Purpose |
|---------------------|-----------------|--------------|---------|
| mode                | string          | helm_command | Indicates the operation to perform. Recommended, but not required. Valid options are `upgrade`, `uninstall`, `lint`, and `help`. |
| update_dependencies | boolean         |              | Calls `helm dependency update` before running the main command.|
| add_repos           | list\<string\>  | helm_repos   | Calls `helm repo add $repo` before running the main command. Each string should be formatted as `repo_name=https://repo.url/`. |
| repo_certificate    | string          |              | Base64 encoded TLS certificate for a chart repository. |
| repo_ca_certificate | string          |              | Base64 encoded TLS certificate for a chart repository certificate authority. |
| namespace           | string          |              | Kubernetes namespace to use for this operation. |
| debug               | boolean         |              | Generate debug output within drone-helm3 and pass `--debug` to all helm commands. Use with care, since the debug output may include secrets. |

## Linting

Linting is only triggered when the `mode` setting is "lint".

| Param name    | Type           | Required | Purpose |
|---------------|----------------|----------|---------|
| chart         | string         | yes      | The chart to be linted. Must be a local path. |
| values        | list\<string\> |          | Chart values to use as the `--set` argument to `helm lint`. |
| string_values | list\<string\> |          | Chart values to use as the `--set-string` argument to `helm lint`. |
| values_files  | list\<string\> |          | Values to use as `--values` arguments to `helm lint`. |
| lint_strictly | boolean        |          | Pass `--strict` to `helm lint`, to turn warnings into errors. |

## Installation

Installations are triggered when the `mode` setting is "upgrade." They can also be triggered when the build was triggered by a `push`, `tag`, `deployment`, `pull_request`, `promote`, or `rollback` Drone event.

| Param name             | Type           | Required | Alias                  | Purpose |
|------------------------|----------------|----------|------------------------|---------|
| chart                  | string         | yes      |                        | The chart to use for this installation. |
| release                | string         | yes      |                        | The release name for helm to use. |
| skip_kubeconfig        | boolean        |          |                        | Whether to skip kubeconfig file creation. |
| kube_api_server        | string         | yes      | api_server             | API endpoint for the Kubernetes cluster. This is ignored if `skip_kubeconfig` is `true`. |
| kube_token             | string         | yes      | kubernetes_token       | Token for authenticating to Kubernetes. This is ignored if `skip_kubeconfig` is `true`. |
| kube_service_account   | string         |          | service_account        | Service account for authenticating to Kubernetes. Default is `helm`. This is ignored if `skip_kubeconfig` is `true`. |
| kube_certificate       | string         |          | kubernetes_certificate | Base64 encoded TLS certificate used by the Kubernetes cluster's certificate authority. This is ignored if `skip_kubeconfig` is `true`. |
| chart_version          | string         |          |                        | Specific chart version to install. |
| dry_run                | boolean        |          |                        | Pass `--dry-run` to `helm upgrade`. |
| dependencies_action    | string         |          |                        | Calls `helm dependency build` OR `helm dependency update` before running the main command. Possible values: `build`, `update`. |
| wait_for_upgrade       | boolean        |          | wait                   | Wait until kubernetes resources are in a ready state before marking the installation successful. |
| timeout                | duration       |          |                        | Timeout for any *individual* Kubernetes operation. The installation's full runtime may exceed this duration. |
| force_upgrade          | boolean        |          | force                  | Pass `--force` to `helm upgrade`. |
| atomic_upgrade         | boolean        |          |                        | Pass `--atomic` to `helm upgrade`. |
| cleanup_failed_upgrade | boolean        |          |                        | Pass `--cleanup-on-fail` to `helm upgrade`. |
| history_max            | int            |          |                        | Pass `--history-max` to `helm upgrade`. |
| values                 | list\<string\> |          |                        | Chart values to use as the `--set` argument to `helm upgrade`. |
| string_values          | list\<string\> |          |                        | Chart values to use as the `--set-string` argument to `helm upgrade`. |
| values_files           | list\<string\> |          |                        | Values to use as `--values` arguments to `helm upgrade`. |
| reuse_values           | boolean        |          |                        | Reuse the values from a previous release. |
| skip_tls_verify        | boolean        |          |                        | Connect to the Kubernetes cluster without checking for a valid TLS certificate. Not recommended in production. This is ignored if `skip_kubeconfig` is `true`. |
| create_namespace       | boolean        |          |                        | Pass --create-namespace to `helm upgrade`. |
| skip_crds              | boolean        |          |                        | Pass --skip-crds to `helm upgrade`. |

## Uninstallation

Uninstallations are triggered when the `mode` setting is "uninstall" or "delete." They can also be triggered when the build was triggered by a `delete` Drone event.

| Param name             | Type     | Required | Alias                  | Purpose |
|------------------------|----------|----------|------------------------|---------|
| release                | string   | yes      |                        | The release name for helm to use. |
| skip_kubeconfig        | boolean  |          |                        | Whether to skip kubeconfig file creation. |
| kube_api_server        | string   | yes      | api_server             | API endpoint for the Kubernetes cluster. This is ignored if `skip_kubeconfig` is `true`. |
| kube_token             | string   | yes      | kubernetes_token       | Token for authenticating to Kubernetes. This is ignored if `skip_kubeconfig` is `true`. |
| kube_service_account   | string   |          | service_account        | Service account for authenticating to Kubernetes. Default is `helm`. This is ignored if `skip_kubeconfig` is `true`. |
| kube_certificate       | string   |          | kubernetes_certificate | Base64 encoded TLS certificate used by the Kubernetes cluster's certificate authority. This is ignored if `skip_kubeconfig` is `true`. |
| keep_history           | boolean  |          |                        | Pass `--keep-history` to `helm uninstall`, to retain the release history. |
| dry_run                | boolean  |          |                        | Pass `--dry-run` to `helm uninstall`. |
| timeout                | duration |          |                        | Timeout for any *individual* Kubernetes operation. The uninstallation's full runtime may exceed this duration. |
| skip_tls_verify        | boolean  |          |                        | Connect to the Kubernetes cluster without checking for a valid TLS certificate. Not recommended in production. This is ignored if `skip_kubeconfig` is `true`. |
| chart                  | string   |          |                        | Required when the global `update_dependencies` parameter is true. No effect otherwise. |

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

### Interpolating secrets into the `values`, `string_values` and `add_repos` settings

If you want to send secrets to your charts, you can use syntax similar to shell variable interpolation--either `$VARNAME` or `$${VARNAME}`. The double dollar-sign is necessary when using curly brackets; using curly brackets with a single dollar-sign will trigger Drone's string substitution (which can't use arbitrary environment variables). If an environment variable is not set, it will be treated as if it were set to the empty string.

```yaml
environment:
  DB_PASSWORD:
    from_secret: db_password
  SESSION_KEY:
    from_secret: session_key
settings:
  values:
    - db_password=$DB_PASSWORD    # db_password will be set to the contents of the db_password secret
    - db_pass=$DB_PASS            # db_pass will be set to "" since $DB_PASS is not set
    - session_key=$${SESSION_KEY} # session_key will be set to the contents of the session_key secret
    - sess_key=${SESSION_KEY}     # sess_key will be set to "" by Drone's variable substitution
```

Variables intended for interpolation must be set in the `environment` section, not `settings`.

### Backward-compatibility aliases

Some settings have alternate names, for backward-compatibility with drone-helm. We recommend using the canonical name unless you require the backward-compatible form.

| Canonical name       | Alias |
|----------------------|-------|
| mode                 | helm_command |
| add_repos            | helm_repos |
| kube_api_server      | api_server |
| kube_service_account | service_account |
| kube_token           | kubernetes_token |
| kube_certificate     | kubernetes_certificate |
| wait_for_upgrade     | wait |
| force_upgrade        | force |
