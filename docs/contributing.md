# Contributing to drone-helm3

We're glad you're interested in contributing! Here are some guidelines that will help make sure everyone has a good experience:

## Submitting a patch

Before you start working on a change, please make sure there's an associated issue. It doesn't need to be thoroughly scrutinized and dissected, but it needs to exist.

Please put the relevant issue number in the first line of your commit messages, e.g. `vorpalize the frabjulator [#42]`. Branch names do not need issue numbers, but feel free to include them if you like.

We encourage you to follow [the guidelines in Pro Git](https://git-scm.com/book/en/v2/Distributed-Git-Contributing-to-a-Project#_commit_guidelines) when making commits. In short:

* Commit early and commit often.
* Make the first line of the commit message concise--no more than 50 characters or so.
* Make the rest of the commit message verbose--information about _why_ you did what you did is particularly helpful.

Once you're satisfied with your work, send us a pull request. If you'd like, you can send the pull request _before_ you're satisfied with your work; just be sure to mark the PR a draft or put `[WIP]` in the title.

## How to run the tests

We use `go test`, `go vet`, and `golint`:

```
go test ./cmd/... ./internal/...
go vet ./cmd/... ./internal/...
golint -set_exit_status ./cmd/... ./internal/...
```

If you have [the Drone cli tool](https://docs.drone.io/cli/install/) installed, you can also use `drone exec --include test --include lint`.

## Testing the plugin end-to-end

Although we aim to make the internal tests as thorough as possible, they can't test drone-helm3's integration with drone and helm themselves. However, you can test a change manually by building an image and running it with a fixture repository.

You will need:

* Access to a docker image registry. This document assumes you'll use [Docker Hub](https://hub.docker.com).
* [The Drone cli tool](https://docs.drone.io/cli/install/).
* A fixture repository--a directory with a `.drone.yml` and a helm chart. If you don't have one handy, try adding a `.drone.yml` to a chart from [Helm's "stable" repository](https://github.com/helm/charts/tree/master/stable/).
* Access to a kubernetes cluster (unless `lint` or `dry_run` is sufficient for your purposes).

Once you have what you need, you can publish and consume an image with your changes:

1. [Create a repository on Docker Hub](https://hub.docker.com/repository/create). This document assumes you've called it drone-helm3-testing.
1. Create a `.secrets` file with your docker credentials (see [example.secrets](./example.secrets) for an example). While you can use your Docker Hub password, it's better to [generate an access token](https://hub.docker.com/settings/security) and use that instead.
1. Use Drone to build and publish an image with your changes: `drone exec --secret-file ./secrets --event push`
1. In the `.drone.yml` of your fixture repository, set the `image` for each relevant stanza to `your_dockerhub_username/drone-helm3-testing`
1. Use `drone exec` in the fixture repo to verify your changes.
