# Pattern Statement Tests for OpenConfig YANG models

Tests [pattern statement](https://tools.ietf.org/html/rfc7950#section-9.4.5)
using a [pyang](https://github.com/mbj4668/pyang) plugin and
[oc-ext:posix-pattern](https://github.com/openconfig/public/blob/master/release/models/openconfig-extensions.yang#L114)
using [goyang](https://github.com/openconfig/goyang).

## How to Contribute to this Repository

Pattern statement tests reside in various modules in the [testdata](testdata)
folder. When adding tests, group pattern tests in test modules named after where
the patterns are found.

### Workflow for Fixing a Bad Pattern in [openconfig/public](https://github.com/openconfig/public)

1.  Add new pattern tests under [testdata](testdata) and post a PR demonstrating
    the failure on openconfig/public.
2.  In the PR, note to the reviewers that the failure is expected due to an
    incorrect existing pattern.
3.  After merge, increment the new minor version of
    openconfig/pattern-regex-tests.
4.  Open a PR in openconfig/public that updates the new version of
    openconfig/pattern-regex-tests in its CI config
    https://github.com/openconfig/public/blob/master/cloudbuild.yaml, while
    making the corresponding pattern fix in the YANG model.

### Limitations

Only typedef tests are currently supported. The current OpenConfig models only
contain patterns in typedef statements.

## Releases

Releases are synchronized with the current OpenConfig YANG models. If new
breaking tests are added (e.g. test cases handled incorrectly by a current
pattern regex), then the minor version must be incremented.

At this time, major version updates are not anticipated, but could occur as a
result of major changes to the repository.

### Demo CI Workflow

There is a demo CI workflow that runs on Pull Requests. They are used to demo
the result of running the tests on the current openconfig/public YANG models.
It's possible test failure are expected (e.g. wrong existing pattern). If this
is the case a minor version increment should be given.

--------------------------------------------------------------------------------

[OpenConfig YANG models](https://github.com/openconfig/public/blob/master/README.md)
