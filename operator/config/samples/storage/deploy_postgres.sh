#!/bin/bash
KUBECONFIG=${1:-$KUBECONFIG}
currentDir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

TARGET_NAMESPACE=${TARGET_NAMESPACE:-"open-cluster-management"}
storageSecret=${STORAGE_SECRET_NAME:-"storage-secret"}
ready=$(kubectl get secret $storageSecret -n $TARGET_NAMESPACE --ignore-not-found=true)
if [ ! -z "$ready" ]; then
  echo "storageSecret $storageSecret already exists in $TARGET_NAMESPACE namespace"
  exit 0
fi

# install postgres operator
POSTGRES_OPERATOR=${POSTGRES_OPERATOR:-"pgo"}
kubectl apply -k ${currentDir}/postgres-operator
kubectl -n postgres-operator wait --for=condition=Available Deployment/$POSTGRES_OPERATOR --timeout=1000s
echo "$POSTGRES_OPERATOR operator is ready!"

# deploy postgres cluster
pgnamespace="hoh-postgres"
userSecret="hoh-pguser-postgres"

kubectl apply -k ${currentDir}/postgres-cluster
matched=$(kubectl get secret $userSecret -n $pgnamespace --ignore-not-found=true)
SECOND=0
while [ -z "$matched" ]; do
  if [ $SECOND -gt 300 ]; then
    echo "Timeout waiting for creating $secret"
    exit 1
  fi
  echo "Waiting for secret $userSecret to be created in pgnamespace $pgnamespace"
  matched=$(kubectl get secret $userSecret -n $pgnamespace --ignore-not-found=true)
  sleep 5
  (( SECOND = SECOND + 5 ))
done
echo "Postgres is ready!"

# create usersecret for postgres
databaseHost="$(kubectl get secrets -n "${pgnamespace}" "${userSecret}" -o go-template='{{index (.data) "host" | base64decode}}')"
databasePort="$(kubectl get secrets -n "${pgnamespace}" "${userSecret}" -o go-template='{{index (.data) "port" | base64decode}}')"
databaseUser="$(kubectl get secrets -n "${pgnamespace}" "${userSecret}" -o go-template='{{index (.data) "user" | base64decode}}')"
databasePassword="$(kubectl get secrets -n "${pgnamespace}" "${userSecret}" -o go-template='{{index (.data) "password" | base64decode}}')"
databasePassword=$(printf %s "$databasePassword" |jq -sRr @uri)

databaseUri="postgres://${databaseUser}:${databasePassword}@${databaseHost}:${pgAdminPort}/hoh"

kubectl create secret generic $storageSecret -n $TARGET_NAMESPACE \
    --from-literal=database_uri=$databaseUri
echo "storage secret is ready in $TARGET_NAMESPACE namespace!"