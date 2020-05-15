# Security Policy

## Supported Versions

This project has finished its initial development cycle. As a rule, the latest version of each major release will be supported.

| Version          | Supported          |
| ---------------- | ------------------ |
| latest (v1.2.1)  | :check_mark:       |
| v1.1.1           | :cross_mark:       |
| v1.0.0           | :cross_mark:       |

## Reporting a Vulnerability

If you find a vulnerability in `go-credentials`, please submit a __Bug report__ or preferably, email security@engi.fyi and the bug will be actioned within 72 hours. Alternately, feel free to submit a PR with code changes to resolve the issue and we will endeavour to merge it within 24 hours.

## Other

We have security scanning enabled on this repository using https://github.com/securego/gosec. As a note, we have rule G101 disabled as we have hardcoded strings for test vars.
