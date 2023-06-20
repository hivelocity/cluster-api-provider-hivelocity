# CI

We use Github Actions to:

* lint the code, so that it follows best practices
* lint yaml files and markdown (documentation)
* run unit tests
* run e2e tests

## Latest Actions

You can see the result of the latest jobs here: [github.com/hivelocity/cluster-api-provider-hivelocity/actions](https://github.com/hivelocity/cluster-api-provider-hivelocity/actions)

## Running Github Actions locally

If you are updating the Github Actions, you can speed up your edit-test feedback loop 
if you run the actions locally.

You need to install [nektos/act](https://github.com/nektos/act). 

Then you can have a look at the available jobs:

```
act -l

Stage  Job ID               Job name                      Workflow name            Workflow file         Events             
0      markdown-link-check  Broken Links                  Check PR Markdown links  lint-docs-pr.yml      pull_request       
0      deploy               deploy                        mdBook github pages      mdbook.yml            push               
0      test-release         Create a Test Release         E2E PR Blocking          pr-e2e.yaml           pull_request_target
0      manager-image        Build and push manager image  E2E PR Blocking          pr-e2e.yaml           pull_request_target
0      golangci             lint                          golangci-lint            pr-golangci-lint.yml  pull_request       
0      shellcheck           shellcheck-lint               shellcheck               pr-shellcheck.yml     pull_request       
0      starlark             run starlark lint             starlark                 pr-starlark-lint.yml  pull_request       
0      unit-test            unit-test                     unit-test                pr-unit-test.yml      pull_request       
0      verify-code          Verify Code                   Verify PR Code           pr-verify-code.yml    pull_request       
0      verify-pr-content    verify PR contents            Verify Pull Request      pr-verify.yml         pull_request_target
0      yamllint             yamllint                      yamllint test            pr-yamllint.yml       pull_request       
1      e2e-basic            End-to-End Test Basic         E2E PR Blocking          pr-e2e.yaml           pull_request_target
```

Then you can call single jobs:

```
act -j yamllint 

[yamllint test/yamllint] 🚀  Start image=catthehacker/ubuntu:act-latest
[yamllint test/yamllint]   🐳  docker pull image=catthehacker/ubuntu:act-latest platform= username= forcePull=true
[yamllint test/yamllint] using DockerAuthConfig authentication for docker pull
[yamllint test/yamllint]   🐳  docker create image=catthehacker/ubuntu:act-latest platform= entrypoint=["tail" "-f" "/dev/null"] cmd=[]
[yamllint test/yamllint]   🐳  docker run image=catthehacker/ubuntu:act-latest platform= entrypoint=["tail" "-f" "/dev/null"] cmd=[]
[yamllint test/yamllint]   ☁  git clone 'https://github.com/actions/setup-python' # ref=v4
[yamllint test/yamllint] ⭐ Run Main actions/checkout@v3
...
```

Actions which are known to work with `act`:
* golangci
* shellcheck
* starlark
* unit-test
* verify-code
* verify-pr-content
* yamllint

Actions which are known to not work with `act`:
* markdown-link-check

Using `act` is optional and not needed, except you want to update the Github Actions. All actions are
available as Makefile targets (see `make help`), too.

## Updating Github Actions

In CI the Github Actions of the main branch get executed.

If you want to update the Github Actions you need to merge the changes to the main branch first.

See docs for [pull_request_target](https://docs.github.com/en/actions/using-workflows/events-that-trigger-workflows#pull_request_target)