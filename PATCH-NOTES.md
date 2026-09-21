# CAP_SYS_RESOURCE fix

This fork contains a small Docker plugin capability fix for LINSTOR-backed ext4 volumes.

## Problem

When Docker mounted a LINSTOR volume, the plugin mounted the DRBD device and then invoked `resize2fs` for the ext4 filesystem.

With the upstream plugin configuration, the resize failed with:

```text
resize2fs: Permission denied to resize filesystem
```

After that failed mount attempt, the DRBD device could remain open, causing follow-up errors such as:

```text
unable to get exclusive open on /dev/drbdXXXX
```

The issue was reproduced independently of the application with a fresh 10 GiB LINSTOR volume and a minimal Alpine container.

## Fix

The upstream plugin configuration grants:

```json
"capabilities": ["CAP_SYS_ADMIN"]
```

This fork adds `CAP_SYS_RESOURCE`:

```json
"capabilities": ["CAP_SYS_ADMIN", "CAP_SYS_RESOURCE"]
```

## What CAP_SYS_RESOURCE does here

Linux capabilities split privileged root operations into smaller permissions.

The LINSTOR Docker plugin runs inside Docker's managed plugin sandbox. Even though it performs privileged storage operations, it can only use the Linux capabilities explicitly granted in `config.json`.

In this setup, the ext4 online-resize operation triggered by `resize2fs` required the additional `CAP_SYS_RESOURCE` capability. Without it, the kernel rejected the resize with `EPERM` / `Permission denied`.

Adding `CAP_SYS_RESOURCE` does **not** change LINSTOR replication, DRBD configuration, volume size, or the storage pool. It only gives the plugin process the permission required to complete its existing ext4 online-resize step.

## Verification

A fresh 10 GiB LINSTOR volume failed with the upstream plugin and succeeded with the patched plugin.

Test command:

```bash
docker run --rm \
  --name capfix-test \
  -v capfix_test:/data \
  alpine:3.20 \
  sh -c 'echo "CAPFIX WORKS" > /data/test.txt && cat /data/test.txt'
```

Expected output:

```text
CAPFIX WORKS
```

## Scope

This patch was added to fix the observed ext4 mount/resize failure in this Docker Swarm + LINSTOR setup.

The relevant change is intentionally limited to the Docker plugin capability list.
# CAP_SYS_RESOURCE fix

This fork contains a small Docker plugin capability fix for LINSTOR-backed ext4 volumes.

## Problem

When Docker mounted a LINSTOR volume, the plugin mounted the DRBD device and then invoked `resize2fs` for the ext4 filesystem.

With the upstream plugin configuration, the resize failed with:

```text
resize2fs: Permission denied to resize filesystem
```

After that failed mount attempt, the DRBD device could remain open, causing follow-up errors such as:

```text
unable to get exclusive open on /dev/drbdXXXX
```

The issue was reproduced independently of the application with a fresh 10 GiB LINSTOR volume and a minimal Alpine container.

## Fix

The upstream plugin configuration grants:

```json
"capabilities": ["CAP_SYS_ADMIN"]
```

This fork adds `CAP_SYS_RESOURCE`:

```json
"capabilities": ["CAP_SYS_ADMIN", "CAP_SYS_RESOURCE"]
```

## What CAP_SYS_RESOURCE does here

Linux capabilities split privileged root operations into smaller permissions.

The LINSTOR Docker plugin runs inside Docker's managed plugin sandbox. Even though it performs privileged storage operations, it can only use the Linux capabilities explicitly granted in `config.json`.

In this setup, the ext4 online-resize operation triggered by `resize2fs` required the additional `CAP_SYS_RESOURCE` capability. Without it, the kernel rejected the resize with `EPERM` / `Permission denied`.

Adding `CAP_SYS_RESOURCE` does **not** change LINSTOR replication, DRBD configuration, volume size, or the storage pool. It only gives the plugin process the permission required to complete its existing ext4 online-resize step.

## Verification

A fresh 10 GiB LINSTOR volume failed with the upstream plugin and succeeded with the patched plugin.

Test command:

```bash
docker run --rm \
  --name capfix-test \
  -v capfix_test:/data \
  alpine:3.20 \
  sh -c 'echo "CAPFIX WORKS" > /data/test.txt && cat /data/test.txt'
```

Expected output:

```text
CAPFIX WORKS
```

## Scope

This patch was added to fix the observed ext4 mount/resize failure in this Docker Swarm + LINSTOR setup.

The relevant change is intentionally limited to the Docker plugin capability list.
