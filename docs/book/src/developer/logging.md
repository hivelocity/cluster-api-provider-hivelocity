# Logging

The interesting logs are from the caphv-controller-manager, because these logs got created by the code of this repository.

You can get the logs via `kubectl`

First you need the exact name of the controller:

```
❯ kubectl get pods -A | grep caphv-controller

capi-hivelocity-system              caphv-controller-manager-7889d9d768-7m8hr                        2/2     Running   0          53m
```

Then you can fetch the logs:
```
❯ kubectl -n capi-hivelocity-system logs -f caphv-controller-manager-7889d9d768-7m8hr   > caphv-controller-manager-7889d9d768-7m8hr.log
```

You can see these logs via Tilt, too.

For debugging you can reduce the log output with this script:

```
❯ ./hack/filter-caphv-controller-manager-logs.py caphv-controller-manager-7889d9d768-7m8hr.log | tail
```

Or:

```
❯ kubectl -n capi-hivelocity-system logs \
    $(kubectl -n capi-hivelocity-system get pods | grep caphv-controller-manager | cut -d' ' -f1) \
    | ./hack/filter-caphv-controller-manager-logs.py - | tail
```

The command [make watch](../topics/make-watch.md) shows the last lines of logs, too.