# Logging

```
kubectl -n capi-hivelocity-system logs -f caphv-controller-manager-5cb6548dfb-9nd4k > caphv-controller-manager-5cb6548dfb-9nd4k.log
```

For debugging you can reduce the log output with this script:

```
❯ ./hack/filter-caphv-controller-manager-logs.py caphv-controller-manager-5cb6548dfb-9nd4k.log | less
```
