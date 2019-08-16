NAMESPACE=oisp-device
SECRET_NAME=oisp-devices-secret
kubectl create namespace ${NAMESPACE} 2>/dev/null || echo "Namespace already exists. Continue."
kubectl delete secret ${SECRET_NAME} -n ${NAMESPACE} 2>/dev/null || echo "Secret not existing. Continue."
kubectl create secret generic ${SECRET_NAME} --from-file=./activationCode.txt -n ${NAMESPACE} 
