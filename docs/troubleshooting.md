# Troubleshooting

## Docker daemon is unavailable

Typical messages include `Cannot connect to the Docker daemon` or a missing Docker Desktop pipe. Start Docker Desktop or the system Docker service, then run:

```bash
docker version
docker compose version
```

Both commands must succeed before the artifact can start.

## Container image cannot be pulled

Confirm the exact lowercase image name and version:

```bash
docker pull ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.2.1
```

No registry login is required for the public artifact. A `manifest unknown` response normally means the tag was typed incorrectly or the release workflow has not completed.

## Apple Silicon or another ARM host

The release workflow publishes Linux AMD64 and ARM64 manifests. Update Docker Desktop if the client selects the wrong platform. As a fallback, add `--platform linux/amd64` to a direct `docker run` command.

## Bind-mounted output is not writable

On Linux, run the container with the host user identity:

```bash
docker run --rm --user "$(id -u):$(id -g)" \
  --mount "type=bind,source=$PWD/output,target=/output" \
  ghcr.io/mb3r-lab/mb3r-software-complex-repro:toolchain-1.2.1 \
  --repeat 2 --work-dir /output
```

Ensure `output/` exists first.

## Reference mismatch

A mismatch is reported with the semantic file and expected/actual SHA-256 values. Do not edit an individual reference file or run the update command merely to make the check pass.

Check, in order:

1. the container tag is `toolchain-1.2.1`;
2. local case files have not replaced files inside a bind-mounted artifact root;
3. a source run uses the exact commits in `toolchain-lock.json`;
4. the mismatch is not caused by an intentional reviewed toolchain change.

Only maintainers accepting a new toolchain release should use `make reproduce-paper-update`.

## Source revision mismatch

This error applies only to the maintainer source-build path. It means one of the sibling repositories is not at its locked commit. The runner will stop before building or executing that component. The container path is unaffected because the binaries are built by CI from locked commits.

## Where to report a problem

Open an issue in the [artifact repository](https://github.com/MB3R-Lab/mb3r-software-complex-repro/issues) and include:

- host operating system and architecture;
- `docker version` output;
- the exact command;
- the first mismatch or error message;
- whether the short demo or complete reproduction was selected.
