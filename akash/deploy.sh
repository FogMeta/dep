#!/bin/bash
source env.sh $1
SDLFILE=$2
provider-services tx cert generate client --from $AKASH_KEY_NAME --overwrite
provider-services tx cert publish client --from $AKASH_KEY_NAME -y
CREATE_LOG=create_deployment.log
provider-services tx deployment create $SDLFILE --from $AKASH_KEY_NAME -y > $CREATE_LOG
export AKASH_DSEQ="$(cat $CREATE_LOG | jq -r '.logs[].events[].attributes[] | select(.key == "dseq") | .value' | sed -n '1p')"
if [ -z "$AKASH_DSEQ" ]; then 
    echo "not found valid dseq"
    exit 1
fi
export AKASH_OSEQ="$(cat $CREATE_LOG | jq -r '.logs[].events[].attributes[] | select(.key == "oseq") | .value' | sed -n '1p')"
export AKASH_GSEQ="$(cat $CREATE_LOG | jq -r '.logs[].events[].attributes[] | select(.key == "gseq") | .value' | sed -n '1p')"
echo $AKASH_DSEQ $AKASH_OSEQ $AKASH_GSEQ
BIDS_LOG=bids.log
provider-services query market bid list --owner=$AKASH_ACCOUNT_ADDRESS --node $AKASH_NODE --dseq $AKASH_DSEQ --state=open > $BIDS_LOG
export AKASH_PROVIDER="$(cat $BIDS_LOG | yq '.bids| sort_by(.bid.escrow_account.balance.amount) | reverse|sort_by(.bid.price.amount)'| yq '.[0].escrow_account.owner')"
echo $AKASH_PROVIDER
if [ -z "$AKASH_PROVIDER" ]; then 
    echo "not found valid provider"
    exit 1
fi
provider-services tx market lease create --dseq $AKASH_DSEQ --provider $AKASH_PROVIDER --from $AKASH_KEY_NAME -y
provider-services send-manifest $SDLFILE --dseq $AKASH_DSEQ --provider $AKASH_PROVIDER --from $AKASH_KEY_NAME
