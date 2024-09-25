```
   ___ _     _         _         __     ___        _ 
  / __(_)___| | ____ _| | /\  /\/__\   / _ \___   / \
 / _\ | / __| |/ / _` | |/ /_/ / \//  / /_\/ _ \ /  /
/ /   | \__ \   < (_| | / __  / _  \ / /_\\ (_) /\_/ 
\/    |_|___/_|\_\__,_|_\/ /_/\/ \_/ \____/\___/\/         
```
[![Test](https://github.com/l-d-t/fiskalhrgo/actions/workflows/test.yml/badge.svg)](https://github.com/l-d-t/fiskalhrgo/actions/workflows/test.yml)
[![CodeQL](https://github.com/l-d-t/fiskalhrgo/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/l-d-t/fiskalhrgo/actions/workflows/github-code-scanning/codeql)
[![Go Report Card](https://goreportcard.com/badge/github.com/l-d-t/fiskalhrgo)](https://goreportcard.com/report/github.com/l-d-t/fiskalhrgo)
![Go version](https://img.shields.io/badge/Go%20version-1.22%2B-blue)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/l-d-t/fiskalhrgo/latest?include_prereleases)
![GitHub License](https://img.shields.io/github/license/l-d-t/fiskalhrgo)

# FiskalHR Go

## Overview

FiskalHR Go is a Go library designed for Croatian fiscalization (fiskalizacija). The goal of this library is to provide a well-maintained, efficient, and hassle-free solution with embedded certificates for checking signatures, making it easy to use.
Based on ["Fiskalizacija" specification v2.5](https://www.porezna-uprava.hr/HR_Fiskalizacija/Documents/Fiskalizacija%20-%20Tehnicka%20specifikacija%20za%20korisnike_v2.5._23_10_23pdf.pdf)

**Note:** This project is currently a work in progress and alpha stage and not feature complete, not recommended in production but coming soon.

[README in Croatian](README-HR.md)

## Why this project?

While there are numerous open-source implementations of Croatian fiscalization libraries available, they tend to focus on other programming languages. However, after some research, it became clear that an open-source solution in Go (Golang) is hard to come by. To fill this gap, we're developing a pure Go open-source Croatian fiscalization library (package). So here we Go! ;)

## Features

- Process and send invoices to CIS (Croatian Tax Administration) for compliance the law.
- Handle and verify responses from CIS.
- Non-intrusive to the host application, leaving business logic entirely to the host.
- Parse and verify embedded certificates.
- Parse and verify client P12 certificate.
- Suitable for single tenant and multitenant application
- Suitable for any type of application (web service, web app, desktop)
- Extract and return certificate details such as public key, issuer, subject, serial number, and validity period.
- Helper function to get data for QR code (that can be passed to a QR code generator of your choice)

## Go Version Compatibility
- Minimum tested and supported version: **Go 1.22**
- Recommended version: **Go 1.23.1+** for best performance

## Installation

In your project root get the module
```
go get github.com/l-d-t/fiskalhrgo
```

## Usage

Minimal simple example of CIS ping using the EchoRequest and get some cert info.

Send a simple minimal invoice, get the JIR

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/l-d-t/fiskalhrgo"
)

func main() {
    // Create a new FiskalEntity
    fiskalEntity, err := NewFiskalEntity(
        "12345678901", // OIB
        true,          // sustavPDV
        "Location1",   // locationID, if not DEMO MODE have to be registered
                      // with ePorezna
        true,          // centralized invoice numbers
        true,          // demoMode, if true expected a valid Fiskal demo
                      // certificate
        true,          // chk_expired
        "cert.p12",    // certPath
        "password",    // certPassword
    )
    if err != nil {
        log.Fatalf("Failed to create FiskalEntity: %v", err)
    }

    errPing := fiskalEntity.PingCIS()
    if errPing != nil {
        log.Fatalf("Failed to make Ping request: %v", err)
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
        "G",           // payment method G - cash, K - credit card, T -
                       // transfer, O - other, C - check (deprecated)
        "12345678901", // operator OIB
        false,         // late delivery, if previous attempt failed but the
                       // invoice was issued with just ZKI
        "",            // receipt book number, if the invoicing system was
                       // unusable and the invoice was issued manually, the
                       // number of the receipt book
        "",            // unused, reserved field for future or temporary
                       // unexpected use by the CIS, should be empty
    )

    if err != nil {
        log.Fatalf("Failed to create invoice: %v", err)
    }

    // Display the ZKI
    fmt.Println("ZKI: ", zki)

    // At this point, the application can save the ZKI with the invoice and
    // commit those changes. An invoice is valid with just ZKI and can be
    // issued even if the next step fails and there is no JIR. But in that case,
    // the process has to be repeated with the same identical data within 48h
    // with the flag for late delivery set to true. A valid JIR has to be added
    // within 48h to an invoice with just ZKI. It is recommended to save the
    // serial of the certificate used to generate it for future reference. You
    // can get the cert serial with fiskalEntity.GetCertSERIAL().

    // Display the invoice
    fmt.Println(invoice)

    // NOW we should have a saved invoice with a valid ZKI and we are ready to
    // send the invoice to the CIS

    // Send test invoice to CIS with InvoiceRequest
    jir, zkiR, err := fiskalEntity.InvoiceRequest(invoice)

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

## Commercial and Professional Support

For commercial support, long-term maintenance contracts, and/or consulting or development services, contact [LDT](https://ldt.hr) or send an email to [info@ldt.hr](mailto:info@ldt.hr).

Custom commercial versions for your specific needs in Go or other technologies (C, Zig, Rust, Python, PHP, Mojo, OCaml, Swift, Objective-C) or for specific platforms or embedded devices are also possible.
Full custom (or not) business/invoicing/erp commercial solutions either as a product or SaaS are also possible.

## Sponsors

<a href="https://ldt.hr" target="_blank">
    <img src="https://ldt.hr/logo.png" alt="LDT Logo">
</a>

If you are interested in sponsoring the project as a company or individual, contact us at [oss@ldt.hr](mailto:oss@ldt.hr).

## Contributors

Contributors are welcome! You can contribute to the development in the following ways:
- Testing
- Reporting issues
- Writing documentation
- Translating documentation
- Submitting pull requests for new features or improving existing ones (recommended to contact and consult before doing the work)

Your contribution is invaluable and helps us create a better product for the community.

## Note for Running Tests

You can run tests with verbose output with

```bash
go test -v
```

Before running some environment variables must be se

The `CIS_P12_BASE64` environment variable must contain a single-line base64 encoded string of the original valid Fiskal certificate in P12 format.
This encoded string is essential for the tests to  interact with the CIS (Croatian Fiscalization System).

To encode your P12 certificate file (e.g., `fiskalDemo1.p12`) to a single-line base64 string on a Linux system, use the following command:

```bash
base64 -w 0 fiskal1.p12
```

Then, set the CIS_P12_BASE64 environment variable with the encoded string.

Additionally, ensure that the `FISKALHRGO_TEST_CERT_PASSWORD` and `FISKALHRGO_TEST_CERT_OIB` environment variables are set with the appropriate certificate password and OIB (Personal Identification Number) respectively.

This system is used for the tests because these tests will run in CI (Continuous Integration), so secrets, for example on GitHub, are passed as environment variables. This makes it easy and convenient to manage. The certificate, password, and OIB for tests can be easily stored as GitHub Action secrets, for example.
