package fiskalhrgo

import "testing"

func TestCISResponseError(t *testing.T) {
	const status = "500 Internal Server Error"
	for _, tt := range []struct {
		name string
		body string
		want string
	}{
		{
			name: "namespaced CIS errors",
			body: `<tns:RacunOdgovor xmlns:tns="http://www.apis-it.hr/fin/2012/types/f73"><tns:Greske><tns:Greska><tns:SifraGreske>s013</tns:SifraGreske><tns:PorukaGreske>Račun sadrži restriktivne greške: 179</tns:PorukaGreske></tns:Greska><tns:Greska><tns:SifraGreske>s004</tns:SifraGreske><tns:PorukaGreske>Neispravan potpis.</tns:PorukaGreske></tns:Greska></tns:Greske></tns:RacunOdgovor>`,
			want: "CIS returned an error: 500 Internal Server Error: s013: Račun sadrži restriktivne greške: 179; s004: Neispravan potpis.",
		},
		{
			name: "SOAP fault",
			body: `<soap:Fault xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><faultcode>soap:Server</faultcode><faultstring>Service unavailable</faultstring></soap:Fault>`,
			want: "CIS returned an error: 500 Internal Server Error: soap:Server: Service unavailable",
		},
		{
			name: "missing error details",
			body: `<RacunOdgovor/>`,
			want: "CIS returned an error: 500 Internal Server Error",
		},
		{
			name: "malformed XML",
			body: `<broken`,
			want: "CIS returned an error: 500 Internal Server Error",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := cisResponseError(status, []byte(tt.body)).Error(); got != tt.want {
				t.Fatalf("got %q, want %q", got, tt.want)
			}
		})
	}
}