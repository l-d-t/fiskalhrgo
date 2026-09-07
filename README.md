# FiskalHR Go

```text
   ___ _     _         _         __     ___        _ 
  / __(_)___| | ____ _| | /\  /\/__\   / _ \___   / \
 / _\ | / __| |/ / _` | |/ /_/ / \//  / /_\/ _ \ /  /
/ /   | \__ \   < (_| | / __  / _  \ / /_\\ (_) /\_/ 
\/    |_|___/_|\_\__,_|_\/ /_/\/ \_/ \____/\___/\/         
```

[![Test](https://github.com/l-d-t/fiskalhrgo/actions/workflows/test.yml/badge.svg)](https://github.com/l-d-t/fiskalhrgo/actions/workflows/test.yml)
[![CodeQL](https://github.com/l-d-t/fiskalhrgo/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/l-d-t/fiskalhrgo/actions/workflows/github-code-scanning/codeql)
[![Go Report Card](https://goreportcard.com/badge/github.com/l-d-t/fiskalhrgo)](https://goreportcard.com/report/github.com/l-d-t/fiskalhrgo)
![Go version](https://img.shields.io/badge/Go-1.27.1-blue)
[![Go Reference](https://pkg.go.dev/badge/github.com/l-d-t/fiskalhrgo.svg)](https://pkg.go.dev/github.com/l-d-t/fiskalhrgo)

## Overview

FiskalHR Go is an open-source Go library for Croatian CIS fiscalization. It creates and submits invoices, generates ZKI, returns JIR, and bundles CIS certificates. It supports buyer OIB for CIS 1.0 B2B invoices, RSA-SHA256 XML request signatures, and optional native XML signature verification on Linux.

For the current protocol documentation, see the official [CIS technical specification v2.7 (21 July 2026)](https://porezna-uprava.gov.hr/UserDocsImages/Fiskalizacija/Tehni%C4%8Dke%20specifikacije/Fiskalizacija%20-%20Tehnicka%20specifikacija%20za%20korisnike_v2.7%20%2821.07.2026.%29.pdf). B2B support here means the CIS buyer-OIB field, **not a full Fiskalizacija 2.0 / eRačun implementation**.

**Note:** This software is provided **“as is”, without warranty of any kind**, under the [MIT License](LICENSE). The host application remains responsible for invoice calculations, business rules, credential protection, and compliance. If you need a supported version for a mission-critical production system, contact **L. D. T. sales**; see [Commercial and Professional Support](#commercial-and-professional-support) below.

[README in Croatian](README-HR.md)

## Why this project?

While there are numerous open-source implementations of Croatian fiscalization libraries available, they tend to focus on other programming languages. However, after some research, it became clear that an open-source solution in Go (Golang) is hard to come by. To fill this gap, we're developing a pure Go open-source Croatian fiscalization library (package). So here we Go! ;)

## Features

- Create and submit invoices to CIS, including buyer OIB for supported CIS 1.0 B2B payments.
- Sign XML requests with RSA-SHA256 and SHA-256 reference digests.
- Generate ZKI using the unchanged RSA-SHA1 → MD5 calculation.
- Verify CIS response XML signatures when the optional `libxml2` backend is enabled.
- Include updated DEMO and production CIS certificate bundles and validate certificate chains.
- Load client P12 certificates and expose issuer, subject, serial number, and validity details.
- Return structured CIS error codes and messages, including errors sent with HTTP 500.
- Support single-tenant and multi-tenant services, web applications, and desktop applications.
- Provide QR-code data for use with a QR-code generator of your choice.
- Leave invoice calculations and business logic to the host application.

## Go Version Compatibility

Target version: **Go 1.27.1** (current stable release).

## Installation

In your project root, install the module:

```bash
go get github.com/l-d-t/fiskalhrgo
```

### Native XML signatures on Linux

The default backend is pure Go and needs no `libxml2`. The optional Linux backend uses system `libxml2` for canonicalization **both when signing outgoing requests and when verifying incoming CIS XML signatures**. RSA signing and cryptographic verification use Go's crypto libraries.

**Important:** With native mode disabled, CIS response XML signature verification is skipped for legacy compatibility. HTTPS/TLS certificate verification still applies, but it is not a substitute for XML signature verification. Building with the tag alone does not enable native mode: call `SetUseLibxml2(true)` on each entity.

Install the development package and `pkg-config` for your distribution:

```bash
# Debian/Ubuntu
sudo apt install build-essential pkg-config libxml2-dev

# Fedora/RHEL
sudo dnf install gcc pkgconf-pkg-config libxml2-devel

# Alpine
apk add build-base pkgconf libxml2-dev
```

Build with cgo and the `libxml2` tag, then enable the backend on the entity:

```bash
CGO_ENABLED=1 go build -tags libxml2 ./...
```

```go
if err := fiskalEntity.SetUseLibxml2(true); err != nil {
    log.Fatal(err)
}
```

With native mode enabled, signed CIS responses are checked automatically during invoice requests: the verifier checks the signed payload, reference digest, RSA signature, and certificate against the selected embedded CIS certificate. Invalid signatures return an error. `VerifyXML` is also available for explicit verification and returns an error if native mode is disabled. Use `Libxml2Available()` to check build support and `UseLibxml2()` to check the entity setting.

### XML signatures and ZKI are separate

- **Outgoing XML:** RSA-SHA256 signatures and SHA-256 reference digests in both backends. SHA-256 describes the hash, not a 256-bit RSA key.
- **Incoming XML:** The native verifier supports RSA-SHA256 and legacy RSA-SHA1 signatures, with SHA-256 or SHA-1 reference digests.
- **ZKI:** Unchanged: concatenate the prescribed invoice fields, sign their SHA-1 hash with RSA PKCS#1 v1.5, then return the MD5 hash of that signature. The XML SHA-256 change does not alter ZKI.

### CIS certificates

Updated DEMO and production CIS certificate bundles are embedded in the library. For the selected environment, the library validates the bundled chains and selects the newest currently valid certificate. Bundles are not downloaded or refreshed at runtime; keep the library up to date for certificate rotations.

You still need your own valid fiscalization P12 certificate and password. DEMO mode requires a DEMO certificate; the issuer OIB must match it. Never commit private keys, P12 files, or passwords.

## Usage

The example pings CIS, displays certificate information, submits a DEMO invoice, and returns JIR and ZKI. Set `FISKAL_OIB`, `FISKAL_OPERATOR_OIB`, `FISKAL_CERT_PATH`, and `FISKAL_CERT_PASSWORD` in your environment first. Use valid issuer and operator OIBs, not placeholder values.

To verify response XML signatures, build with native support and enable it immediately after creating the entity, as shown above.

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/l-d-t/fiskalhrgo"
)

func main() {
    // Create a new FiskalEntity
    fiskalEntity, err := fiskalhrgo.NewFiskalEntity(
        os.Getenv("FISKAL_OIB"), // issuer OIB matching the certificate
        true,          // sustavPDV
        "Location1",   // locationID, if not DEMO MODE have to be registered
                      // with ePorezna
        true,          // centralized invoice numbers
        true,          // demoMode, if true expected a valid Fiskal demo
                      // certificate
        true,          // chk_expired
        os.Getenv("FISKAL_CERT_PATH"),     // certPath
        os.Getenv("FISKAL_CERT_PASSWORD"), // certPassword
    )
    if err != nil {
        log.Fatalf("Failed to create FiskalEntity: %v", err)
    }

    errPing := fiskalEntity.PingCIS()
    if errPing != nil {
        log.Fatalf("Failed to make Ping request: %v", errPing)
    }

    // Get certificate basic info
    if fiskalEntity.IsExpiringSoon() {
        fmt.Println("WARNING: Certificate is expiring soon")
        fmt.Printf("Certificate expires in %d days",
            fiskalEntity.DaysUntilExpire())
    }

    // Display certificate info
    fmt.Println(fiskalEntity.DisplayCertInfoText())

    invoice, zki, err := fiskalEntity.NewCISInvoice(
        time.Now(),
        uint(1236), // invoice number
        uint(1),    // register id number
        [][]interface{}{ // PDV
            {"25.00", "1000.00", "250.00"},
        },
        nil, // PNP
        nil, // Other taxes
        "0.00", // total amount of exemptions on the issued invoice. 
                // Exemptions in cases where goods are delivered or services are
                // provided that are exempt from VAT payment.
        "0.00", // amount subject to the special margin taxation procedure 
                // if exist
        "0.00", // total amount not subject to taxation on the issued invoice. 
                // This information is submitted to the Tax Administration only
                // if there is an amount on the invoice that is not subject to
                // taxation.
        nil,           // naknade
        "1250.00",     // total
        fiskalhrgo.CISCash, // payment method G - cash, K - credit card, T -
                       // transfer, O - other, C - check (deprecated)
        os.Getenv("FISKAL_OPERATOR_OIB"), // valid operator OIB
    )

    if err != nil {
        log.Fatalf("Failed to create invoice: %v", err)
    }

    // Display the ZKI
    fmt.Println("ZKI: ", zki)

    // Persist the invoice data and ZKI before submission. If submission fails,
    // retain the original data and follow the applicable late-delivery rules
    // and deadlines. Keep the signing certificate serial for future reference:
    // fiskalEntity.GetCertSERIAL().

    // Display the invoice for test
    fmt.Println(invoice)

    // NOW we should have a saved invoice with a valid ZKI and we are ready to
    // send the invoice to the CIS

    // Send test invoice to CIS with InvoiceRequest
    jir, zkiR, err := invoice.InvoiceRequest()

    if err != nil {
        log.Fatalf("Failed to send invoice: %v", err)
    }

    // Display the JIR and ZKI
    fmt.Println("JIR: ", jir)
    fmt.Println("ZKI: ", zkiR)

    // At this point the application can save the JIR with the invoice and
    // commit those changes,

    // Display/send/print the invoice to the user with all elements required by
    // law
}
```

### B2B: buyer OIB in CIS 1.0

After creating the invoice and before calling `InvoiceRequest()`, set the buyer OIB:

```go
if err := invoice.SetOibPrimateljaRacuna(buyerOIB); err != nil {
    log.Fatal(err)
}
```

- `buyerOIB` must be a valid 11-digit OIB with a correct checksum.
- Supported payment methods are `CISCash` (`G`), `CISCard` (`K`), and `CISMixOther` (`O`). Bank transfer (`T`) cannot carry this field in this implementation.
- Pass an empty string to clear the field; ordinary consumer invoices omit it.
- Buyer OIB is not part of the ZKI calculation, so setting it does not regenerate ZKI.
- This sends `OibPrimateljaRacuna` through the CIS invoice protocol; it does not provide eRačun exchange or full Fiskalizacija 2.0 support.

### Scope and limitations

- This package implements the CIS invoice protocol. eRačun exchange and the wider Fiskalizacija 2.0 workflow are outside its scope.
- XML signature verification requires Linux, `libxml2`, a cgo build, and `SetUseLibxml2(true)`. The portable Go backend signs requests but retains the legacy behavior of skipping response XML signature verification.
- CIS certificates are embedded and selected at runtime, but are not downloaded automatically. Update the library when CIS rotates its certificates.

## Commercial and Professional Support

For commercial support, long-term maintenance contracts, and/or consulting or development services, contact [LDT](https://ldt.hr) or send an email to [info@ldt.hr](mailto:info@ldt.hr).

Custom commercial versions for your specific needs in Go or other technologies (C, Zig, Rust, Python, PHP, Mojo, OCaml, Swift, Objective-C) or for specific platforms or embedded devices are also possible.
Full custom (or not) business/invoicing/erp commercial solutions either as a product or SaaS are also possible.

## Sponsors

[![LDT Logo](https://ldt.hr/logo.png)](https://ldt.hr)

If you are interested in sponsoring the project as a company or individual, contact us at [oss@ldt.hr](mailto:oss@ldt.hr).

## Contributors

Contributors are welcome! You can contribute to the development in the following ways:

- Testing
- Reporting issues
- Writing documentation
- Translating documentation
- Submitting pull requests for new features or improving existing ones (recommended to contact and consult before doing the work)

Your contribution is invaluable and helps us create a better product for the community.

## Running tests

The suite includes **live CIS DEMO requests and invoice submissions**, including B2B invoices. It requires network access and a valid DEMO certificate. Set these environment variables, or load them from your private local setup script:

| Variable | Value |
| --- | --- |
| `CIS_P12_BASE64` | Single-line base64-encoded DEMO P12 certificate |
| `FISKALHRGO_TEST_CERT_PASSWORD` | P12 password |
| `FISKALHRGO_TEST_CERT_OIB` | Valid OIB matching the certificate; also used as the DEMO operator OIB |

If you have the local setup script, source it in the same shell before testing:

```bash
source testcert/set_env.sh
go test -count=1 -v ./...

# Linux with libxml2 dependencies installed: also verify CIS XML signatures
CGO_ENABLED=1 go test -count=1 -v -tags libxml2 ./...
```

The test setup automatically enables native verification when the tagged backend is available; applications must enable it explicitly. The default build does not verify response XML signatures.

On Linux, `base64 -w 0` produces the single-line encoding of a P12 file. Base64 is not encryption: keep the certificate, password, and local setup script private and out of version control. In CI, supply these values through secret storage.

HTTP 500 can contain a structured CIS validation error rather than indicate an outage. Inspect the returned CIS code and message; invalid operator OIB test data, for example, can trigger a rejection.
