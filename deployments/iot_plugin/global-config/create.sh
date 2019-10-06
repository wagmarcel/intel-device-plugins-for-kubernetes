#NAMESPACE=oisp-devices
NAMESPACE=${NAMESPACE:-default}
CONFIG_MAP_NAME=global-devices-config
SECRET_NAME=global-devices-secret
kubectl create namespace ${NAMESPACE} 2>/dev/null || echo "Namespace $NAMESPACE already exists. Continue."
kubectl delete configmap ${CONFIG_MAP_NAME} -n ${NAMESPACE} 2>/dev/null || echo "ConfigMap not existing. Continue."
kubectl create configmap ${CONFIG_MAP_NAME} --from-file=./config.json.r1 --from-file=./config.json.r2 -n ${NAMESPACE}
kubectl delete secret ${SECRET_NAME} -n ${NAMESPACE} 2>/dev/null || echo "Secret not existing. Continue."
kubectl create secret generic ${SECRET_NAME} --from-file=./activationCode --from-file=./activationCode-r2 -n ${NAMESPACE}

