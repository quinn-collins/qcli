# QCLI

## Install

With Go
```
go install github.com/quinn-collins/qcli@latest
```

With Homebrew
```
brew install quinn-collins/tap/qcli
```

## Configuration

Configuration operates in order of precedence - flags, environment variables, and configuration file.

Flags:

```bash
> qcli -h
NAME
  qcli

  -h, --help                          Display help information.

QCLI COMMANDS
  qcli mfa

  -h, --help                          Display help information.
  -p, --aws-profile STRING            The credentials profile used by AWS. (default: default)
  -t, --aws-target-profile STRING     Sets the target for AWS credentials MFA. (default: default)
  -r, --aws-region STRING             The region AWS commands will operate within. (default: us-east-1)

  qcli list-buckets

  -h, --help                          Display help information.
  -p, --aws-profile STRING            The credentials profile used by AWS. (default: default)
  -r, --aws-region STRING             The region AWS commands will operate within. (default: us-east-1)
```

Environment variables:
Can be specified by prefixing a flag with `QCLI_` and replacing hyphens with underscores, e.g. `--aws-target-profile` is set via `QCLI_AWS_TARGET_PROFILE`.

Configuration file:
Expected to be found at `$HOME/.config/qcli/config.yaml`. Items in the configuration file can be specified using the flags minus the `--` prefix, e.g. `--aws-target-profile` is set via `aws-target-profile: my-named-profile`


## AWS Integration

Currently supports authentication to AWS accounts and Localstack using the [configuration and credentials](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html#cli-configure-files-format) files.

## Datadog Integration
## Github Integration
## Slack Integration
## Sonar Integration

