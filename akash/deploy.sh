#!/bin/bash
if [ $# -lt 2 ]; then
    echo "invalid arguments"
    exit 1
fi
# set environment variables
ENV_SHELL=${0/deploy.sh/env.sh} 
echo "init env"
source $ENV_SHELL $1
SDLFILE=$2
IMAGE=$(cat $SDLFILE| grep "image:" | yq '.image')
echo -e "Deployment: $IMAGE\n"

# generate & publish cert
echo "Generating & publishing cert"
provider-services tx cert generate client --from $AKASH_KEY_NAME --overwrite
provider-services tx cert publish client --from $AKASH_KEY_NAME -y >/dev/null 2>&1
echo -e "success\n"

# create deployment
echo "Start creating deployment: $IMAGE"

CREATE_LOG=create_deployment.log
provider-services tx deployment create $SDLFILE --from $AKASH_KEY_NAME -y > $CREATE_LOG
export AKASH_DSEQ="$(cat $CREATE_LOG | jq -r '.logs[].events[].attributes[] | select(.key == "dseq") | .value' | sed -n '1p')"
if [ -z "$AKASH_DSEQ" ]; then 
    echo "not found valid dseq"
    exit 1
fi
export AKASH_OSEQ="$(cat $CREATE_LOG | jq -r '.logs[].events[].attributes[] | select(.key == "oseq") | .value' | sed -n '1p')"
export AKASH_GSEQ="$(cat $CREATE_LOG | jq -r '.logs[].events[].attributes[] | select(.key == "gseq") | .value' | sed -n '1p')"
echo "    dseq: $AKASH_DSEQ"
echo "    oseq: $AKASH_OSEQ"
echo "    gseq: $AKASH_GSEQ"
echo

echo "Waitting for providers bids..."
SLEEP_TIME=15
if [ $# -eq 3 ]; then
   SLEEP_TIME=$3
fi
sleep $SLEEP_TIME

# query provider bids
echo "Querying provider bids..."
BIDS_LOG=bids.log
provider-services query market bid list --owner=$AKASH_ACCOUNT_ADDRESS --node $AKASH_NODE --dseq $AKASH_DSEQ --state=open > $BIDS_LOG
echo "Choosing the best provider..."
BIDS_SORTED_LOG=bids-sorted.log
cat $BIDS_LOG | yq '.bids| sort_by(.bid.escrow_account.balance.amount) | reverse | sort_by(.bid.price.amount)' > $BIDS_SORTED_LOG
export AKASH_PROVIDER="$(cat $BIDS_SORTED_LOG | yq '.[0].escrow_account.owner')"
PRICE="$(cat $BIDS_SORTED_LOG | yq '.[0].bid.price.amount')"
AMOUNT="$(cat $BIDS_SORTED_LOG | yq '.[0].escrow_account.balance.amount')"
echo "The best provider is:"
echo "    Address: $AKASH_PROVIDER"
echo "    Bid Price: $PRICE"
echo "    Available Balance: $AMOUNT"
if [ -z "$AKASH_PROVIDER" ]; then 
    echo "not found valid provider"
    exit 1
fi
echo 

# create lease
echo "Start creating the lease"
LEASE_LOG=lease.log
provider-services tx market lease create --dseq $AKASH_DSEQ --provider $AKASH_PROVIDER --from $AKASH_KEY_NAME -y > $LEASE_LOG
provider-services send-manifest $SDLFILE --dseq $AKASH_DSEQ --provider $AKASH_PROVIDER --from $AKASH_KEY_NAME >> $LEASE_LOG
echo "Start creating the lease"
echo "    provider: $AKASH_PROVIDER"
echo "    status:       PASS"
echo

# query deployment url
echo "Waitting for lease url..."
sleep $SLEEP_TIME
echo "Querying the lease url..."
LEASE_STATUS_LOG=lease-status.log
provider-services lease-status --dseq $AKASH_DSEQ --from $AKASH_KEY_NAME --provider $AKASH_PROVIDER > $LEASE_STATUS_LOG
LEASE_URL="$(cat $LEASE_STATUS_LOG | jq -r '.services.web.uris[]| select(. != null)')"
echo "Deploy is finished!"
echo "    Provider: $AKASH_PROVIDER"
echo "    Price: $PRICE"
echo "    URIs: $LEASE_URL"