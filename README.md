```
   ___ _     _         _         __     ___        _ 
  / __(_)___| | ____ _| | /\  /\/__\   / _ \___   / \
 / _\ | / __| |/ / _` | |/ /_/ / \//  / /_\/ _ \ /  /
/ /   | \__ \   < (_| | / __  / _  \ / /_\\ (_) /\_/ 
\/    |_|___/_|\_\__,_|_\/ /_/\/ \_/ \____/\___/\/         
```

# FiskalHR Go

## Overview

FiskalHR Go is a Go library designed for Croatian fiscalization (fiskalizacija). The goal of this library is to provide a well-maintained, efficient, and hassle-free solution with embedded certificates for checking signatures, making it easy to use.
Based on ["Fiskalizacija" specification v2.5](https://www.porezna-uprava.hr/HR_Fiskalizacija/Documents/Fiskalizacija%20-%20Tehnicka%20specifikacija%20za%20korisnike_v2.5._23_10_23pdf.pdf)

**Note:** This project is currently a work in progress and is not yet usable.

[README in Croatian](README-HR.md)

## Features

- Process and send invoices to CIS (Croatian Tax Administration) for compliance the law.
- Handle and verify responses from CIS.
- Non-intrusive to the host application, leaving business logic entirely to the host.
- Inform whether a response is validated against CIS public certificates without deciding what to do with unverified responses.
- Parse and verify embedded certificates.
- Parse and verify client P12 certificate.
- Suitable for single tenant and multitenant application
- Suitable for any type of application (web service, web app, desktop)
- Extract and return certificate details such as public key, issuer, subject, serial number, and validity period.
- Helper functions for generating QR code for printing on invoices in various formats

## Installation

This project is not yet available for installation as it is still under development.

## Usage

Guide and documentation coming soon, when breaking changes stop and the library is feature complete

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