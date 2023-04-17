# Plundered 🏴‍☠

Arr, me hearty!

This be me booty of Go packages, pried free from grander projects and plundered for me own use.

Sail where the wind takes ya! Pirates live free, ya scallywag!

## Usage

```go
import "go.teddydd.me/plundered/"
```

[Docs](https://godocs.io/go.teddydd.me/plundered)

## Source code

Ye best be finding this here repository at https://git.sr.ht/~teddy/plundered
If ye got patches, send 'em on over to https://lists.sr.ht/~teddy/public-inbox
And if ye prefer a read-only mirror, head to https://github.com/TeddyDD/plundered

## 'Tis a tally of our plundered packages:

- recorder
    - https://github.com/carlmjohnson/requests
    - round tripper for recording and replaying HTTP requests
    - added filter functions that let you delete some parts of request and
    response objects (like credentials) before saving them as fixtures
- templates
    - https://github.com/gofiber/template/tree/master/html
    - use `io.FS`
    - removed `Parse`
- signals
    - https://github.com/kubernetes-sigs/controller-runtime/tree/master/pkg/manager/signals
    - aye, removed test since they use Ginko/Gomega
- cmdline
    - https://pkg.go.dev/github.com/u-root/u-root@v0.11.0/pkg/cmdline#CmdLine.Flag
    - for parsing simple cmdline arguments in kernel style `foo=bar feature1 load-cannons=yes`
    - only CmdLine struct
    - new AsBool method with simple heuristics
