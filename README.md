# telegram-forwarder

```
export $(cat .env | xargs) && go run ./cmd/tgforwarder
```

```
go mod tidy
```

```
go build -v -x ./cmd/tgforwarder && rm -f tgforwarder
```

```
docker buildx build --no-cache --progress=plain .
```

```
set -a && source .env && set +a && go run ./cmd/tgforwarder
```

```
git tag v0.1.1
git push origin v0.1.1
```

```
git tag -l
git tag -d vX.Y.Z
git push --delete origin vX.Y.Z
git ls-remote --tags origin | grep 'refs/tags/vX.Y.Z$'
```
