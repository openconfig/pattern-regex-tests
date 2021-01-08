# Pattern Statement Tests for OpenConfig YANG models

Tests [pattern statement](https://tools.ietf.org/html/rfc7950#section-9.4.5)
using a [pyang](https://github.com/mbj4668/pyang) plugin and
[oc-ext:posix-pattern](https://github.com/openconfig/public/blob/master/release/models/openconfig-extensions.yang#L114)
using [goyang](https://github.com/openconfig/goyang).

## Releases

Releases are synchronized with the current OpenConfig YANG models. If new
breaking tests are added (e.g. a test cases handled incorrectly by a current
pattern regex), then the minor version must be incremented.

At this time, major version updates are not anticipated, but could occur as a
result of major changes to the repository.

--------------------------------------------------------------------------------

[OpenConfig YANG models](https://github.com/openconfig/public/blob/master/README.md)
