# Kubernetes Cluster API Provider Hivelocity
> :warning: This project is in the development stage. DO NOT USE IN PRODUCTION! :warning:

The Hivelocity Provider is a Kubernetes-native tool that allows you to easily create and manage declarative infrastructure for your Kubernetes clusters on Hivelocity's infrastructure. It offers options for high availability on instant bare metal or custom dedicated setups and simplifies the process of creating, updating, and operating production-ready clusters.

You can find more information about Hivelocity and their infrastructure at https://www.hivelocity.net/.


> If you have questions or are interested in running production-ready Kubernetes clusters on Hivelocity, then please contact us via e-mail: [info@syself.com](mailto:info@syself.com?subject=cluster-api-provider-hivelocity).

## :newspaper: What is the Cluster API Provider Hivelocity?
The Cluster API is an operator that manages infrastructure in a similar way to how Kubernetes manages containers. It uses a declarative API and includes controllers that ensure the desired state of the infrastructure is maintained. This approach, called Infrastructure as Software, allows for more automatic reactions to changes and problems compared to Infrastructure as Code solutions.

The Hivelocity Provider is the infrastructure component of the Cluster API stack that allows the Cluster API to be used on Hivelocity's infrastructure. It enables the creation of stable and highly available Kubernetes clusters, allowing organizations to benefit from the advantages of declarative infrastructure and cost-effectiveness on a global scale. The Hivelocity Provider allows for the creation of stable and highly available Kubernetes clusters on certified HIPAA, PCI, ISAE-3402, SSAE 16 SOC1, and SOC2 infrastructure around the globe.

With the Hivelocity Provider, you can trust that your infrastructure is in good hands with a provider that has a track record of dynamic performance, static pricing, and a global presence.

---
## :book: Documentation

Please see our [book](https://hivelocity.github.io/cluster-api-provider-hivelocity) for in-depth documentation.

## :sparkles: Features

* Native Kubernetes manifests and API
* Manages the bootstrapping of Networking, Loadbalancers and devices.
* Choice of Linux distribution
* Support for single and multi-node control plane clusters (HA Kubernetes)
* Doesn't use SSH for bootstrapping nodes.
* Day 2 operations including: updating Kubernetes and nodes, scaling, and self-healing
* Custom CSR approver for approving kubelet-serving certificate signing requests
* Support for both Hivelocity instant bare metal and custom dedicated setups


## :rocket: Get Started

If you're looking to jump straight into it, check out the [Quick Start Guide](https://hivelocity.github.io/cluster-api-provider-hivelocity/user/getting-started.html)

## :fire: Compatibility with Cluster API and Kubernetes Versions

Please see: https://hivelocity.github.io/cluster-api-provider-hivelocity/reference/versions.html

**NOTE:** As the versioning for this project is tied to the versioning of Cluster API, future modifications to this policy may be made to more closely align with other providers in the Cluster API ecosystem.

---

## :busts_in_silhouette: Getting Involved and Contributing

Are you interested in contributing to Cluster API Provider Hivelocity? We, the
maintainers and community, would love your suggestions, contributions, and help!
If you want to learn more about how to get involved, you can contact the maintainers at any time.

To set up your environment, try out the development guide.

In the interest of getting more new people involved, we tag issues with
[`good first issue`][good_first_issue].
These are typically issues that have a smaller scope, but are good to get acquainted with the codebase.

We also encourage ALL active community participants to act as if they are
maintainers, even if you don't have "official" write permissions. This is a
community effort, we are here to serve the Kubernetes community. If you have an
active interest and you want to get involved, you have real power! Don't assume
that the only people who can get things done around here are the "maintainers".

We would also love to add more "official" maintainers, so show us what you can
do!

## :dizzy: Code of Conduct

Participation in the Kubernetes community is governed by the [Kubernetes Code of Conduct](https://github.com/hivelocity/cluster-api-provider-hivelocity/blob/main/code-of-conduct.md).

## :construction: Github Issues

### :bug: Bugs

If you think you have found a bug, please follow these steps:

- Take some time to give due diligence to the issue tracker. Your issue might be a duplicate.
- Get the logs from the cluster controllers. Paste this into your issue.
- Open a [bug report][bug_report].
- Give it a meaningful title to help others who might be searching for your issue in the future.
- If you have questions, reach out to the Cluster API community on the [Kubernetes Slack channel][slack_info].

### :star: Tracking New Features

We also use the issue tracker to track features. If you have an idea for a feature or think that you can help Cluster API Provider Hivelocity become even more awesome, then follow these steps:

- Open a [feature request][feature_request].
- Give it a meaningful title to help others who might be searching for your issue in the future.
- Define clearly the use case. Use concrete examples, e.g. "I type `this` and
  Cluster API Provider Hivelocity does `that`".
- Some of our larger features will require some design. If you would like to
  include a technical design for your feature, please include it in the issue.
- After the new feature is well understood and the design is agreed upon, we can
  start coding the feature. We would love if you code it. So please open
  up a **WIP** *(work in progress)* pull request. Happy coding!

<!-- References -->

[good_first_issue]: https://github.com/hivelocity/cluster-api-provider-hivelocity/issues?q=is%3Aissue+is%3Aopen+sort%3Aupdated-desc+label%3A%22good+first+issue%22
[bug_report]: https://github.com/hivelocity/cluster-api-provider-hivelocity/issues/new?template=bug_report.md
[feature_request]: https://github.com/hivelocity/cluster-api-provider-hivelocity/issues/new?template=feature_request.md
[slack_info]: https://github.com/kubernetes/community/tree/master/communication#slack
[cluster_api]: https://github.com/kubernetes-sigs/cluster-api
[quickstart]: https://cluster-api.sigs.k8s.io/user/quick-start.html