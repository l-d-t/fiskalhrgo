```
   ___ _     _         _         __     ___        _ 
  / __(_)___| | ____ _| | /\  /\/__\   / _ \___   / \
 / _\ | / __| |/ / _` | |/ /_/ / \//  / /_\/ _ \ /  /
/ /   | \__ \   < (_| | / __  / _  \ / /_\\ (_) /\_/ 
\/    |_|___/_|\_\__,_|_\/ /_/\/ \_/ \____/\___/\/                                     
```

# FiskalHR Go

## Pregled

FiskalHR Go je Go modul dizajniran za fiskalizaciju. Cilj ovog modula je pružiti dobro održavano, učinkovito i jednostavno rješenje s ugrađenim certifikatima za provjeru potpisa, čineći je jednostavnom za korištenje.
Temeljeno na ["Fiskalizacija" specifikaciji v2.5](https://www.porezna-uprava.hr/HR_Fiskalizacija/Documents/Fiskalizacija%20-%20Tehnicka%20specifikacija%20za%20korisnike_v2.5._23_10_23pdf.pdf)

**Napomena:** Ovaj projekt je trenutno u izradi i još nije upotrebljiv.

[README in English](README.md)

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

## Instalacija

Ovaj projekt još nije dostupan za instalaciju jer je još uvijek u razvoju.

## Korištenje

Vodič i dokumentacija dolaze uskoro, kada prestanu promjene koje narušavaju kompatibilnost i biblioteka bude kompletna i stabilna.

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