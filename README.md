# Plundered 🏴‍☠

Various Go packages extracted from bigger projects and modified for my
use case.

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
