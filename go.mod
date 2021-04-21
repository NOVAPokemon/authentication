module github.com/NOVAPokemon/authentication

go 1.13

require (
	github.com/NOVAPokemon/utils v0.0.62
	github.com/bruno-anjos/archimedesHTTPClient v0.0.2
	github.com/bruno-anjos/cloud-edge-deployment v0.0.1
	github.com/golang/geo v0.0.0-20200730024412-e86565bf3f35
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.7.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
)

replace (
	github.com/NOVAPokemon/utils v0.0.62 => ../utils
	github.com/bruno-anjos/archimedesHTTPClient v0.0.2 => ../../bruno-anjos/archimedesHTTPClient
	github.com/bruno-anjos/cloud-edge-deployment v0.0.1 => ../../bruno-anjos/cloud-edge-deployment
	github.com/nm-morais/demmon-client v1.0.0 => ../../nm-morais/demmon-client
	github.com/nm-morais/demmon-common v1.0.0 => ../../nm-morais/demmon-common
	github.com/nm-morais/demmon-exporter v1.0.2 => ../../nm-morais/demmon-exporter
)
