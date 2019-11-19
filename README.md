# AWS Account Prepare
simple script to prepare a newly created AWS Account.

The program expects a valid aws session to paseed.

I recommend `aws-vault` like so:

`aws-vault exec <<profile>>  -- go run main.go `