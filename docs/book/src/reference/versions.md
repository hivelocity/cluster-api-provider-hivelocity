# Version Support

## Release Versioning

CAPHV follows the [semantic versionining][semver] specification:

MAJOR version release for incompatible API changes, 
MINOR version release for backwards compatible feature additions, 
and PATCH version release for only bug fixes.

**Example versions:**

- Minor release: `v0.1.0`
- Patch release: `v0.1.1`
- Major release: `v1.0.0`


## Compatibility with Cluster API Versions

CAPHV's versions are compatible with the following versions of Cluster API

|                         | Cluster API `v1beta1` (`v1.3.x`) |
|-------------------------|----------------------------------|
| CAPHV v1alpha1 `v1.0.x` | ✓                                |


CAPHV versions are not released in lock-step with Cluster API releases.
Multiple CAPHV minor releases can use the same Cluster API minor release.

For compatibility, check the release notes [here](https://github.com/hivelocity/cluster-api-provider-hivelocity/releases/) to see which v1beta1 Cluster API version each CAPHV version is compatible with.

For example:
- CAPHV v1.0.x is compatible with Cluster API v1.3.x

## End-of-Life Timeline
TBD

## Compatibility with Kubernetes Versions

 CAPHV API versions support all Kubernetes versions that is supported by its compatible Cluster API version:

|     API Versions             | CAPI v1beta1 (v1.x) |
| ---------------------------- | -------------- |
| CAPHV v1alpha1 (v0.1)        | ✓             |


(See [Kubernetes support matrix][cluster-api-supported-v] of Cluster API versions).

Tested Kubernetes Versions:

|                   | Hivelocity Provider `v0.1.x` |
|-------------------|------------------------------|
| Kubernetes 1.25.x | ✓                            |


[cluster-api-supported-v]: https://cluster-api.sigs.k8s.io/reference/versions.html
[semver]: https://semver.org/#semantic-versioning-200
