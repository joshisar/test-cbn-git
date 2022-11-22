# Empty GoLang

## Project structure

This project is set up based on SettleMint's experience and knowledge in Hyperledge Fabric projects aiming to provide the best
developer experience possible. In the Explorer panel on the left, you will find all the files and folders of which we describe below.

```plain
root
|    README.md                          # This file
|
+--+ network                            # This folder contains the network configuration
|  |--+ ca                              # CA TLS certificate
|  |--+ orderer                         # Orderer node TLS certificate
|  |--+ organization                    # Your organization's user's MSP + TLS certificates
|  |  +--+ admin                        # Admin user's MSP + TLS certificates and NodeOus configuration
|  |  +--+ peer                         # Peer user's MSP + TLS certificates and NodeOus configuration
+--+ node
|  |--+ config                          # Peer node's configuration folder A.K.A. "FABRIC_CFG_PATH"
|  |  +--+ core.yaml                    # Peer node's configuration file
|  |  +--+ configtx.yaml                # File that contains the information that is required to build the channel configurations
|
+--+ src                                # This folder contains your Chaincode smart contracts written in Typescript.
|  |--+ assetTransfer                   # This smart contract folder was set up to show how you should structure
|                                       # your project. You can delete it and create your own smart contract folder.
|                                       # You can have as many smart contract folders as you want, think of them as
|                                       # packages.
|
+--+ scripts                            # This folder contains the scripts to deploy the chaincode to your
|                                       # blockchain node.
|
+--+ .env                               # This file contains the environment variables you want to use when running
|                                       # the scripts. These variables are crucial to the scripts to work.
|                                       # The most important variables are:
|                                       # - CC_SRC_PATH: This is the path to your chaincode smart contract folder.
|                                       # - CC_NAME: This is the name of your chaincode.
|                                       # - CC_INIT_FCN: This is the function that will be called when initializing your chaincode.
|
```

## Tasks

While the terminal (shortcut: `^⇧` or via the hamburger menu top left corner) is a fully functional linux based terminal to execute all scripts folder's `.sh` files, on the left bottom in the Explorer panel you will find the **TASK EXPLORER** panel that provides one click access to predefined commands for the most common actions, expand the `bash (scripts)` section.

The following tasks run using the Admin user and the environment variables defined in the `.env` file.

### package-chaincode (runs `./scripts/package-chaincode.sh`)

It packages the smart contract code into a gzipped tar file and stores it in the smart contract folder.

### install-chaincode (runs `./scripts/install-chaincode.sh`)

It installs the chaincode definition on the peer node.

### approve-chaincode (runs `./scripts/approve-chaincode.sh`)

It approves for your organization the chaincode definition previously installed on the peer node.

### check-chaincode commit (runs `./scripts/check-commit-chaincode.sh`)

It checks the approval status of the chaincode definition.

### commit-chaincode (runs `./scripts/commit-chaincode.sh`)

Run this task when a sufficient number of channel member have approved the chaincode definition. If the chaincode definition has been approved by the majority (or all, it depends on the network's configuration) of the channel members, it will commit the chaicode definition.

### invoke-init-chaincode (runs `./scripts/invoke-init-chaincode.sh`)

Run this task when the chaincode definition has been commited and if it has an init function. Adjust (if needed) the arguments passed to the init function in the `./scripts/invoke-init-chaincode.sh` file.

## Extra utilities

### `./scripts/register-client-user.sh`

Run this script to register a new user of type `client` on the network.
