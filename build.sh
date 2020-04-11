VERSION=0.1.0

GOOS=linux GOARCH=amd64 go build -o "terraform-provider-akc_v${VERSION}"
GOOS=windows GOARCH=amd64 go build -o "terraform-provider-akc_v${VERSION}.exe"