# go-aws-lambda-logshipper
## shipping aws cloudwatch log stream to your log collector endpoint (over tcp)

## Variable.
bucket=NAME-OF-S3-BUCKET

## Create binary that we told Lambda to use as the handler (via Cloud Formation).
```
GOOS=linux go build -o main main.go
zip deployment.zip main
```

## Trigger a Cloud Formation update to update the infrastructure and the deployed code.
```
aws s3 mb s3://$(bucket)
aws cloudformation package \
    --template-file formation.yml \
    --output-template-file formation.compiled.yml \
    --s3-bucket $(bucket)
aws cloudformation deploy \
    --template-file formation.compiled.yml \
    --stack-name $(bucket) --capabilities CAPABILITY_IAM
```