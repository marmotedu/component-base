# component-base

Scheme, typing, encoding, decoding, and conversion packages for IAM and IAM-like API objects.

## Purpose

This library is a shared dependency for servers and clients to work with IAM API infrastructure without direct
type dependencies. Its first consumers are `github.com/marmotedu/component-base`, `github.com/marmotedu/marmotedu-sdk-go`.


## Compatibility

There are *NO compatibility guarantees* for this repository. It is in direct support of IAM, so branches
will track IAM and be compatible with that repo. As we more cleanly separate the layers, we will review the
compatibility guarantee.
