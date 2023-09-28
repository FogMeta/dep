# DEP

## Project Overview

DEP (Deployment Enabler Platform) is a cross-platform application deployment project designed to help users easily deploy their applications to the [AKASH Network](https://akash.network) and [Lagrange Platform](https://lagrangedao.org). DEP offers various deployment methods, including building images using Dockerfiles and deploying applications using SDL (YAML) files.

## Key Features

DEP provides the following key features:

1. **Building Image From Dockerfile:** DEP allows users to build application images using Dockerfiles, making it easy to deploy and manage applications.

2. **Deploy to AKASH Network:** DEP supports deploying your applications to the [AKASH Network](https://akash.network), a distributed cloud computing platform that provides a reliable hosting environment for your applications.

3. **Deploy to Lagrange Platform:** DEP integrates with the [Lagrange Platform](https://lagrangedao.org), enabling you to deploy your applications to this powerful blockchain platform.(**coming soon**)

4. **One-Click Deploy Lagrange Space to AKASH Network:** DEP offers one-click deployment, allowing you to quickly deploy applications, including Dockerfile and images, from Lagrange Platform to the [AKASH Network](https://akash.network), simplifying the cross-platform deployment process.

5. **Retrieve Access URLs:** DEP provides the functionality to retrieve access URLs for your deployed applications.

## Prerequisites

Before using DEP, ensure that you have the following prerequisites in place:

- [Docker](https://www.docker.com/) installed and started
- A valid account on [AKASH Network](https://docs.akash.network/guides/cli/detailed-steps/part-2.-create-an-account) if you intend to deploy there.
- [GO](https://go.dev/dl/) version v1.19.0+ installed
- [Lag-cli](https://github.com/lagrangedao/lagrange-cli) installed
- [Akash-cli](https://docs.akash.network/guides/cli/detailed-steps/part-1.-install-akash) installed 
- [jq](https://jqlang.github.io/jq/download/) installed 
- [yq](https://github.com/mikefarah/yq#install) installed

## Compilation and Installation

To compile and install DEP, follow these steps:

1. Clone the [DEP](https://github.com/fogmeta/dep) repository to your local machine.
   ```
   git clone https://github.com/FogMeta/dep.git
   ```

3. Navigate to the DEP project directory `cd dep`.

4. Run `go build`.

## Framework

<img width="979" alt="image" src="https://github.com/FogMeta/dep/assets/102578774/91a9f49f-a0af-44e5-8a4c-5062e097dd3b">


## Usage

### Init 

init configuration

```bash
dep init 
```

after `init`, `dep.conf` file generated in current directory, which content like below:

```yml
work_dir = "."

[registry]
  server_address = ""
  user_name = ""
  password = ""
```

| parameter      | description              |
| -------------- | ------------------------ |
| work_dir       | directory to build image |
| registry       | registry to pull image   |
| server_address | registry server          |
| user_name      | registry user name       |
| password       | registry password        |

you can set it in your config

### Build

build docker image

```bash
dep build [lag_url]
```

### Create Account

create a new account

```bash
dep create-account [account_name]
```

result like below:

```bash
**Important** write this mnemonic phrase in a safe place.
It is the only way to recover your account if you ever forget your password.

****
 ----------------------------------------------------------------------------------
| your akash account address: akash**         |
| you need get funds into your account before using it.                            |
 ----------------------------------------------------------------------------------
```

### Deploy

deploy the deployment file

```bash
dep deploy [account_name] [deployment_file]
```

if deployed successfully, you can get the deployment's `desq`,`provider`,`price` and `uris`

### Status

query the status of the deployment

```bash
dep status [account_name] [dseq] [provider]
```

### Close

close the deployment

```bash
dep close [account_name] [dseq]
```

## Use Case
Here is a case to deploy the `Hellow World` in [Lagrange Platform](https://lagrangedao.org) to the [AKASH Network](https://akash.network) using `Dockerfile`


```


```


## Contributions and Support

If you're interested in the DEP project and would like to contribute code or report issues, please visit our [GitHub repository](https://github.com/fogmeta/dep). If you have any questions or need technical support, feel free to contact our team.

Thank you for choosing DEP, and we look forward to seeing your applications successfully deployed and running!
