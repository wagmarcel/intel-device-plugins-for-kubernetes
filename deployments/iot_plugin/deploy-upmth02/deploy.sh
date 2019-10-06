#NAMESPACE=oisp-devices
NODENAME=${1:-NODENAME}
NAMESPACE=${NAMESPACE:-default}
UPMTH02_DEVICESPEC=${UPMTH02_DEVICESPEC:-'[{"id":"1234567890","name":"upm-th02","hostPath":"/dev/i2c-5","containerPath":"/dev/i2c-5","permission":"rw","gid":1001}]'}
#DRYRUN="--dry-run"

echo device spec: ${UPMTH02_DEVICESPEC}
kubectl label node --overwrite $NODENAME deviceType=iot $DRYRUN || echo "Label already exists. Continue."
kubectl create namespace ${NAMESPACE} $DRYRUN 2>/dev/null || echo "Namespace $NAMESPACE already exists. Continue."
kubectl annotate --overwrite node $NODENAME oisp.net/deviceSpec=$UPMTH02_DEVICESPEC $DRYRUN
kubectl apply -f $NODENAME/all.yaml -n $NAMESPACE $DRYRUN
