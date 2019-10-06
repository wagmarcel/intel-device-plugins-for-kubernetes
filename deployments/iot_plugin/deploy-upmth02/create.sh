#NAMESPACE=oisp-devices
NODENAME=${1:-NODENAME}
NAMESPACE=${NAMESPACE:-default}
CONFIG_MAP_NAME=upmth02-config

echo Creating subdir for node $NODENAME
mkdir -p $NODENAME
echo Copy and adapt the templates
for file in $(ls *.yaml *.json); do 
	echo Processing $file
	sed 's|<NODENAME>|'$NODENAME'|g' $file > $NODENAME/$file
done
#kubectl create namespace ${NAMESPACE} 2>/dev/null || echo "Namespace $NAMESPACE already exists. Continue."
#kubectl delete configmap ${CONFIG_MAP_NAME} -n ${NAMESPACE} 2>/dev/null || echo "ConfigMap not existing. Continue."
kubectl create configmap ${CONFIG_MAP_NAME}  --from-file=./$NODENAME/sensorSpecs.json -o yaml --dry-run > $NODENAME/config.yaml 
cat $NODENAME/*.yaml > $NODENAME/all.yaml
