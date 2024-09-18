This directory contains the certificates embedded for signature verification. The library automatically selects the newest valid certificate for verifying signatures, ensuring seamless transitions when certificate rotations occur.

By including upcoming server certificates before they become active, the switch to the new certificate is seamless on the effective date, avoiding the need for manual intervention.

Each certificate should be stored as a full chain, combining the leaf certificate, intermediate certificates, and the root CA, like so:

```
cat leafCert.pem intermediateCert.pem rootCA.pem > fullchain.pem
openssl crl2pkcs7 -nocrl -certfile fullchain.pem | openssl pkcs7 -print_certs -out fullchain.pem
```

Production certificates should be named sequentially as fiskalcis1.pem, fiskalcis2.pem, etc. When a new certificate is added, the oldest one can be removed, ensuring that at most two certificates (the current and upcoming) are stored at any time.

See ciscert.go for details