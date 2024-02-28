aws s3api create-bucket \
    --bucket $BUCKET \
    --region $REGION \
    --create-bucket-configuration LocationConstraint=$REGION

aws iam create-user --user-name dgrigorev-velero

cat > velero-policy.json <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "ec2:DescribeVolumes",
                "ec2:DescribeSnapshots",
                "ec2:CreateTags",
                "ec2:CreateVolume",
                "ec2:CreateSnapshot",
                "ec2:DeleteSnapshot"
            ],
            "Resource": "*"
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:GetObject",
                "s3:DeleteObject",
                "s3:PutObject",
                "s3:AbortMultipartUpload",
                "s3:ListMultipartUploadParts"
            ],
            "Resource": [
                "arn:aws:s3:::${BUCKET}/*"
            ]
        },
        {
            "Effect": "Allow",
            "Action": [
                "s3:ListBucket"
            ],
            "Resource": [
                "arn:aws:s3:::${BUCKET}"
            ]
        }
    ]
}
EOF

aws iam put-user-policy \
  --user-name dgrigorev-velero \
  --policy-name dgrigorev-velero \
  --policy-document file://velero-policy.json

out=`aws iam create-access-key --user-name dgrigorev-velero`
key=`echo "$out" | jq -r ".AccessKey.AccessKeyId"`
secret=`echo "$out" | jq -r ".AccessKey.SecretAccessKey"`

cat > credentials-velero <<EOF
[default]
aws_access_key_id=$key
aws_secret_access_key=$secret
EOF