# Helm setup

## Installing Smallstep CA
```shell
helm install step-certificates smallstep/step-certificates \
  --namespace step-ca --create-namespace
```

After installing
```shell
NAME: step-certificates
LAST DEPLOYED: <date_time>
NAMESPACE: step-ca
STATUS: deployed
REVISION: 1
NOTES:
Thanks for installing Step CA.

1. Get the PKI and Provisioner secrets running these commands:
   kubectl get -n step-ca -o jsonpath='{.data.password}' secret/step-certificates-ca-password | base64 --decode
   kubectl get -n step-ca -o jsonpath='{.data.password}' secret/step-certificates-provisioner-password | base64 --decode
2. Get the CA URL and the root certificate fingerprint running this command:
   kubectl -n step-ca logs job.batch/step-certificates

3. Delete the configuration job running this command:
   kubectl -n step-ca delete job.batch/step-certificates
```