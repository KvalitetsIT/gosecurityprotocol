module github.com/KvalitetsIT/gosecurityprotocol

go 1.13

replace github.com/russellhaering/goxmldsig => github.com/evtr/goxmldsig v0.0.0-20190907195011-53d9398322c5

require (
	github.com/beevik/etree v1.1.0 // indirect
	github.com/google/go-cmp v0.4.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/russellhaering/gosaml2 v0.3.1
	github.com/russellhaering/goxmldsig v0.0.0-00010101000000-000000000000
	go.uber.org/atomic v1.5.0
	go.uber.org/multierr v1.3.0
	go.uber.org/zap v1.13.0
	gotest.tools v2.2.0+incompatible

	go.mongodb.org/mongo-driver v1.3.2
)
