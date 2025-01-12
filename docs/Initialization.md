# Initialization Guide

This document outlines the process for initializing the server during its first startup.

## Overview

When the server is started for the first time, it detects an uninitialized state and automatically enters **Initialization Mode**. During this mode, the server generates a unique token and prints it to `STDOUT`. This token, along with an admin username and password, must be provided to initialize the server.

Once initialization is complete:

- The token becomes invalid.
- The initialization endpoint is permanently disabled.

## Steps to Initialize the Server

### Start the Server in Initialization Mode

Start the server. Upon detecting that it is uninitialized, the server will enter Initialization Mode and display a token.

Example:
```shell
$ ./hoodie
...
Initialization Mode Enabled, Token: asdasdasd
...
```

### Create an Admin Account

Use the generated token to create an admin account by sending a POST request to the server's initialization endpoint.

Example:
```shell
curl -X POST \
    -H 'Authorization: Bearer {token}' \
    -H 'Content-Type: application/json' \
    -d '{"username": "admin", "password": "admin1234"}' \
    https://hoodie-server/api/v0.1/init/
```

- Replace `{token}` with the token displayed during Initialization Mode.
- Customize the `username` and `password` fields as required.
