NAMESPACE=oisp-devices
CONFIG_MAP_NAME=oisp-devices-config
kubectl create namespace ${NAMESPACE} 2>/dev/null || echo "Namespace already exists. Continue."
kubectl delete configmap ${CONFIG_MAP_NAME} -n ${NAMESPACE} 2>/dev/null || echo "ConfigMap not existing. Continue."
kubectl create configmap ${CONFIG_MAP_NAME} --from-file=./config.json -n ${NAMESPACE} 
