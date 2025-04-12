module github.com/koshka_backend/routes

go 1.23.7

replace github.com/koshka_backend/config => ../config

replace github.com/koshka_backend/helpers => ../helpers

replace github.com/koshka_backend/middleware => ../middleware

require (
	github.com/gofiber/fiber/v2 v2.52.6
	github.com/koshka_backend/config v0.0.0-00010101000000-000000000000
	github.com/koshka_backend/helpers v0.0.0-00010101000000-000000000000
	github.com/koshka_backend/middleware v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.37.0
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/gofiber/storage/memory v1.3.4 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sendgrid/rest v2.6.9+incompatible // indirect
	github.com/sendgrid/sendgrid-go v3.16.0+incompatible // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
)
