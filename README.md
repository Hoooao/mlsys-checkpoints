## Checkpoint 1

**Dockerfile:** `dockerfile`

**K8S configuration file:** `fib-deployment.yaml` (app's config file and that of istio are put together; the first segment would be that of app.)

**Sample trace** `traces.txt` (gained from microk8s and ran 1000 commands)



### Performance: run 10000 commands/computations in total

- No Sampling:
  1. CPU: ~2.92% 
  2. Memory(before exit): 61.5MB

- 50% Sampling:
  1. CPU: ~3.9%
  2. Memory(before exit): 66.5MB

- Always Sampling:
  1. CPU: ~4.89%
  2. Memory(before exit): 68.5MB



