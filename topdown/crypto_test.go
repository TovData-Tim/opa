// Copyright 2018 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package topdown

import (
	"fmt"
	"testing"
)

func TestCryptoX509ParseCertificates(t *testing.T) {

	rule := `
		p = x {
			parsed := crypto.x509.parse_certificates(certs)
			x := [ x | x := parsed[_].Subject.CommonName ]
		}
	`

	tests := []struct {
		note     string
		certs    string
		rule     string
		expected interface{}
	}{
		{
			note:     "DER, single cert, b64",
			certs:    `MIIDujCCAqKgAwIBAgIIE31FZVaPXTUwDQYJKoZIhvcNAQEFBQAwSTELMAkGA1UEBhMCVVMxEzARBgNVBAoTCkdvb2dsZSBJbmMxJTAjBgNVBAMTHEdvb2dsZSBJbnRlcm5ldCBBdXRob3JpdHkgRzIwHhcNMTQwMTI5MTMyNzQzWhcNMTQwNTI5MDAwMDAwWjBpMQswCQYDVQQGEwJVUzETMBEGA1UECAwKQ2FsaWZvcm5pYTEWMBQGA1UEBwwNTW91bnRhaW4gVmlldzETMBEGA1UECgwKR29vZ2xlIEluYzEYMBYGA1UEAwwPbWFpbC5nb29nbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfRrObuSW5T7q5CnSEqefEmtH4CCv6+5EckuriNr1CjfVvqzwfAhopXkLrq45EQm8vkmf7W96XJhC7ZM0dYi1/qOCAU8wggFLMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAaBgNVHREEEzARgg9tYWlsLmdvb2dsZS5jb20wCwYDVR0PBAQDAgeAMGgGCCsGAQUFBwEBBFwwWjArBggrBgEFBQcwAoYfaHR0cDovL3BraS5nb29nbGUuY29tL0dJQUcyLmNydDArBggrBgEFBQcwAYYfaHR0cDovL2NsaWVudHMxLmdvb2dsZS5jb20vb2NzcDAdBgNVHQ4EFgQUiJxtimAuTfwb+aUtBn5UYKreKvMwDAYDVR0TAQH/BAIwADAfBgNVHSMEGDAWgBRK3QYWG7z2aLV29YG2u2IaulqBLzAXBgNVHSAEEDAOMAwGCisGAQQB1nkCBQEwMAYDVR0fBCkwJzAloCOgIYYfaHR0cDovL3BraS5nb29nbGUuY29tL0dJQUcyLmNybDANBgkqhkiG9w0BAQUFAAOCAQEAH6RYHxHdcGpMpFE3oxDoFnP+gtuBCHan2yE2GRbJ2Cw8Lw0MmuKqHlf9RSeYfd3BXeKkj1qO6TVKwCh+0HdZk283TZZyzmEOyclm3UGFYe82P/iDFt+CeQ3NpmBg+GoaVCuWAARJN/KfglbLyyYygcQq0SgeDh8dRKUiaW3HQSoYvTvdTuqzwK4CXsr3b5/dAOY8uMuG/IAR3FgwTbZ1dtoWRvOTa8hYiU6A475WuZKyEHcwnGYe57u2I2KbMgcKjPniocj4QzgYsVAVKW3IwaOhyE+vPxsiUkvQHdO2fojCkY8jg70jxM+gu59tPDNbw3Uh/2Ij310FgTHsnGQMyA==`,
			rule:     rule,
			expected: `["mail.google.com"]`,
		},
		{
			note:     "DER, chain, b64",
			certs:    `MIIDIjCCAougAwIBAgIQbt8NlJn9RTPdEpf8Qqk74TANBgkqhkiG9w0BAQUFADBMMQswCQYDVQQGEwJaQTElMCMGA1UEChMcVGhhd3RlIENvbnN1bHRpbmcgKFB0eSkgTHRkLjEWMBQGA1UEAxMNVGhhd3RlIFNHQyBDQTAeFw0wOTAzMjUxNjQ5MjlaFw0xMDAzMjUxNjQ5MjlaMGkxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1Nb3VudGFpbiBWaWV3MRMwEQYDVQQKEwpHb29nbGUgSW5jMRgwFgYDVQQDEw9tYWlsLmdvb2dsZS5jb20wgZ8wDQYJKoZIhvcNAQEBBQADgY0AMIGJAoGBAMXW+JL8yvVhSwZBSegKLJWBohjvQew1vXpYElrnb56lTdyJOrvrAp9rc2Fr8P/YaHkfunr5xK6/Nwa6Puru0nQ1tN3PsVfAXzUdZqqH/uDeBy1m13Ov+9Nqt4vvCQ4MyGGpA6yQ3Zi1HJxBVmwBfwvuw7/zkQUf+6D1zGhQrSpZAgMBAAGjgecwgeQwKAYDVR0lBCEwHwYIKwYBBQUHAwEGCCsGAQUFBwMCBglghkgBhvhCBAEwNgYDVR0fBC8wLTAroCmgJ4YlaHR0cDovL2NybC50aGF3dGUuY29tL1RoYXd0ZVNHQ0NBLmNybDByBggrBgEFBQcBAQRmMGQwIgYIKwYBBQUHMAGGFmh0dHA6Ly9vY3NwLnRoYXd0ZS5jb20wPgYIKwYBBQUHMAKGMmh0dHA6Ly93d3cudGhhd3RlLmNvbS9yZXBvc2l0b3J5L1RoYXd0ZV9TR0NfQ0EuY3J0MAwGA1UdEwEB/wQCMAAwDQYJKoZIhvcNAQEFBQADgYEAYvHzBQ68EF5JfHrt+H4k0vSphrs7g3vRm5HrytmLBlmS9r0rSbfW08suQnqZ1gbHsdRjUlJ/rDnmqLZybeW/cCEqUsugdjSl4zIBG9GGjnjrXjyTzwMHInZ4byB0lP6qDtnVOyEQp2Vx+QIJza6IQ4XIglhwMO4V8z12Hi5FprwwggMjMIICjKADAgECAgQwAAACMA0GCSqGSIb3DQEBBQUAMF8xCzAJBgNVBAYTAlVTMRcwFQYDVQQKEw5WZXJpU2lnbiwgSW5jLjE3MDUGA1UECxMuQ2xhc3MgMyBQdWJsaWMgUHJpbWFyeSBDZXJ0aWZpY2F0aW9uIEF1dGhvcml0eTAeFw0wNDA1MTMwMDAwMDBaFw0xNDA1MTIyMzU5NTlaMEwxCzAJBgNVBAYTAlpBMSUwIwYDVQQKExxUaGF3dGUgQ29uc3VsdGluZyAoUHR5KSBMdGQuMRYwFAYDVQQDEw1UaGF3dGUgU0dDIENBMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDU02fQjRV/rs0x/n0dkaE/C3E8rMzIZPtj/DJLB5S9b4C6L+EEk8Az/AkzI+kLdCtxxAPG0s3iL/UJY83/SKUAv+Dn84i3LTLemDbmCq0Ae8RkSjuEdQPycJJ9DmL1IatpNoQxdZD4v8dsiBsGlXzJ5ajedaEsemjf1coch1hgGQIDAQABo4H+MIH7MBIGA1UdEwEB/wQIMAYBAf8CAQAwCwYDVR0PBAQDAgEGMBEGCWCGSAGG+EIBAQQEAwIBBjAoBgNVHREEITAfpB0wGzEZMBcGA1UEAxMQUHJpdmF0ZUxhYmVsMy0xNTAxBgNVHR8EKjAoMCagJKAihiBodHRwOi8vY3JsLnZlcmlzaWduLmNvbS9wY2EzLmNybDAyBggrBgEFBQcBAQQmMCQwIgYIKwYBBQUHMAGGFmh0dHA6Ly9vY3NwLnRoYXd0ZS5jb20wNAYDVR0lBC0wKwYIKwYBBQUHAwEGCCsGAQUFBwMCBglghkgBhvhCBAEGCmCGSAGG+EUBCAEwDQYJKoZIhvcNAQEFBQADgYEAVaxj6t6h3dKQX58Lzna+E1GPk9kFK8gbd0utaVCh7t7c/dsH6eg5lNyrcnkvBr+rgXDEqO3qUzTt7x5T2QbHVivRXPTRio60K7E3kEgIQiXFPorLf+tvBNFtxXSi96J8e2A8d80OzkgCfwEvtps34CoqNtzVhdas5T9Ub5YeBa8=`,
			rule:     rule,
			expected: `["mail.google.com","Thawte SGC CA"]`,
		},
		{
			note:     "PEM, single cert, b64",
			certs:    `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tDQpNSUlGZHpDQ0JGK2dBd0lCQWdJU0EzTnJpQUV1cy8rY3ZmbHZoVlFPVzV6VE1BMEdDU3FHU0liM0RRRUJDd1VBDQpNRW94Q3pBSkJnTlZCQVlUQWxWVE1SWXdGQVlEVlFRS0V3MU1aWFFuY3lCRmJtTnllWEIwTVNNd0lRWURWUVFEDQpFeHBNWlhRbmN5QkZibU55ZVhCMElFRjFkR2h2Y21sMGVTQllNekFlRncweU1EQTNNVEF4TmpBd016QmFGdzB5DQpNREV3TURneE5qQXdNekJhTUI0eEhEQWFCZ05WQkFNVEUyOXdaVzV3YjJ4cFkzbGhaMlZ1ZEM1dmNtY3dnZ0VpDQpNQTBHQ1NxR1NJYjNEUUVCQVFVQUE0SUJEd0F3Z2dFS0FvSUJBUUN5eThIWlhWVEoyVFNIWFlub0wrQ0tZcG80DQp3ejF3b3dVY2R0L1hCZ04wOGYzN054YU5rK1ZBajhHRDJzNnpob0hMU2h5WVMyUFZvc2Y3eHVtdnlHOTE0UExwDQpJSE85V21DYVpNcXdFeXZNTS9WRTlkQmtLZmFUbzc4QlQ2YVh5Sm1ua2pwZUZtQk9HczN1UDViVUFSajNPbm5yDQo3QW9zOWo0NXJncnl0cGVsWVRNbExpNmpWdEJ2NVJJWnVNb0oxNVcyNTJ0OGVJZ3NPcTU3YWQwQm9iZXl5NFR1DQpHaHZlUDBWM3ZVSnZJM2licUg1RTljV3pJMmY4VXRvaXJVTmYwSjN0Y25nOEpxU091dXpXRFlXclJEQXpRYkpZDQpxS3p2VkRjTitwdHFWN0daNkp1cUhoZHdnRGVxQk9zdmVEYnpBQXlZU1ZQSmpSV1llYThNeGxNN09YYnRBZ01CDQpBQUdqZ2dLQk1JSUNmVEFPQmdOVkhROEJBZjhFQkFNQ0JhQXdIUVlEVlIwbEJCWXdGQVlJS3dZQkJRVUhBd0VHDQpDQ3NHQVFVRkJ3TUNNQXdHQTFVZEV3RUIvd1FDTUFBd0hRWURWUjBPQkJZRUZIRHdlYjZLcHJTdldydy92UjZrDQp3VFZwdWRQdE1COEdBMVVkSXdRWU1CYUFGS2hLYW1NRWZkMjY1dEU1dDZaRlplL3pxT3loTUc4R0NDc0dBUVVGDQpCd0VCQkdNd1lUQXVCZ2dyQmdFRkJRY3dBWVlpYUhSMGNEb3ZMMjlqYzNBdWFXNTBMWGd6TG14bGRITmxibU55DQplWEIwTG05eVp6QXZCZ2dyQmdFRkJRY3dBb1lqYUhSMGNEb3ZMMk5sY25RdWFXNTBMWGd6TG14bGRITmxibU55DQplWEIwTG05eVp5OHdOd1lEVlIwUkJEQXdMb0lUYjNCbGJuQnZiR2xqZVdGblpXNTBMbTl5WjRJWGQzZDNMbTl3DQpaVzV3YjJ4cFkzbGhaMlZ1ZEM1dmNtY3dUQVlEVlIwZ0JFVXdRekFJQmdabmdRd0JBZ0V3TndZTEt3WUJCQUdDDQozeE1CQVFFd0tEQW1CZ2dyQmdFRkJRY0NBUllhYUhSMGNEb3ZMMk53Y3k1c1pYUnpaVzVqY25sd2RDNXZjbWN3DQpnZ0VFQmdvckJnRUVBZFo1QWdRQ0JJSDFCSUh5QVBBQWRnQmVwM1A1MzFiQTU3VTJTSDNRU2VBeWVwR2FESVNoDQpFaEtFR0hXV2dYRkZXQUFBQVhNNXE5dkRBQUFFQXdCSE1FVUNJUUNSSHFncnRsMDdZNlRyeWZNbVFONlROS1JWDQptMUxUeTl2STNNaC9rcmJTUVFJZ1lnVkFLd1hSb1BSK0JOMXBjSmJKdjNBaXZiaDZFN0w5ODdyTVNFUWs1Vm9BDQpkZ0N5SGdYTWk2TE5paUJPaDJiNUs3bUtKU0JuYTlyNmNPZXlTVk10NzR1UVhnQUFBWE01cTl1dUFBQUVBd0JIDQpNRVVDSVFEZHJ1VHV0US9VY2hja3FZUSsycDltdXRuclNublFYYTh4TEE0MVlHelpIZ0lnWFhFVEZiR2ZuczJDDQo3WUo4Y0RvWVlBam1kek1nOGs3aEtYUUd1L0tzQWI0d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFHazlwNXl0DQpPYURJUFJQazVJbXBIMWY2ZjAxMG1VTFdQVjVQam42a3pNSFA5ejVuZE16KysxTk92SFY0R1ZCQ29ldUtxMWJwDQpGQ0QrSWdBOXBjSkFFWFEvdTRHcG1iQUtVWnptZk1JYjg5YVJnbkpwMG14OVk0QkJkNDVFeFVXczh3NGNmZ0ZaDQp5WlVlSHZXczFhbnBBY1IyRklacEFWTVFDYUlnak90MmRkUjF4djRhY0N3K21EL0I5b0tmR1pFVWd5SUFOdnBCDQpJRGFiZ2dMU3dGYTlPS0tYUkJWUkFhZm83T2FjMjFIUVU3RTNzWHBoYUhaR2ZuMkYyN2REL3FvcVVjTHFyNGxDDQpjN2xORTBZR3A2cithUG85VkxjSDJWMGxONHQrMVZiVkFyd0t6bnNOZGNRbndLQmV0Z3F2WnJnTGc0K3FqbzR5DQp1aXhKWTM4WFUvYjdiYVU9DQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tDQo=`,
			rule:     rule,
			expected: `["openpolicyagent.org"]`,
		},
		{
			note: "PEM, single cert, string",
			certs: `-----BEGIN CERTIFICATE-----
MIIFdzCCBF+gAwIBAgISA3NriAEus/+cvflvhVQOW5zTMA0GCSqGSIb3DQEBCwUA
MEoxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MSMwIQYDVQQD
ExpMZXQncyBFbmNyeXB0IEF1dGhvcml0eSBYMzAeFw0yMDA3MTAxNjAwMzBaFw0y
MDEwMDgxNjAwMzBaMB4xHDAaBgNVBAMTE29wZW5wb2xpY3lhZ2VudC5vcmcwggEi
MA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCyy8HZXVTJ2TSHXYnoL+CKYpo4
wz1wowUcdt/XBgN08f37NxaNk+VAj8GD2s6zhoHLShyYS2PVosf7xumvyG914PLp
IHO9WmCaZMqwEyvMM/VE9dBkKfaTo78BT6aXyJmnkjpeFmBOGs3uP5bUARj3Onnr
7Aos9j45rgrytpelYTMlLi6jVtBv5RIZuMoJ15W252t8eIgsOq57ad0Bobeyy4Tu
GhveP0V3vUJvI3ibqH5E9cWzI2f8UtoirUNf0J3tcng8JqSOuuzWDYWrRDAzQbJY
qKzvVDcN+ptqV7GZ6JuqHhdwgDeqBOsveDbzAAyYSVPJjRWYea8MxlM7OXbtAgMB
AAGjggKBMIICfTAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEG
CCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFHDweb6KprSvWrw/vR6k
wTVpudPtMB8GA1UdIwQYMBaAFKhKamMEfd265tE5t6ZFZe/zqOyhMG8GCCsGAQUF
BwEBBGMwYTAuBggrBgEFBQcwAYYiaHR0cDovL29jc3AuaW50LXgzLmxldHNlbmNy
eXB0Lm9yZzAvBggrBgEFBQcwAoYjaHR0cDovL2NlcnQuaW50LXgzLmxldHNlbmNy
eXB0Lm9yZy8wNwYDVR0RBDAwLoITb3BlbnBvbGljeWFnZW50Lm9yZ4IXd3d3Lm9w
ZW5wb2xpY3lhZ2VudC5vcmcwTAYDVR0gBEUwQzAIBgZngQwBAgEwNwYLKwYBBAGC
3xMBAQEwKDAmBggrBgEFBQcCARYaaHR0cDovL2Nwcy5sZXRzZW5jcnlwdC5vcmcw
ggEEBgorBgEEAdZ5AgQCBIH1BIHyAPAAdgBep3P531bA57U2SH3QSeAyepGaDISh
EhKEGHWWgXFFWAAAAXM5q9vDAAAEAwBHMEUCIQCRHqgrtl07Y6TryfMmQN6TNKRV
m1LTy9vI3Mh/krbSQQIgYgVAKwXRoPR+BN1pcJbJv3Aivbh6E7L987rMSEQk5VoA
dgCyHgXMi6LNiiBOh2b5K7mKJSBna9r6cOeySVMt74uQXgAAAXM5q9uuAAAEAwBH
MEUCIQDdruTutQ/UchckqYQ+2p9mutnrSnnQXa8xLA41YGzZHgIgXXETFbGfns2C
7YJ8cDoYYAjmdzMg8k7hKXQGu/KsAb4wDQYJKoZIhvcNAQELBQADggEBAGk9p5yt
OaDIPRPk5ImpH1f6f010mULWPV5Pjn6kzMHP9z5ndMz++1NOvHV4GVBCoeuKq1bp
FCD+IgA9pcJAEXQ/u4GpmbAKUZzmfMIb89aRgnJp0mx9Y4BBd45ExUWs8w4cfgFZ
yZUeHvWs1anpAcR2FIZpAVMQCaIgjOt2ddR1xv4acCw+mD/B9oKfGZEUgyIANvpB
IDabggLSwFa9OKKXRBVRAafo7Oac21HQU7E3sXphaHZGfn2F27dD/qoqUcLqr4lC
c7lNE0YGp6r+aPo9VLcH2V0lN4t+1VbVArwKznsNdcQnwKBetgqvZrgLg4+qjo4y
uixJY38XU/b7baU=
-----END CERTIFICATE-----`,
			rule:     rule,
			expected: `["openpolicyagent.org"]`,
		},
		{
			note:     "PEM, chain, b64",
			certs:    `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tDQpNSUlGZHpDQ0JGK2dBd0lCQWdJU0EzTnJpQUV1cy8rY3ZmbHZoVlFPVzV6VE1BMEdDU3FHU0liM0RRRUJDd1VBTUVveEN6QUpCZ05WQkFZVEFsVlRNUll3RkFZRFZRUUtFdzFNWlhRbmN5QkZibU55ZVhCME1TTXdJUVlEVlFRREV4cE1aWFFuY3lCRmJtTnllWEIwSUVGMWRHaHZjbWwwZVNCWU16QWVGdzB5TURBM01UQXhOakF3TXpCYUZ3MHlNREV3TURneE5qQXdNekJhTUI0eEhEQWFCZ05WQkFNVEUyOXdaVzV3YjJ4cFkzbGhaMlZ1ZEM1dmNtY3dnZ0VpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCRHdBd2dnRUtBb0lCQVFDeXk4SFpYVlRKMlRTSFhZbm9MK0NLWXBvNHd6MXdvd1VjZHQvWEJnTjA4ZjM3TnhhTmsrVkFqOEdEMnM2emhvSExTaHlZUzJQVm9zZjd4dW12eUc5MTRQTHBJSE85V21DYVpNcXdFeXZNTS9WRTlkQmtLZmFUbzc4QlQ2YVh5Sm1ua2pwZUZtQk9HczN1UDViVUFSajNPbm5yN0FvczlqNDVyZ3J5dHBlbFlUTWxMaTZqVnRCdjVSSVp1TW9KMTVXMjUydDhlSWdzT3E1N2FkMEJvYmV5eTRUdUdodmVQMFYzdlVKdkkzaWJxSDVFOWNXekkyZjhVdG9pclVOZjBKM3Rjbmc4SnFTT3V1eldEWVdyUkRBelFiSllxS3p2VkRjTitwdHFWN0daNkp1cUhoZHdnRGVxQk9zdmVEYnpBQXlZU1ZQSmpSV1llYThNeGxNN09YYnRBZ01CQUFHamdnS0JNSUlDZlRBT0JnTlZIUThCQWY4RUJBTUNCYUF3SFFZRFZSMGxCQll3RkFZSUt3WUJCUVVIQXdFR0NDc0dBUVVGQndNQ01Bd0dBMVVkRXdFQi93UUNNQUF3SFFZRFZSME9CQllFRkhEd2ViNktwclN2V3J3L3ZSNmt3VFZwdWRQdE1COEdBMVVkSXdRWU1CYUFGS2hLYW1NRWZkMjY1dEU1dDZaRlplL3pxT3loTUc4R0NDc0dBUVVGQndFQkJHTXdZVEF1QmdnckJnRUZCUWN3QVlZaWFIUjBjRG92TDI5amMzQXVhVzUwTFhnekxteGxkSE5sYm1OeWVYQjBMbTl5WnpBdkJnZ3JCZ0VGQlFjd0FvWWphSFIwY0RvdkwyTmxjblF1YVc1MExYZ3pMbXhsZEhObGJtTnllWEIwTG05eVp5OHdOd1lEVlIwUkJEQXdMb0lUYjNCbGJuQnZiR2xqZVdGblpXNTBMbTl5WjRJWGQzZDNMbTl3Wlc1d2IyeHBZM2xoWjJWdWRDNXZjbWN3VEFZRFZSMGdCRVV3UXpBSUJnWm5nUXdCQWdFd053WUxLd1lCQkFHQzN4TUJBUUV3S0RBbUJnZ3JCZ0VGQlFjQ0FSWWFhSFIwY0RvdkwyTndjeTVzWlhSelpXNWpjbmx3ZEM1dmNtY3dnZ0VFQmdvckJnRUVBZFo1QWdRQ0JJSDFCSUh5QVBBQWRnQmVwM1A1MzFiQTU3VTJTSDNRU2VBeWVwR2FESVNoRWhLRUdIV1dnWEZGV0FBQUFYTTVxOXZEQUFBRUF3QkhNRVVDSVFDUkhxZ3J0bDA3WTZUcnlmTW1RTjZUTktSVm0xTFR5OXZJM01oL2tyYlNRUUlnWWdWQUt3WFJvUFIrQk4xcGNKYkp2M0FpdmJoNkU3TDk4N3JNU0VRazVWb0FkZ0N5SGdYTWk2TE5paUJPaDJiNUs3bUtKU0JuYTlyNmNPZXlTVk10NzR1UVhnQUFBWE01cTl1dUFBQUVBd0JITUVVQ0lRRGRydVR1dFEvVWNoY2txWVErMnA5bXV0bnJTbm5RWGE4eExBNDFZR3paSGdJZ1hYRVRGYkdmbnMyQzdZSjhjRG9ZWUFqbWR6TWc4azdoS1hRR3UvS3NBYjR3RFFZSktvWklodmNOQVFFTEJRQURnZ0VCQUdrOXA1eXRPYURJUFJQazVJbXBIMWY2ZjAxMG1VTFdQVjVQam42a3pNSFA5ejVuZE16KysxTk92SFY0R1ZCQ29ldUtxMWJwRkNEK0lnQTlwY0pBRVhRL3U0R3BtYkFLVVp6bWZNSWI4OWFSZ25KcDBteDlZNEJCZDQ1RXhVV3M4dzRjZmdGWnlaVWVIdldzMWFucEFjUjJGSVpwQVZNUUNhSWdqT3QyZGRSMXh2NGFjQ3crbUQvQjlvS2ZHWkVVZ3lJQU52cEJJRGFiZ2dMU3dGYTlPS0tYUkJWUkFhZm83T2FjMjFIUVU3RTNzWHBoYUhaR2ZuMkYyN2REL3FvcVVjTHFyNGxDYzdsTkUwWUdwNnIrYVBvOVZMY0gyVjBsTjR0KzFWYlZBcndLem5zTmRjUW53S0JldGdxdlpyZ0xnNCtxam80eXVpeEpZMzhYVS9iN2JhVT0NCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0NCi0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQ0KTUlJRWtqQ0NBM3FnQXdJQkFnSVFDZ0ZCUWdBQUFWT0ZjMm9MaGV5bkNEQU5CZ2txaGtpRzl3MEJBUXNGQURBL01TUXdJZ1lEVlFRS0V4dEVhV2RwZEdGc0lGTnBaMjVoZEhWeVpTQlVjblZ6ZENCRGJ5NHhGekFWQmdOVkJBTVREa1JUVkNCU2IyOTBJRU5CSUZnek1CNFhEVEUyTURNeE56RTJOREEwTmxvWERUSXhNRE14TnpFMk5EQTBObG93U2pFTE1Ba0dBMVVFQmhNQ1ZWTXhGakFVQmdOVkJBb1REVXhsZENkeklFVnVZM0o1Y0hReEl6QWhCZ05WQkFNVEdreGxkQ2R6SUVWdVkzSjVjSFFnUVhWMGFHOXlhWFI1SUZnek1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBbk5NTThGcmxMa2UzY2wwM2c3Tm9ZekRxMXpVbUdTWGh2YjQxOFhDU0w3ZTRTMEVGcTZtZU5RaFk3TEVxeEdpSEM2UGpkZVRtODZkaWNicDVnV0FmMTVHYW4vUFFlR2R4eUdrT2xaSFAvdWFaNldBOFNNeCt5azEzRWlTZFJ4dGE2N25zSGpjQUhKeXNlNmNGNnM1SzY3MUI1VGFZdWN2OWJUeVdhTjhqS2tLUURJWjBaOGgvcFpxNFVtRVVFejlsNllLSHk5djZEbGIyaG9uemhUK1hocSt3M0JydmF3MlZGbjNFSzZCbHNwa0VObldBYTZ4Szh4dVFTWGd2b3BaUEtpQWxLUVRHZE1EUU1jMlBNVGlWRnJxb003aEQ4YkVmd3pCL29ua3hFejB0TnZqai9QSXphcms1TWNXdnhJME5IV1FXTTZyNmhDbTIxQXZBMkgzRGt3SURBUUFCbzRJQmZUQ0NBWGt3RWdZRFZSMFRBUUgvQkFnd0JnRUIvd0lCQURBT0JnTlZIUThCQWY4RUJBTUNBWVl3ZndZSUt3WUJCUVVIQVFFRWN6QnhNRElHQ0NzR0FRVUZCekFCaGlab2RIUndPaTh2YVhOeVp5NTBjblZ6ZEdsa0xtOWpjM0F1YVdSbGJuUnlkWE4wTG1OdmJUQTdCZ2dyQmdFRkJRY3dBb1l2YUhSMGNEb3ZMMkZ3Y0hNdWFXUmxiblJ5ZFhOMExtTnZiUzl5YjI5MGN5OWtjM1J5YjI5MFkyRjRNeTV3TjJNd0h3WURWUjBqQkJnd0ZvQVV4S2V4cEhzc2NmcmI0VXVRZGYvRUZXQ0ZpUkF3VkFZRFZSMGdCRTB3U3pBSUJnWm5nUXdCQWdFd1B3WUxLd1lCQkFHQzN4TUJBUUV3TURBdUJnZ3JCZ0VGQlFjQ0FSWWlhSFIwY0RvdkwyTndjeTV5YjI5MExYZ3hMbXhsZEhObGJtTnllWEIwTG05eVp6QThCZ05WSFI4RU5UQXpNREdnTDZBdGhpdG9kSFJ3T2k4dlkzSnNMbWxrWlc1MGNuVnpkQzVqYjIwdlJGTlVVazlQVkVOQldETkRVa3d1WTNKc01CMEdBMVVkRGdRV0JCU29TbXBqQkgzZHV1YlJPYmVtUldYdjg2anNvVEFOQmdrcWhraUc5dzBCQVFzRkFBT0NBUUVBM1RQWEVmTmpXRGpkR0JYN0NWVytkbGE1Y0VpbGFVY25lOElrQ0pMeFdoOUtFaWszSkhSUkhHSm91TTJWY0dmbDk2UzhUaWhSelp2b3JvZWQ2dGk2V3FFQm10enczV29kYXRnK1Z5T2VwaDRFWXByLzF3WEt0eDgvd0FwSXZKU3d0bVZpNE1GVTVhTXFyU0RFNmVhNzNNajJ0Y015bzVqTWQ2am1lV1VISzhzby9qb1dVb0hPVWd3dVg0UG8xUVl6KzNkc3prRHFNcDRma2x4QndYUnNXMTBLWHpQTVRaK3NPUEF2ZXl4aW5kbWprVzhsR3krUXNSbEdQZlorRzZaNmg3bWplbTBZK2lXbGtZY1Y0UElXTDFpd0JpOHNhQ2JHUzVqTjJwOE0rWCtRN1VOS0VrUk9iM042S09xa3FtNTdUSDJIM2VESkFrU25oNi9ETkZ1MFFnPT0NCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0NCi0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQ0KTUlJRFNqQ0NBaktnQXdJQkFnSVFSSyt3Z05hako3cUpNRG1HTHZoQWF6QU5CZ2txaGtpRzl3MEJBUVVGQURBL01TUXdJZ1lEVlFRS0V4dEVhV2RwZEdGc0lGTnBaMjVoZEhWeVpTQlVjblZ6ZENCRGJ5NHhGekFWQmdOVkJBTVREa1JUVkNCU2IyOTBJRU5CSUZnek1CNFhEVEF3TURrek1ESXhNVEl4T1ZvWERUSXhNRGt6TURFME1ERXhOVm93UHpFa01DSUdBMVVFQ2hNYlJHbG5hWFJoYkNCVGFXZHVZWFIxY21VZ1ZISjFjM1FnUTI4dU1SY3dGUVlEVlFRREV3NUVVMVFnVW05dmRDQkRRU0JZTXpDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTit2NlpkUUNJTlh0TXhpWmZhUWd1ekgweXhyTU1wYjdObkRmY2RBd1JnVWkrRG9NM1pKS3VNL0lVbVRyRTRPcno1SXkyWHUvTk1oRDJYU0t0a3lqNHpsOTNld0VudTFsY0NKbzZtNjdYTXVlZ3dHTW9PaWZvb1VNTTBSb09FcU9MbDVDakg5VUwyQVpkKzNVV09EeU9LSVllcExZWUhzVW11NW91SkxHaWlmU0tPZUROb0pqajRYTGg3ZElOOWJ4aXFLcXk2OWNLM0ZDeG9sa0hSeXhYdHFxelRXTUluLzVXZ1RlMVFMeU5hdTdGcWNraDQ5WkxPTXh0Ky95VUZ3N0JaeTFTYnNPRlU1UTlEOC9SaGNRUEdYNjlXYW00MGR1dG9sdWNiWTM4RVZBanFyMm03eFBpNzFYQWljUE5hRGFlUVFteGtxdGlsWDQrVTltNS93QWwwQ0F3RUFBYU5DTUVBd0R3WURWUjBUQVFIL0JBVXdBd0VCL3pBT0JnTlZIUThCQWY4RUJBTUNBUVl3SFFZRFZSME9CQllFRk1TbnNhUjdMSEg2MitGTGtIWC94QlZnaFlrUU1BMEdDU3FHU0liM0RRRUJCUVVBQTRJQkFRQ2pHaXliRndCY3FSN3VLR1kzT3IrRHh6OUx3d21nbFNCZDQ5bFpSTkkrRFQ2OWlrdWdkQi9PRUlLY2RCb2RmcGdhM2NzVFM3TWdST1NSNmN6OGZhWGJhdVgrNXYzZ1R0MjNBRHExY0Vtdjh1WHJBdkhSQW9zWnk1UTZYa2pFR0I1WUdWOGVBbHJ3RFBHeHJhbmNXWWFMYnVtUjlZYksrcmxtTTZwWlc4N2lweFp6UjhzcnpKbXdOMGpQNDFaTDljOFBESEl5aDhid1JMdFRjbTFEOVNaSW1sSm50MWlyL21kMmNYamJEYUpXRkJNNUpER0ZvcWdDV2pCSDRkMVFCN3dDQ1pBQTYyUmpZSnNXdklqSkV1YlNmWkdMK1QweWpXVzA2WHl4VjNicXhiWW9PYjhWWlJ6STluZVdhZ3FOZHd2WWtRc0VqZ2ZiS2JZSzdwMkNOVFVRDQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tDQo=`,
			rule:     rule,
			expected: `["openpolicyagent.org", "Let's Encrypt Authority X3", "DST Root CA X3"]`,
		},
		{
			note: "PEM, chain, string",
			certs: `-----BEGIN CERTIFICATE-----
MIIFdzCCBF+gAwIBAgISA3NriAEus/+cvflvhVQOW5zTMA0GCSqGSIb3DQEBCwUAMEoxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MSMwIQYDVQQDExpMZXQncyBFbmNyeXB0IEF1dGhvcml0eSBYMzAeFw0yMDA3MTAxNjAwMzBaFw0yMDEwMDgxNjAwMzBaMB4xHDAaBgNVBAMTE29wZW5wb2xpY3lhZ2VudC5vcmcwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCyy8HZXVTJ2TSHXYnoL+CKYpo4wz1wowUcdt/XBgN08f37NxaNk+VAj8GD2s6zhoHLShyYS2PVosf7xumvyG914PLpIHO9WmCaZMqwEyvMM/VE9dBkKfaTo78BT6aXyJmnkjpeFmBOGs3uP5bUARj3Onnr7Aos9j45rgrytpelYTMlLi6jVtBv5RIZuMoJ15W252t8eIgsOq57ad0Bobeyy4TuGhveP0V3vUJvI3ibqH5E9cWzI2f8UtoirUNf0J3tcng8JqSOuuzWDYWrRDAzQbJYqKzvVDcN+ptqV7GZ6JuqHhdwgDeqBOsveDbzAAyYSVPJjRWYea8MxlM7OXbtAgMBAAGjggKBMIICfTAOBgNVHQ8BAf8EBAMCBaAwHQYDVR0lBBYwFAYIKwYBBQUHAwEGCCsGAQUFBwMCMAwGA1UdEwEB/wQCMAAwHQYDVR0OBBYEFHDweb6KprSvWrw/vR6kwTVpudPtMB8GA1UdIwQYMBaAFKhKamMEfd265tE5t6ZFZe/zqOyhMG8GCCsGAQUFBwEBBGMwYTAuBggrBgEFBQcwAYYiaHR0cDovL29jc3AuaW50LXgzLmxldHNlbmNyeXB0Lm9yZzAvBggrBgEFBQcwAoYjaHR0cDovL2NlcnQuaW50LXgzLmxldHNlbmNyeXB0Lm9yZy8wNwYDVR0RBDAwLoITb3BlbnBvbGljeWFnZW50Lm9yZ4IXd3d3Lm9wZW5wb2xpY3lhZ2VudC5vcmcwTAYDVR0gBEUwQzAIBgZngQwBAgEwNwYLKwYBBAGC3xMBAQEwKDAmBggrBgEFBQcCARYaaHR0cDovL2Nwcy5sZXRzZW5jcnlwdC5vcmcwggEEBgorBgEEAdZ5AgQCBIH1BIHyAPAAdgBep3P531bA57U2SH3QSeAyepGaDIShEhKEGHWWgXFFWAAAAXM5q9vDAAAEAwBHMEUCIQCRHqgrtl07Y6TryfMmQN6TNKRVm1LTy9vI3Mh/krbSQQIgYgVAKwXRoPR+BN1pcJbJv3Aivbh6E7L987rMSEQk5VoAdgCyHgXMi6LNiiBOh2b5K7mKJSBna9r6cOeySVMt74uQXgAAAXM5q9uuAAAEAwBHMEUCIQDdruTutQ/UchckqYQ+2p9mutnrSnnQXa8xLA41YGzZHgIgXXETFbGfns2C7YJ8cDoYYAjmdzMg8k7hKXQGu/KsAb4wDQYJKoZIhvcNAQELBQADggEBAGk9p5ytOaDIPRPk5ImpH1f6f010mULWPV5Pjn6kzMHP9z5ndMz++1NOvHV4GVBCoeuKq1bpFCD+IgA9pcJAEXQ/u4GpmbAKUZzmfMIb89aRgnJp0mx9Y4BBd45ExUWs8w4cfgFZyZUeHvWs1anpAcR2FIZpAVMQCaIgjOt2ddR1xv4acCw+mD/B9oKfGZEUgyIANvpBIDabggLSwFa9OKKXRBVRAafo7Oac21HQU7E3sXphaHZGfn2F27dD/qoqUcLqr4lCc7lNE0YGp6r+aPo9VLcH2V0lN4t+1VbVArwKznsNdcQnwKBetgqvZrgLg4+qjo4yuixJY38XU/b7baU=
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIEkjCCA3qgAwIBAgIQCgFBQgAAAVOFc2oLheynCDANBgkqhkiG9w0BAQsFADA/MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMTDkRTVCBSb290IENBIFgzMB4XDTE2MDMxNzE2NDA0NloXDTIxMDMxNzE2NDA0NlowSjELMAkGA1UEBhMCVVMxFjAUBgNVBAoTDUxldCdzIEVuY3J5cHQxIzAhBgNVBAMTGkxldCdzIEVuY3J5cHQgQXV0aG9yaXR5IFgzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnNMM8FrlLke3cl03g7NoYzDq1zUmGSXhvb418XCSL7e4S0EFq6meNQhY7LEqxGiHC6PjdeTm86dicbp5gWAf15Gan/PQeGdxyGkOlZHP/uaZ6WA8SMx+yk13EiSdRxta67nsHjcAHJyse6cF6s5K671B5TaYucv9bTyWaN8jKkKQDIZ0Z8h/pZq4UmEUEz9l6YKHy9v6Dlb2honzhT+Xhq+w3Brvaw2VFn3EK6BlspkENnWAa6xK8xuQSXgvopZPKiAlKQTGdMDQMc2PMTiVFrqoM7hD8bEfwzB/onkxEz0tNvjj/PIzark5McWvxI0NHWQWM6r6hCm21AvA2H3DkwIDAQABo4IBfTCCAXkwEgYDVR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAYYwfwYIKwYBBQUHAQEEczBxMDIGCCsGAQUFBzABhiZodHRwOi8vaXNyZy50cnVzdGlkLm9jc3AuaWRlbnRydXN0LmNvbTA7BggrBgEFBQcwAoYvaHR0cDovL2FwcHMuaWRlbnRydXN0LmNvbS9yb290cy9kc3Ryb290Y2F4My5wN2MwHwYDVR0jBBgwFoAUxKexpHsscfrb4UuQdf/EFWCFiRAwVAYDVR0gBE0wSzAIBgZngQwBAgEwPwYLKwYBBAGC3xMBAQEwMDAuBggrBgEFBQcCARYiaHR0cDovL2Nwcy5yb290LXgxLmxldHNlbmNyeXB0Lm9yZzA8BgNVHR8ENTAzMDGgL6AthitodHRwOi8vY3JsLmlkZW50cnVzdC5jb20vRFNUUk9PVENBWDNDUkwuY3JsMB0GA1UdDgQWBBSoSmpjBH3duubRObemRWXv86jsoTANBgkqhkiG9w0BAQsFAAOCAQEA3TPXEfNjWDjdGBX7CVW+dla5cEilaUcne8IkCJLxWh9KEik3JHRRHGJouM2VcGfl96S8TihRzZvoroed6ti6WqEBmtzw3Wodatg+VyOeph4EYpr/1wXKtx8/wApIvJSwtmVi4MFU5aMqrSDE6ea73Mj2tcMyo5jMd6jmeWUHK8so/joWUoHOUgwuX4Po1QYz+3dszkDqMp4fklxBwXRsW10KXzPMTZ+sOPAveyxindmjkW8lGy+QsRlGPfZ+G6Z6h7mjem0Y+iWlkYcV4PIWL1iwBi8saCbGS5jN2p8M+X+Q7UNKEkROb3N6KOqkqm57TH2H3eDJAkSnh6/DNFu0Qg==
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
MIIDSjCCAjKgAwIBAgIQRK+wgNajJ7qJMDmGLvhAazANBgkqhkiG9w0BAQUFADA/MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMTDkRTVCBSb290IENBIFgzMB4XDTAwMDkzMDIxMTIxOVoXDTIxMDkzMDE0MDExNVowPzEkMCIGA1UEChMbRGlnaXRhbCBTaWduYXR1cmUgVHJ1c3QgQ28uMRcwFQYDVQQDEw5EU1QgUm9vdCBDQSBYMzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN+v6ZdQCINXtMxiZfaQguzH0yxrMMpb7NnDfcdAwRgUi+DoM3ZJKuM/IUmTrE4Orz5Iy2Xu/NMhD2XSKtkyj4zl93ewEnu1lcCJo6m67XMuegwGMoOifooUMM0RoOEqOLl5CjH9UL2AZd+3UWODyOKIYepLYYHsUmu5ouJLGiifSKOeDNoJjj4XLh7dIN9bxiqKqy69cK3FCxolkHRyxXtqqzTWMIn/5WgTe1QLyNau7Fqckh49ZLOMxt+/yUFw7BZy1SbsOFU5Q9D8/RhcQPGX69Wam40dutolucbY38EVAjqr2m7xPi71XAicPNaDaeQQmxkqtilX4+U9m5/wAl0CAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYwHQYDVR0OBBYEFMSnsaR7LHH62+FLkHX/xBVghYkQMA0GCSqGSIb3DQEBBQUAA4IBAQCjGiybFwBcqR7uKGY3Or+Dxz9LwwmglSBd49lZRNI+DT69ikugdB/OEIKcdBodfpga3csTS7MgROSR6cz8faXbauX+5v3gTt23ADq1cEmv8uXrAvHRAosZy5Q6XkjEGB5YGV8eAlrwDPGxrancWYaLbumR9YbK+rlmM6pZW87ipxZzR8srzJmwN0jP41ZL9c8PDHIyh8bwRLtTcm1D9SZImlJnt1ir/md2cXjbDaJWFBM5JDGFoqgCWjBH4d1QB7wCCZAA62RjYJsWvIjJEubSfZGL+T0yjWW06XyxV3bqxbYoOb8VZRzI9neWagqNdwvYkQsEjgfbKbYK7p2CNTUQ
-----END CERTIFICATE-----`,
			rule:     rule,
			expected: `["openpolicyagent.org", "Let's Encrypt Authority X3", "DST Root CA X3"]`,
		},
		{
			note:     "invalid DER or PEM data, b64",
			certs:    `YmFkc3RyaW5n`,
			rule:     rule,
			expected: &Error{Code: BuiltinErr, Message: "asn1: structure error"},
		},
		{
			note:     "invalid DER or PEM data, string",
			certs:    `foobar`,
			rule:     rule,
			expected: &Error{Code: BuiltinErr, Message: "illegal base64"},
		},
	}

	data := loadSmallTestData()

	for _, tc := range tests {
		rules := []string{
			fmt.Sprintf("certs = %q { true }", tc.certs),
			tc.rule,
		}
		runTopDownTestCase(t, data, tc.note, rules, tc.expected)
	}

}

