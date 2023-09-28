#!/bin/bash
AKASH_KEY_NAME=$1
AKASH_ACCOUNT_ADDRESS="$(provider-services keys show $AKASH_KEY_NAME -a)"
CODE="$(provider-services tx deployment close --dseq $2  --owner $AKASH_ACCOUNT_ADDRESS --from $AKASH_KEY_NAME | jq '.code')" 
if [ "$CODE" -eq "0" ] ; then
    echo "deployment $2 closed"
else
    echo "failed to close, code: $CODE"
fi