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
![Go version](https://img.shields.io/badge/Go-1.22%2B-blue)
![GitHub commits since latest release](https://img.shields.io/github/commits-since/l-d-t/fiskalhrgo/latest?include_prereleases)
![GitHub License](https://img.shields.io/github/license/l-d-t/fiskalhrgo)
[![Go Reference](https://pkg.go.dev/badge/github.com/l-d-t/fiskalhrgo.svg)](https://pkg.go.dev/github.com/l-d-t/fiskalhrgo)

# FiskalHR Go

## Pregled

FiskalHR Go je Go modul dizajniran za fiskalizaciju. Cilj ovog modula je pružiti dobro održavano, učinkovito i jednostavno rješenje s ugrađenim certifikatima za provjeru potpisa, čineći je jednostavnom za korištenje.
Temeljeno na ["Fiskalizacija" specifikaciji v2.5](https://www.porezna-uprava.hr/HR_Fiskalizacija/Documents/Fiskalizacija%20-%20Tehnicka%20specifikacija%20za%20korisnike_v2.5._23_10_23pdf.pdf)

**Napomena:** Ovaj projekt je trenutno u tijeku i u alfa fazi te nije potpuno dovršen. Nije preporučeno za korištenje u produkciji, ali uskoro dolazi.

[README in English](README.md)

## Zašto ovaj projekt?

Iako postoji mnogo open-source implementacija libraryja za fiskalizaciju, one su uglavnom usmjerene na druge programske jezike. Međutim, nakon istraživanja postalo je jasno da je teško pronaći open-source rješenje u Go (Golang) jeziku. Kako bismo popunili tu prazninu, razvijamo čisti Go open-source library (paket) za fiskalizaciju. Pa, "let's Go!" ;)

## Značajke

- Obrada i slanje računa CIS-u (Poreznoj upravi) radi usklađenosti sa zakonom.
- Obrada i provjera odgovora od CIS-a.
- Neometajuće za glavnu aplikaciju, ostavljajući poslovnu logiku potpuno glavnoj aplikaciji.
- Informira je li odgovor validiran protiv javnih certifikata CIS-a bez odlučivanja što učiniti s nevalidiranim odgovorima.
- Parsiranje i provjera ugrađenih certifikata.
- Parsiranje i provjera klijentskog P12 certifikata.
- Prikladno za aplikacije s jednim ili više poslovnih subjekta i OIB-a.
- Prikladno za bilo koju vrstu aplikacije (web servis, web aplikacija, desktop aplikacija).
- Ekstrakcija i vraćanje detalja certifikata kao što su javni ključ, izdavatelj, subjekt, serijski broj i razdoblje valjanosti.
- Pomoćne funkcije za generiranje QR koda za ispis na računima u raznim formatima.

### Go Verzija Kompatibilnost
- Minimalna testirana i podržana verzija: **Go 1.22**
- Preporučena verzija: **Go 1.23.1+** za najbolju izvedbu

## Instalacija

U korijenu vašeg projekta preuzmite modul
```
go get github.com/l-d-t/fiskalhrgo
```

## Korištenje

Minimalni jednostavni primjer CIS pinga koristeći EchoRequest i dohvaćanje nekih informacija o certifikatu.

Pošaljite jednostavan minimalni račun, dobijte JIR.

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

## Komercijalna i profesionalna podrška

Za komercijalnu podršku, ugovoru o dugoročnom održavanju, i/ili konzultacije ili usluge razvoja, kontaktirajte [LDT](https://ldt.hr) ili pošaljite email na [info@ldt.hr](mailto:info@ldt.hr).

Moguće su i posebne prilagođene komercijalne verzije za vaše specifične potrebe u Go ili drugim tehnologijama (C, Zig, Rust, Python, PHP, Mojo, OCaml, Swift, Objective-C) ili za specifičnu platformu ili embeded uređaje.
Potpuno custom (ili ne) poslovna/računovodstvena/erp komercijalna rješenja kao proizvod ili SaaS također su moguća.

## Sponzori

<a href="https://ldt.hr" target="_blank">
    <img src="https://ldt.hr/logo.png" alt="LDT Logo">
</a>

Ako ste zainteresirani za sponzorirat projekt kao tvrtka ili individua oss@ldt.hr

## Pridonosioci

Pridonosioci su dobrodošli! Možete doprinijeti razvoju na sljedeće načine:
- Testiranjem
- Prijavom problema
- Pisanjem dokumentacije
- Prijevodom dokumentacije
- Slanjem pull requesta za nove značajke ili poboljšanje postojećih (preporučeno prije kontaktirat i konzultirat se nego se posao napravi)

Vaš doprinos je neprocjenjiv i pomaže nam u stvaranju boljeg proizvoda za zajdnicu.

## Napomena za pokretanje testova

Testove možete pokrenuti s detaljnim ispisom pomoću

```bash
go test -v
```

Prije pokretanja potrebno je postaviti određene varijable okoline.

Varijabla okoline `CIS_P12_BASE64` mora sadržavati jednolinijski base64 enkodirani niz izvorne važeće Fiskalne potvrde u P12 formatu.
Ovaj enkodirani niz je ključan za interakciju testova s CIS-om (Hrvatskim sustavom fiskalizacije).

Za enkodiranje vaše P12 potvrde (npr. `fiskalDemo1.p12`) u jednolinijski base64 niz na Linux sustavu, upotrijebite sljedeću naredbu:

```bash
base64 -w 0 fiskal1.p12
```

Zatim postavite varijablu okoline `CIS_P12_BASE64` s enkodiranim nizom.

Dodatno, provjerite da su varijable okoline `FISKALHRGO_TEST_CERT_PASSWORD` i `FISKALHRGO_TEST_CERT_OIB` postavljene s odgovarajućom lozinkom potvrde i OIB-om (Osobnim identifikacijskim brojem).

Ovaj sustav se koristi za testove jer će se testovi izvoditi u CI (Kontinuiranoj integraciji), gdje se tajne, kao što su one na GitHubu, prenose putem varijabli okoline. Ovo čini upravljanje jednostavnim i praktičnim. Potvrda, lozinka i OIB za testove mogu se lako pohraniti kao GitHub Action tajne, na primjer.
