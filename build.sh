set -e

VERSION=0.2.0

rm -rf .publish/ && mkdir .publish

GOOS=linux GOARCH=amd64 go build -o .publish/linux_amd64/terraform-provider-akc_v${VERSION}

cd .publish/linux_amd64
zip -r terraform-provider-akc_${VERSION}_linux_amd64.zip terraform-provider-akc_v${VERSION}
shasum -a 256 terraform-provider-akc_${VERSION}_linux_amd64.zip > terraform-provider-akc_${VERSION}_SHA256SUMS
#gpg --detach-sign terraform-provider-akc_${VERSION}_SHA256SUMS
rm terraform-provider-akc_v${VERSION}

cd ../..

GOOS=windows GOARCH=amd64 go build -o .publish/win_amd64/terraform-provider-akc_v${VERSION}.exe

cd .publish/win_amd64
zip -r terraform-provider-akc_${VERSION}_win_amd64.zip terraform-provider-akc_v${VERSION}.exe
shasum -a 256 terraform-provider-akc_${VERSION}_win_amd64.zip > terraform-provider-akc_${VERSION}_SHA256SUMS
#gpg --detach-sign terraform-provider-akc_${VERSION}_SHA256SUMS
rm terraform-provider-akc_v${VERSION}.exe
