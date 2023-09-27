#!/bin/bash
AKASH_KEY_NAME=$1
AKASH_ACCOUNT_ADDRESS="$(provider-services keys show $AKASH_KEY_NAME -a)"
provider-services tx deployment close --dseq $2  --owner $AKASH_ACCOUNT_ADDRESS --from $AKASH_KEY_NAME