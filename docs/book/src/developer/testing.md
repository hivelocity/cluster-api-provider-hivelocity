# Testing

We use (from simple to complex):

* Go unit tests
* [envtest](https://github.com/kubernetes-sigs/controller-runtime/tree/main/pkg/envtest) from controller-runtime
* [End-to-End Tests](e2e.md)

The unit tests and envtests get executed via `make test-unit`.

Please add new tests, if you add new features. 

Try to use a simple form. For example prefer a unit test to a test which needs envtest.