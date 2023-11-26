[![main](https://github.com/alphauslabs/iam/actions/workflows/main.yml/badge.svg)](https://github.com/alphauslabs/iam/actions/workflows/main.yml)

`iam` is the command line client for our internal **IAM** service [(iamd)](https://github.com/mobingilabs/ouchan/tree/master/cloudrun/iamd).

To install using [HomeBrew](https://brew.sh/), run the following command:

```bash
$ brew install alphauslabs/tap/iam
```

To setup authentication, set your `GOOGLE_APPLICATION_CREDENTIALS` env variable using your credentials file. You can validate by running the following command:

```bash
$ iam whoami
```

Explore more available subcommands and flags though:

```bash
$ iam -h
# or
$ iam <subcmd> -h
```
