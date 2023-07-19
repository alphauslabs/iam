`iam` is the command line client for our internal **IAM** service [(iamd)](https://github.com/mobingilabs/ouchan/tree/master/cloudrun/iamd).

To install using [HomeBrew](https://brew.sh/), run the following command:

```bash
$ brew install alphauslabs/tap/iam
```

To setup authentication, set your `GOOGLE_APPLICATION_CREDENTIALS` env variable. You can validate by running the following command:

```bash
$ iam whoami
```

Run `iam -h` or `iam <subcommand> -h` to know more about the available subcommands and flags.
