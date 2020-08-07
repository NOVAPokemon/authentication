module github.com/NOVAPokemon/authentication

go 1.13

require (
	github.com/NOVAPokemon/utils v0.0.62
	github.com/bruno-anjos/archimedesHTTPClient v0.0.2
	github.com/pkg/errors v0.9.1
	github.com/sirupsen/logrus v1.6.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
)

replace (
	github.com/NOVAPokemon/utils v0.0.62 => ../utils
	github.com/bruno-anjos/archimedes v0.0.2 => ./../../bruno-anjos/archimedes
	github.com/bruno-anjos/archimedesHTTPClient v0.0.2 => ./../../bruno-anjos/archimedesHTTPClient
)
