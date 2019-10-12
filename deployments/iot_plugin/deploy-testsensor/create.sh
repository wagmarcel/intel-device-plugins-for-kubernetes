#NAMESPACE=oisp-devices
NODENAME=${1:-NODENAME}
NAMESPACE=${NAMESPACE:-default}
CONFIG_MAP_NAME=testsensor-config
if [ ! -z "$2" ]; then
	ID=-$2
fi

echo Creating subdir for node $NODENAME$ID
mkdir -p $NODENAME$ID
echo Copy and adapt the templates
for file in $(ls *.yaml *.json); do 
	echo Processing $file
	sed 's|<NODENAME>|'$NODENAME'|g' $file | sed 's|<ID>|'$ID'|g' > $NODENAME$ID/$file
done
#kubectl create namespace ${NAMESPACE} 2>/dev/null || echo "Namespace $NAMESPACE already exists. Continue."
#kubectl delete configmap ${CONFIG_MAP_NAME} -n ${NAMESPACE} 2>/dev/null || echo "ConfigMap not existing. Continue."
kubectl create configmap ${CONFIG_MAP_NAME}  --from-file=./$NODENAME$ID/sensorSpecs.json -o yaml --dry-run > $NODENAME$ID/configmap.yaml 
cat $NODENAME$ID/configmap.yaml > $NODENAME$ID/all.yaml
echo --- >> $NODENAME$ID/all.yaml
cat $NODENAME$ID/pvc.yaml >> $NODENAME$ID/all.yaml
echo --- >> $NODENAME$ID/all.yaml
cat $NODENAME$ID/deployment.yaml >> $NODENAME$ID/all.yaml
