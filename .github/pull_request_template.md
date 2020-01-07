**Please replace this line with "fixes #ISSUE_NUMBER" (or "relates to #ISSUE_NUMBER", if it is not a complete fix)**

Pre-merge checklist:

* [ ] Code changes have tests
* [ ] Any config changes are documented:
    * If the change touches _required_ config, there's a corresponding update to `README.md`
    * There's a corresponding update to `docs/parameter_reference.md`
    * There's a pull request to update [the parameter reference in drone-plugin-index](https://github.com/drone/drone-plugin-index/blob/master/content/pelotech/drone-helm3/index.md)
* [ ] Any large changes have been verified by running a Drone job