func TestCryptoX509ParseCertificateRequest(t *testing.T) {
	rule := `
		p = x {
		  parsed := crypto.x509.parse_certificate_request(csr)
		  x := parsed.Subject.CommonName
		}
	`

	tests := []struct {
		note     string
		csr      string
		rule     string
		expected interface{}
	}{
		{
			note:     "PEM, b64",
			csr:      `LS0tLS1CRUdJTiBDRVJUSUZJQ0FURSBSRVFVRVNULS0tLS0KTUlJQ21EQ0NBWUFDQVFBd1V6RUxNQWtHQTFVRUJoTUNWVk14RkRBU0JnTlZCQU1NQzJWNFlXMXdiR1V1WTI5dApNUW93Q0FZRFZRUUhEQUVnTVFvd0NBWURWUVFLREFFZ01Rb3dDQVlEVlFRSURBRWdNUW93Q0FZRFZRUUxEQUVnCk1JSUJJakFOQmdrcWhraUc5dzBCQVFFRkFBT0NBUThBTUlJQkNnS0NBUUVBMlpkaG1zaERBVTBYYnhnTk1GQWsKeEdWQnNjaHdWb2s5dXBBU2ZVWDA4VFlqMFZrV0VxNitmemdOdmRQSnd6Nm1lUDlnL01hRmhPYW91Nmh1UEhmbwpTVTlKN1FiTW56Uktsc0VJTzNodEM1QUt3OXYyZldVZGpCQS92Q1dZdXU1aUc1ZTdtUHNXWjd1cGxuVGZSekM4ClJLK0srWXJtNEQ4NHE1bHR5NEMzS2tRc0FjU0xQZk9MMXMvYjJyV21KR0FoV3NSa2doTVk2V3dza3VYWXRINTkKRzl5VURHUUhoalprcHFlZFY0OUM4c0NwMU8vWVpvU0hncDdHK0JiaFRta05CRzY3OFZHREplTnB3SG96dnRjVQpyQVNGRFJ4WnhPdTFHRzE3L1FiVW9SNVVkOTNwaUtaU0U2UHVDU2VCcy9UQmFJc3ZwUGtudVhkOXI4WGovbVd5CmtRSURBUUFCb0FBd0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFBeDJkaCtkMU1CaEwwaDJYZklxaDVEYy9lYWoKU0xadGFNTWlJY1h1cC96UTl2eENXSkZlSGYzczBJdXliMEhkMlZNZ1BSYU8ydWRkY2JZdFFlKzJnWUtrTzFMWApCdHdQcXcwWHAweUF2dDUxRzJvZmVCbCtFa0ptNjk3RlNtemg4eDJJZFFBSkMzWi9ROFdMVmh3NFg2WlVicnhqCjJnTjJmaVhjS0RKbGVkcUgxY2V4WVVvbnlLSDZubG4wbzQzUUtEOFlSZG9hNVFqb3Ixb0JkY3dSTTA0VDM4ak0KV1B3d2JZTjNrVE9Ea0tiaVFVVWxVeFZuNnFnZTlNTWt0c0lOWkc0eDY1QmIwaWxTdHExRWQwN2Y5NmVnbHNKaApZVE9VRnZpZDZVSkVEcEJzcjhyZFROSW1JQkhCdkkra1BHS2FqcW83Z0VNc3hFYkNkemFHUTNZZnNYWT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUgUkVRVUVTVC0tLS0t`,
			rule:     rule,
			expected: `"example.com"`,
		},
		{
			note: "PEM, string",
			csr: `-----BEGIN CERTIFICATE REQUEST-----
MIICmDCCAYACAQAwUzELMAkGA1UEBhMCVVMxFDASBgNVBAMMC2V4YW1wbGUuY29t
MQowCAYDVQQHDAEgMQowCAYDVQQKDAEgMQowCAYDVQQIDAEgMQowCAYDVQQLDAEg
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2ZdhmshDAU0XbxgNMFAk
xGVBschwVok9upASfUX08TYj0VkWEq6+fzgNvdPJwz6meP9g/MaFhOaou6huPHfo
SU9J7QbMnzRKlsEIO3htC5AKw9v2fWUdjBA/vCWYuu5iG5e7mPsWZ7uplnTfRzC8
RK+K+Yrm4D84q5lty4C3KkQsAcSLPfOL1s/b2rWmJGAhWsRkghMY6WwskuXYtH59
G9yUDGQHhjZkpqedV49C8sCp1O/YZoSHgp7G+BbhTmkNBG678VGDJeNpwHozvtcU
rASFDRxZxOu1GG17/QbUoR5Ud93piKZSE6PuCSeBs/TBaIsvpPknuXd9r8Xj/mWy
kQIDAQABoAAwDQYJKoZIhvcNAQELBQADggEBAAx2dh+d1MBhL0h2XfIqh5Dc/eaj
SLZtaMMiIcXup/zQ9vxCWJFeHf3s0Iuyb0Hd2VMgPRaO2uddcbYtQe+2gYKkO1LX
BtwPqw0Xp0yAvt51G2ofeBl+EkJm697FSmzh8x2IdQAJC3Z/Q8WLVhw4X6ZUbrxj
2gN2fiXcKDJledqH1cexYUonyKH6nln0o43QKD8YRdoa5Qjor1oBdcwRM04T38jM
WPwwbYN3kTODkKbiQUUlUxVn6qge9MMktsINZG4x65Bb0ilStq1Ed07f96eglsJh
YTOUFvid6UJEDpBsr8rdTNImIBHBvI+kPGKajqo7gEMsxEbCdzaGQ3YfsXY=
-----END CERTIFICATE REQUEST-----`,
			rule:     rule,
			expected: `"example.com"`,
		},
		{
			note:     "DER, b64",
			csr:      `MIICmDCCAYACAQAwUzELMAkGA1UEBhMCVVMxFDASBgNVBAMMC2V4YW1wbGUuY29tMQowCAYDVQQHDAEgMQowCAYDVQQKDAEgMQowCAYDVQQIDAEgMQowCAYDVQQLDAEgMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2ZdhmshDAU0XbxgNMFAkxGVBschwVok9upASfUX08TYj0VkWEq6+fzgNvdPJwz6meP9g/MaFhOaou6huPHfoSU9J7QbMnzRKlsEIO3htC5AKw9v2fWUdjBA/vCWYuu5iG5e7mPsWZ7uplnTfRzC8RK+K+Yrm4D84q5lty4C3KkQsAcSLPfOL1s/b2rWmJGAhWsRkghMY6WwskuXYtH59G9yUDGQHhjZkpqedV49C8sCp1O/YZoSHgp7G+BbhTmkNBG678VGDJeNpwHozvtcUrASFDRxZxOu1GG17/QbUoR5Ud93piKZSE6PuCSeBs/TBaIsvpPknuXd9r8Xj/mWykQIDAQABoAAwDQYJKoZIhvcNAQELBQADggEBAAx2dh+d1MBhL0h2XfIqh5Dc/eajSLZtaMMiIcXup/zQ9vxCWJFeHf3s0Iuyb0Hd2VMgPRaO2uddcbYtQe+2gYKkO1LXBtwPqw0Xp0yAvt51G2ofeBl+EkJm697FSmzh8x2IdQAJC3Z/Q8WLVhw4X6ZUbrxj2gN2fiXcKDJledqH1cexYUonyKH6nln0o43QKD8YRdoa5Qjor1oBdcwRM04T38jMWPwwbYN3kTODkKbiQUUlUxVn6qge9MMktsINZG4x65Bb0ilStq1Ed07f96eglsJhYTOUFvid6UJEDpBsr8rdTNImIBHBvI+kPGKajqo7gEMsxEbCdzaGQ3YfsXY=`,
			rule:     rule,
			expected: `"example.com"`,
		},
		{
			note:     "invalid DER or PEM data, b64",
			csr:      `YmFkc3RyaW5n`,
			rule:     rule,
			expected: &Error{Code: BuiltinErr, Message: "asn1: structure error"},
		},
		{
			note:     "invalid DER or PEM data, string",
			csr:      `foobar`,
			rule:     rule,
			expected: &Error{Code: BuiltinErr, Message: "illegal base64"},
		},
	}

	data := loadSmallTestData()

	for _, tc := range tests {
		rules := []string{
			fmt.Sprintf("csr = %q { true }", tc.csr),
			tc.rule,
		}
		runTopDownTestCase(t, data, tc.note, rules, tc.expected)
	}
}

