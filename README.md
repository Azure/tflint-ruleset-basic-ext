# TFLint Ruleset for terraform-provider-basic-ext

TFLint ruleset extension plugin for common terraform code syntax check

## Requirements

- TFLint v0.35+
- Go v1.18

## Building the plugin

Clone the repository locally and run the following command:

```
$ make
```

You can easily install the built plugin with the following:

```
$ make install
```

Note that if you install the plugin with make install, you must omit the `version` and `source` attributes in `.tflint.hcl`:

```hcl
plugin "basic-ext" {
    enabled = true
}
```

If you install the plugin with `tflint --init`, you must specify the following content in `.tflint.hcl`:
```hcl
plugin "basic-ext" {
    enabled = true
    version = "0.0.1"
    source  = "github.com/DrikoldLun/tflint-ruleset-basic-ext"
    signing_key = <<-KEY
    -----BEGIN PGP PUBLIC KEY BLOCK-----

    mQGNBGLtC9wBDADUty0QYzsGe7+Q1e2slFNFw6+SVQKvj+/K3CGrKt1sl8evPLOy
    Q+kNUmCtTSunPgRlxPWICNa/QiV5AhR5Jn7B0o92OG8LjHJ1ThrvzbHVDyCkKPb/
    PzpsSNjnWVxZV3uTO3brzaUOn9nbmfOGBtJQc4V2a9U8Z8X3WVh2LCn3qXb6xRtM
    SP/DIccZHChMQpgYheh3kZIPOvRyLA4Kp8MRB1MOi44aQyfneOA5ryVa7yLNHn1h
    QyKQ0HiGGnCjYS5NEtyBWtzm2xRWEmJRWydsDDDvBLI+44ymhriwZUrS8FVUgiYo
    GjIGp6JOxP/V3HqfCP22I1sfheXFX4jzr0qWasKYGoSn9CoAUxvoO1k19jLCL3cV
    Fw1JC26wIHlTE9GBLNWklgfBJapK34Iqwy9Q5SAuw6gi5ZoGN5kVMqmzWacsxcZZ
    MAAhOdSwOsGNx/Zx6D2eJ8BWsybvdz8RQ3AXx7YCGSwHgNbcgjjRKNQlxs8X60XN
    bQYip6Xu8V2S3Y0AEQEAAbQjTHVuIFpoYW5nIDx6aGFuZ2xuanUyMDE1QGdtYWls
    LmNvbT6JAdEEEwEIADsWIQSNp1KuIvc/Uuj2uScp5xG7/lWJIgUCYu0L3AIbAwUL
    CQgHAgIiAgYVCgkICwIEFgIDAQIeBwIXgAAKCRAp5xG7/lWJIjZxC/0SG46lBuyj
    8og2fdcCTiKdI1IllIjKnrHjh1Hp92ohfD+N9Q729qwN7CmC/MSlhLB2bTUn6HZ0
    Pp6Go+h5V0mGO9CtN15icbvhEpO+k48wGLfX9qNfc5Wod2a8zkyTZ2ZAlEudxEwT
    A7mdpGhi6zLCnsSlC+/A293KnzyIm3/8GwZbOoGBmPOFYb2FLCRvIupFj+VcOsDz
    5ZnmrLVNSC69dDmcXRqRpgdd60J7wS093qqsd0G67Mj7FCxOc4BLfDMKE/MYSJmq
    xRpFak6hjQVYGmUQhkyu+pPAdz+oR5fRCDwlmjtClnTE2IfP6bi4IX3+g+s3GJU0
    bHSh1ImWUvnjJ/G/ZmeYtWWAxt1ZzJj0fjnF6nS+uzh/9bstY4EklMd8cY5v/vmU
    ittjB1KH5IJueDHL5JNoWi0tJVh8vSLuM3HbZ7yWfvE+uaextyfpOPJmIt+z5yOM
    4MHeu12TCIxKrSkdq/Z+3b1FWQTUDPDBr/41fpt22H+Lxf4F2Jcn7Tc=
    =T2nt
    -----END PGP PUBLIC KEY BLOCK-----
    KEY
}
```

If you want to specify the ignore/retain filename patterns for some rules in this plugin, write the following ignore config in `.tflintignore.basic-ext.json`: 
```json
{
  "terraform_required_version_declaration" : [
    "^*$",
    "!^main\\.tf$"
  ]
}
```
The pattern regular expression follows Go syntax, and the prefix `!` means the file with such name pattern would still be checked.

Follow the instructions to edit the generated files and open a new pull request.