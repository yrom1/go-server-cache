curl -X POST -H "Content-Type: application/json" -d '{"key": "apple", "value": "42"}' http://localhost:3333/cache/apple
curl -X POST -H "Content-Type: application/json" -d '{"key": "banana", "value": "666"}' http://localhost:3333/cache/apple
curl http://localhost:3333/cache/apple # get apple returns 42
curl http://localhost:3333/cache/banana # get banana returns 666
curl http://localhost:3333/cache/ # lists apple banana
