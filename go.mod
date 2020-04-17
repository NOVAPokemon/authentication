module github.com/NOVAPokemon/authentication

go 1.13

require (
	github.com/NOVAPokemon/utils v0.0.62
	github.com/sirupsen/logrus v1.5.0
	go.mongodb.org/mongo-driver v1.3.1
	golang.org/x/crypto v0.0.0-20190530122614-20be4c3c3ed5
)

replace github.com/NOVAPokemon/utils v0.0.62 => ../utils
