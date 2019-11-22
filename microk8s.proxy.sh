#!/bin/bash
exec microk8s.kubectl proxy --accept-hosts=.* --address=0.0.0.0 &
echo http://127.0.0.1:8001/api/v1/namespaces/kube-system/services/https:kubernetes-dashboard:/proxy/#/overview?namespace=default
