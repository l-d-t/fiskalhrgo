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

## Pregled

FiskalHR Go je biblioteka otvorenog koda za fiskalizaciju putem hrvatskog CIS-a. Omogućuje izradu i slanje računa, generiranje ZKI-ja i dohvat JIR-a te sadrži ugrađene CIS certifikate. Podržava OIB primatelja za B2B račune u CIS-u 1.0, RSA-SHA256 potpise XML zahtjeva i opcionalnu izvornu provjeru XML potpisa na Linuxu.

Aktualna dokumentacija protokola dostupna je u službenoj [tehničkoj specifikaciji CIS-a v2.7 (21. srpnja 2026.)](https://porezna-uprava.gov.hr/UserDocsImages/Fiskalizacija/Tehni%C4%8Dke%20specifikacije/Fiskalizacija%20-%20Tehnicka%20specifikacija%20za%20korisnike_v2.7%20%2821.07.2026.%29.pdf). B2B podrška ovdje označava polje OIB-a primatelja u CIS-u, **a ne cjelovitu implementaciju Fiskalizacije 2.0 / eRačuna**.

**Napomena:** Ovaj softver pruža se **„kakav jest” (as is), bez ikakvog jamstva**, u skladu s [MIT licencom](LICENSE). Aplikacija koja koristi biblioteku odgovorna je za izračune računa, poslovna pravila, zaštitu vjerodajnica i usklađenost s propisima. Ako trebate verziju s podrškom za produkcijski sustav kritičan za poslovanje, obratite se **prodaji L. D. T.-a**; pogledajte odjeljak [Komercijalna i profesionalna podrška](#komercijalna-i-profesionalna-podrška) u nastavku.

[README in English](README.md)

## Zašto ovaj projekt?

Postoje brojne biblioteke otvorenog koda za fiskalizaciju, ali većina je namijenjena drugim programskim jezicima. Ovaj projekt donosi rješenje za Go (Golang), s prenosivom Go implementacijom i opcionalnom podrškom za `libxml2` na Linuxu.

## Značajke

- Izrada i slanje računa CIS-u, uključujući OIB primatelja za podržane načine plaćanja B2B računa u CIS-u 1.0.
- Potpisivanje XML zahtjeva algoritmom RSA-SHA256 uz SHA-256 sažetke referenci.
- Generiranje ZKI-ja nepromijenjenim postupkom RSA-SHA1 → MD5.
- Provjera XML potpisa CIS odgovora kada je uključena opcionalna podrška za `libxml2`.
- Ažurirani ugrađeni DEMO i produkcijski CIS certifikati uz provjeru certifikacijskih lanaca.
- Učitavanje klijentskih P12 certifikata i dohvat podataka o izdavatelju, subjektu, serijskom broju i valjanosti.
- Vraćanje strukturiranih CIS šifri i poruka grešaka, uključujući odgovore s HTTP statusom 500.
- Podrška za servise, web i desktop aplikacije s jednim ili više poslovnih subjekata.
- Priprema podataka za QR kod, uz generator QR koda po vašem izboru.
- Izračuni računa i poslovna logika ostaju odgovornost aplikacije koja koristi biblioteku.

## Kompatibilnost s verzijama Go-a

Ciljana verzija: **Go 1.27.1** (trenutačna stabilna verzija).

## Instalacija

U korijenu svojeg projekta instalirajte modul:

```bash
go get github.com/l-d-t/fiskalhrgo
```

### Izvorna podrška za XML potpise na Linuxu

Zadana implementacija koristi samo Go i ne zahtijeva `libxml2`. Opcionalna Linux implementacija koristi sistemski `libxml2` za kanonikalizaciju **pri potpisivanju izlaznih zahtjeva i pri provjeri XML potpisa ulaznih CIS odgovora**. RSA potpisivanje i kriptografska provjera koriste Go kriptografske biblioteke.

**Važno:** Kada izvorna podrška nije uključena, provjera XML potpisa CIS odgovora preskače se radi kompatibilnosti s prethodnim ponašanjem. Provjera HTTPS/TLS certifikata i dalje se provodi, ali nije zamjena za provjeru XML potpisa. Sama kompilacija s oznakom nije dovoljna: pozovite `SetUseLibxml2(true)` na svakom entitetu.

Instalirajte razvojni paket i `pkg-config` za svoju distribuciju:

```bash
# Debian/Ubuntu
sudo apt install build-essential pkg-config libxml2-dev

# Fedora/RHEL
sudo dnf install gcc pkgconf-pkg-config libxml2-devel

# Alpine
apk add build-base pkgconf libxml2-dev
```

Kompilirajte s podrškom za cgo i oznakom `libxml2`, a zatim uključite izvornu podršku na entitetu:

```bash
CGO_ENABLED=1 go build -tags libxml2 ./...
```

```go
if err := fiskalEntity.SetUseLibxml2(true); err != nil {
    log.Fatal(err)
}
```

Kada je izvorna podrška uključena, potpisi CIS odgovora automatski se provjeravaju pri slanju računa: provjeravaju se potpisani sadržaj, sažetak reference, RSA potpis i certifikat u odnosu na odabrani ugrađeni CIS certifikat. Neispravan potpis vraća grešku. Za izričitu provjeru dostupan je i `VerifyXML`, koji vraća grešku ako izvorna podrška nije uključena. `Libxml2Available()` provjerava je li podrška ugrađena pri kompilaciji, a `UseLibxml2()` je li uključena na entitetu.

### XML potpisi i ZKI zasebni su postupci

- **Izlazni XML:** RSA-SHA256 potpisi i SHA-256 sažeci referenci u obje implementacije. SHA-256 označava sažetak, a ne 256-bitni RSA ključ.
- **Ulazni XML:** Izvorna provjera podržava RSA-SHA256 i naslijeđene RSA-SHA1 potpise, uz SHA-256 ili SHA-1 sažetke referenci.
- **ZKI:** Nepromijenjen: propisana polja računa spajaju se, njihov SHA-1 sažetak potpisuje se postupkom RSA PKCS#1 v1.5, a zatim se vraća MD5 sažetak tog potpisa. Prelazak XML potpisa na SHA-256 ne mijenja ZKI.

### CIS certifikati

Ažurirani DEMO i produkcijski CIS certifikati ugrađeni su u biblioteku. Za odabrano okruženje biblioteka provjerava ugrađene certifikacijske lance i odabire najnoviji trenutačno važeći certifikat. Certifikati se ne preuzimaju niti osvježavaju tijekom rada; ažurirajte biblioteku radi zamjene CIS certifikata.

I dalje trebate vlastiti važeći fiskalni P12 certifikat i lozinku. DEMO način rada zahtijeva DEMO certifikat, a OIB izdavatelja mora mu odgovarati. Privatne ključeve, P12 datoteke i lozinke nikada nemojte spremati u repozitorij.

## Korištenje

Primjer provjerava vezu s CIS-om, prikazuje podatke o certifikatu, šalje DEMO račun i vraća JIR i ZKI. Prethodno postavite varijable okoline `FISKAL_OIB`, `FISKAL_OPERATOR_OIB`, `FISKAL_CERT_PATH` i `FISKAL_CERT_PASSWORD`. Koristite važeće OIB-e izdavatelja i operatera, a ne ogledne nevažeće vrijednosti.

Za provjeru XML potpisa odgovora kompilirajte s izvornom podrškom i uključite je odmah nakon izrade entiteta, kao u prethodnom primjeru.

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

### B2B: OIB primatelja računa u CIS-u 1.0

Nakon izrade računa, a prije poziva `InvoiceRequest()`, postavite OIB primatelja:

```go
if err := invoice.SetOibPrimateljaRacuna(buyerOIB); err != nil {
    log.Fatal(err)
}
```

- `buyerOIB` mora biti važeći 11-znamenkasti OIB s ispravnom kontrolnom znamenkom.
- Podržani načini plaćanja su `CISCash` (`G`), `CISCard` (`K`) i `CISMixOther` (`O`). Plaćanje na transakcijski račun (`T`) u ovoj implementaciji ne može sadržavati to polje.
- Prazan niz uklanja polje; obični računi krajnjim potrošačima izostavljaju ga.
- OIB primatelja nije dio izračuna ZKI-ja pa njegovo postavljanje ne generira novi ZKI.
- Polje `OibPrimateljaRacuna` šalje se kroz CIS protokol za račune; ovo ne obuhvaća razmjenu eRačuna niti potpunu podršku za Fiskalizaciju 2.0.

### Opseg i ograničenja

- Paket implementira CIS protokol za račune. Razmjena eRačuna i širi tijek Fiskalizacije 2.0 nisu dio njegova opsega.
- Provjera XML potpisa zahtijeva Linux, `libxml2`, cgo kompilaciju i `SetUseLibxml2(true)`. Prenosiva Go implementacija potpisuje zahtjeve, ali zadržava naslijeđeno ponašanje preskakanja provjere XML potpisa odgovora.
- CIS certifikati ugrađeni su i odabiru se pri radu, ali se ne preuzimaju automatski. Ažurirajte biblioteku kada CIS promijeni certifikate.

## Komercijalna i profesionalna podrška

Za komercijalnu podršku, ugovore o dugoročnom održavanju, savjetovanje ili usluge razvoja kontaktirajte [LDT](https://ldt.hr) ili pošaljite e-poruku na [info@ldt.hr](mailto:info@ldt.hr).

Moguće su prilagođene komercijalne verzije u Go-u ili drugim tehnologijama (C, Zig, Rust, Python, PHP, Mojo, OCaml, Swift, Objective-C), kao i rješenja za posebne platforme ili ugrađene uređaje.
Dostupna su i poslovna, računovodstvena i ERP rješenja po mjeri, kao proizvod ili SaaS.

## Sponzori

[![LDT Logo](https://ldt.hr/logo.png)](https://ldt.hr)

Ako želite sponzorirati projekt kao tvrtka ili pojedinac, javite se na [oss@ldt.hr](mailto:oss@ldt.hr).

## Doprinos projektu

Suradnici su dobrodošli! Razvoju možete doprinijeti na sljedeće načine:

- Testiranjem
- Prijavom problema
- Pisanjem dokumentacije
- Prijevodom dokumentacije
- Slanjem pull requestova za nove značajke ili poboljšanja postojećih (preporučujemo prethodni dogovor prije početka rada)

Vaš doprinos pomaže nam u stvaranju boljeg proizvoda za zajednicu.

## Pokretanje testova

Testovi uključuju **stvarne zahtjeve CIS DEMO servisu i slanje računa**, uključujući B2B račune. Potrebni su mrežni pristup i važeći DEMO certifikat. Postavite ove varijable okoline ili ih učitajte iz privatne lokalne skripte:

| Varijabla | Vrijednost |
| --- | --- |
| `CIS_P12_BASE64` | DEMO P12 certifikat kodiran kao jednolinijski base64 niz |
| `FISKALHRGO_TEST_CERT_PASSWORD` | Lozinka P12 certifikata |
| `FISKALHRGO_TEST_CERT_OIB` | Važeći OIB koji odgovara certifikatu; koristi se i kao OIB operatera u DEMO testovima |

Ako imate lokalnu skriptu za postavljanje okruženja, učitajte je u istoj ljusci prije testiranja:

```bash
source testcert/set_env.sh
go test -count=1 -v ./...

# Linux s instaliranim libxml2 ovisnostima: provjerava i XML potpise CIS odgovora
CGO_ENABLED=1 go test -count=1 -v -tags libxml2 ./...
```

Testno okruženje automatski uključuje izvornu provjeru kada je podrška dostupna u kompiliranoj verziji; aplikacije je moraju uključiti izričito. Zadana Go implementacija ne provjerava XML potpise odgovora.

Na Linuxu `base64 -w 0` stvara jednolinijski zapis P12 datoteke. Base64 nije enkripcija: certifikat, lozinku i lokalnu skriptu čuvajte privatno, izvan sustava za upravljanje verzijama. U CI okruženju koristite pohranu tajni.

HTTP 500 može sadržavati strukturiranu CIS grešku validacije, a ne nužno značiti nedostupnost servisa. Provjerite vraćenu šifru i poruku greške; primjerice, nevažeći OIB operatera u testnim podacima može uzrokovati odbijanje računa.
