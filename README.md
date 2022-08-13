## Testing
```bash
SHORTCOMMAND_CONFIG=shortcommands-example.yml go run . Test test
SHORTCOMMAND_CONFIG=shortcommands-example.yml go run . Test test2
SHORTCOMMAND_CONFIG=shortcommands-example.yml go run . Test test3
```

## Build Testing
```bash
goreleaser build --snapshot --rm-dist
```

## Build for release
```bash
go tag vX.X.X
git push origin vX.X.X
goreleaser release --rm-dist
```

## Installing on Linux (tested on Ubuntu)
```bash
wget https://github.com/flawiddsouza/shortcommand/releases/download/v0.0.1/shortcommand_0.0.1_Linux_x86_64.tar.gz
tar -xf shortcommand_0.0.1_Linux_x86_64.tar.gz
mv shortcommand ~/.local/bin
rm shortcommand_0.0.1_Linux_x86_64.tar.gz
```
You should then be able to use the `shortcommand` command anywhere you are.