func TestCryptoMd5(t *testing.T) {

	tests := []struct {
		note     string
		rule     []string
		expected interface{}
	}{
		{
			note:     "crypto.md5 with string",
			rule:     []string{`p[hash] { hash := crypto.md5("lorem ipsum") }`},
			expected: `["80a751fde577028640c419000e33eba6"]`,
		},
	}

	data := loadSmallTestData()

	for _, tc := range tests {
		runTopDownTestCase(t, data, tc.note, tc.rule, tc.expected)
	}

}

func TestCryptoSha1(t *testing.T) {

	tests := []struct {
		note     string
		rule     []string
		expected interface{}
	}{
		{
			note:     "crypto.sha1 with string",
			rule:     []string{`p[hash] { hash := crypto.sha1("lorem ipsum") }`},
			expected: `["bfb7759a67daeb65410490b4d98bb9da7d1ea2ce"]`,
		},
	}

	data := loadSmallTestData()

	for _, tc := range tests {
		runTopDownTestCase(t, data, tc.note, tc.rule, tc.expected)
	}

}

func TestCryptoSha256(t *testing.T) {

	tests := []struct {
		note     string
		rule     []string
		expected interface{}
	}{
		{
			note:     "crypto.sha256 with string",
			rule:     []string{`p[hash] { hash := crypto.sha256("lorem ipsum") }`},
			expected: `["5e2bf57d3f40c4b6df69daf1936cb766f832374b4fc0259a7cbff06e2f70f269"]`,
		},
	}

	data := loadSmallTestData()

	for _, tc := range tests {
		runTopDownTestCase(t, data, tc.note, tc.rule, tc.expected)
	}

}
