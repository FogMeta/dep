#!/bin/bash
AKASH_KEY_NAME=$1
echo $AKASH_KEY_NAME
AKASH_KEYRING_BACKEND=os
provider-services keys add $AKASH_KEY_NAME
export AKASH_ACCOUNT_ADDRESS="$(provider-services keys show $AKASH_KEY_NAME -a)"
echo " ----------------------------------------------------------------------------------"
echo "| your akash account address: $AKASH_ACCOUNT_ADDRESS         |"                             
echo "| you need get funds into your account before using it.                            |"
echo " ----------------------------------------------------------------------------------"
