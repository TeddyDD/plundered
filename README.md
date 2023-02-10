# Plundered 🏴‍☠

Various Go packages extracted from bigger projects and modified for my
use cases.

Do what you want cause a pirate is free, you are a pirate! 

## List of packages

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
    - removed test since they use Ginko/Gomega
