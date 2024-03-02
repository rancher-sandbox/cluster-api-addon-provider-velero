AZURE_BACKUP_RESOURCE_GROUP=dgrigore_backups
az group create -n $AZURE_BACKUP_RESOURCE_GROUP --location francecentral

AZURE_STORAGE_ACCOUNT_ID="velerodgrigorev" #-$(uuidgen | cut -d '-' -f5 | tr '[A-Z]' '[a-z]')"
az storage account create \
    --name $AZURE_STORAGE_ACCOUNT_ID \
    --resource-group $AZURE_BACKUP_RESOURCE_GROUP \
    --sku Standard_GRS \
    --encryption-services blob \
    --https-only true \
    --min-tls-version TLS1_2 \
    --kind BlobStorage \
    --access-tier Hot

BLOB_CONTAINER=velero-dgrigore
az storage container create -n $BLOB_CONTAINER --public-access off --account-name $AZURE_STORAGE_ACCOUNT_ID

AZURE_STORAGE_ACCOUNT_ACCESS_KEY=`az storage account keys list --account-name $AZURE_STORAGE_ACCOUNT_ID --query "[?keyName == 'key1'].value" -o tsv`
cat << EOF  > ./credentials-velero-az
AZURE_STORAGE_ACCOUNT_ACCESS_KEY=${AZURE_STORAGE_ACCOUNT_ACCESS_KEY}
AZURE_CLOUD_NAME=AzurePublicCloud
EOF

kubectl create secret generic -n default az-credentials --from-file=azure=credentials-velero-az

rm credentials-velero-az
