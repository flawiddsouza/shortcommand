# shortcommand

Easily run a set of commands quickly using a yaml configuration file

## Usage

Create a yaml file named shortcommand.yml or anything you want. Here's an example:
```
shortcommands:
  - name: Journals
    commands:
      - name: ui
        description: Pull, build and deploy Web-UI
        cwd: ~/Apps/Journals/Web-UI
        do:
          - git pull
          - npm run build
      - name: api
        description: Pull, build and deploy API
        cwd: ~/Apps/Journals/API
        do:
          - git pull
          - shards build --release
          - pm2 restart "Journals API"
          - pm2 reset "Journals API"
  - name: QuickNote
    commands:
      - name: api
        description: Pull, build and deploy API
        cwd: ~/Apps/quick-note/API
        do:
          - git pull
          - pm2 restart "Quick Note API"
          - pm2 reset "Quick Note API"
      - name: restart-api
        description: Restart API in pm2
        do:
          - pm2 restart "Quick Note API"
          - pm2 reset "Quick Note API"
```

Then add this in your ~/.bashrc or ~/.zshrc:
```
export SHORTCOMMAND_CONFIG=/path/to/shortcommand.yml
```

Now you'll be able to use the above defined short commands to run those sets of commands wherever you are using just this:
```bash
shortcommand Journals ui
shortcommand Journals api
shortcommand QuickNote api
```

Note: Commands mentioned under `do:` will execute in sequence and will not continue to the next one if the one that's running fails.

This project was created mainly because I have several apps that I run on my server and finding the right set of commands for deploying an app is a hassle. So this basically documents the set of commands for each of my projects, as well as gives me quick access to them.

## Installing on Linux (tested on Ubuntu)
```bash
wget https://github.com/flawiddsouza/shortcommand/releases/download/v0.0.3/shortcommand_0.0.3_Linux_x86_64.tar.gz
tar -xf shortcommand_0.0.3_Linux_x86_64.tar.gz
mv shortcommand ~/.local/bin
rm shortcommand_0.0.3_Linux_x86_64.tar.gz
```
You should then be able to use the `shortcommand` command anywhere you are.

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
