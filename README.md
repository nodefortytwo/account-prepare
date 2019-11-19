# AWS Account Prepare
simple script to prepare a newly created AWS Account.

## Run
The program expects a valid aws session to paseed.

I recommend `aws-vault` like so:

`aws-vault exec <<profile>>  -- go run main.go `

## Tasks
1. Delete default VPCs
2. Enable long arn formats on ECS