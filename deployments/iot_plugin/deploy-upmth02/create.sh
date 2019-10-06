#NAMESPACE=oisp-devices
NODENAME=${1:-NODENAME}
NAMESPACE=${NAMESPACE:-default}
CONFIG_MAP_NAME=oisp-devices-config

echo Creating subdir for node $NODENAME
mkdir -p $NODENAME
echo Copy and adapt the templates
for yamlfile in $(ls *.yaml); do 
	echo Processing $yamlfile
	sed 's|<NODENAME>|'$NODENAME'|g' $yamlfile > $NODENAME/$yamlfile
done
#kubectl create namespace ${NAMESPACE} 2>/dev/null || echo "Namespace $NAMESPACE already exists. Continue."
#kubectl delete configmap ${CONFIG_MAP_NAME} -n ${NAMESPACE} 2>/dev/null || echo "ConfigMap not existing. Continue."
#kubectl create configmap ${CONFIG_MAP_NAME} --from-file=./config.json.r1 --from-file=./config.json.r2 --from-file=./sensorSpecs.json -n ${NAMESPACE}
