set -e

VERSION=0.2.0

rm -rf .publish/ && mkdir .publish

GOOS=linux GOARCH=amd64 go build -o .publish/terraform-provider-akc_v${VERSION}

cd .publish
zip -r terraform-provider-akc_${VERSION}_linux_amd64.zip terraform-provider-akc_v${VERSION}
shasum -a 256 terraform-provider-akc_${VERSION}_linux_amd64.zip > terraform-provider-akc_${VERSION}_SHA256SUMS
gpg --detach-sign terraform-provider-akc_${VERSION}_SHA256SUMS
rm terraform-provider-akc_v${VERSION}

cd ..

# GOOS=windows GOARCH=amd64 go build -o "terraform-provider-akc_v${VERSION}.exe"