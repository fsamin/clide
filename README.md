# Clide

Cloud storage files management CLI

Supported providers:

- Openstack Swift: `swift`
- Amazon S3: `s3`

## Commands

### Upload

Upload files to a container/bucket with:

```bash
    clide <provider> upload <file 0> [file 1] ... [file n] <destination>
```

- `provider` must be on the supported providers
- `destination` is the destination container/bucket. Is it doesn't exist, it will be created by default as a private container/bucket.

## About Authentication

Authentication settings can be set with command flags or environment variables. We strongly suggest to set environment variables to use `clide`easily.

### Openstack swift environment variables

```bash
    OS_USERNAME
    OS_PASSWORD
    OS_TENANT_NAME
    OS_AUTH_URL
```

### Amazon S3 environment variables

```bash
    AWS_ACCESS_KEY_ID
    AWS_SECRET_ACCESS_KEY
    AWS_DEFAULT_REGION
```