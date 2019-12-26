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

## Using the plugin locally

The internal tests can't test drone-helm3's integration with drone and helm themselves. However, you can build and run a local image to test a change end-to-end.

You will need:

* A Docker image registry. See [docs.docker.com/registry/](https://docs.docker.com/registry/) for information on standing up a local registry.
* You will also need [the Drone cli tool](https://docs.drone.io/cli/install/).
* A fixture repo--a repo with a `.drone.yml` and a helm chart.
* Access to a kubernetes cluster.

Once you have a local registry, uncomment the `publish_locally` step in `.drone.yml` and replace the `0.0.0.0`s with your computer's local IP address.

Now you can run `drone exec --include test --include lint --include publish_locally` to build an image and publish it to your local registry.

Finally, configure your fixture repo to use the locally-published image, e.g. `image: 192.168.0.1:5000/drone-helm3`.

Now you can use `drone exec` in the fixture repo to verify your changes.
