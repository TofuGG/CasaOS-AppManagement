# CasaOS-AppManagement

> ⚠️ **UNOFFICIAL FORK — NOT THE OFFICIAL RELEASE.** This is a **TofuGG** community
> fork of the CasaOS AppManagement service, not affiliated with or endorsed by the
> official CasaOS / IceWhaleTech team. It is provided **AS-IS** with **no warranty**
> and **no official support**. **USE AT YOUR OWN RISK.**

[![Go Reference](https://pkg.go.dev/badge/github.com/IceWhaleTech/CasaOS-AppManagement.svg)](https://pkg.go.dev/github.com/IceWhaleTech/CasaOS-AppManagement)
[![Go Report Card](https://goreportcard.com/badge/github.com/IceWhaleTech/CasaOS-AppManagement)](https://goreportcard.com/report/github.com/IceWhaleTech/CasaOS-AppManagement)
[![goreleaser](https://github.com/IceWhaleTech/CasaOS-AppManagement/actions/workflows/release.yml/badge.svg)](https://github.com/IceWhaleTech/CasaOS-AppManagement/actions/workflows/release.yml)
[![codecov](https://codecov.io/gh/IceWhaleTech/CasaOS-AppManagement/branch/main/graph/badge.svg?token=ZCWZOFKXJT)](https://codecov.io/gh/IceWhaleTech/CasaOS-AppManagement)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=IceWhaleTech_CasaOS-AppManagement&metric=vulnerabilities)](https://sonarcloud.io/summary/new_code?id=IceWhaleTech_CasaOS-AppManagement)
[![Bugs](https://sonarcloud.io/api/project_badges/measure?project=IceWhaleTech_CasaOS-AppManagement&metric=bugs)](https://sonarcloud.io/summary/new_code?id=IceWhaleTech_CasaOS-AppManagement)
[![Code Smells](https://sonarcloud.io/api/project_badges/measure?project=IceWhaleTech_CasaOS-AppManagement&metric=code_smells)](https://sonarcloud.io/summary/new_code?id=IceWhaleTech_CasaOS-AppManagement)
[![Lines of Code](https://sonarcloud.io/api/project_badges/measure?project=IceWhaleTech_CasaOS-AppManagement&metric=ncloc)](https://sonarcloud.io/summary/new_code?id=IceWhaleTech_CasaOS-AppManagement)
[![Duplicated Lines (%)](https://sonarcloud.io/api/project_badges/measure?project=IceWhaleTech_CasaOS-AppManagement&metric=duplicated_lines_density)](https://sonarcloud.io/summary/new_code?id=IceWhaleTech_CasaOS-AppManagement)

App management service manages CasaOS apps lifecycle, such as installation, running, etc.

## Security Hardening

This codebase has been hardened against the following attack classes:

| # | Fix | Files |
|---|-----|-------|
| 7 | **JWT skippers removed** — both v1 and v2 routes require a valid bearer token; no localhost bypass | `route/v1.go`, `route/v2.go` |
| 8 | **SSRF guard** — `go-getter` restricted to `https` scheme only; app-store URLs validated against public-IP resolution + DNS denylist | `pkg/utils/downloadHelper/getter.go`, `service/appstore.go` |
| 9 | **TLS verification restored** — `InsecureSkipVerify: true` removed from Docker digest HTTP client | `pkg/docker/digest.go` |
| 10 | **Default password replaced** — host-unique random secret generated at first start, stored with mode `0600` | `pkg/utils/envHelper/env.go` |

Build target for verification: `GOOS=linux GOARCH=amd64 go build ./...`
