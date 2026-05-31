# date-tracker
Track multiple recurring tasks in daily granularity.
Use simple subcommands to `show` activities in past days and `tick` any tasks.
You can pass the names or alias of one or more tasks.
The yaml data file is human readable and editable for more flexibility.
Tab completion is enabled.

## Build
`go build -o ~/go/bin/dt`
You can add `export PATH=$PATH:$HOME/go/bin` to ~/.zshrc to run `dt` command everywhere.

## Create yaml for your data
You can refer to the format of example.yaml and create your own task.yaml in ~/go/bin/data
Remember to use absolute path if you need to symlink it.

## Generate tab completion script
This utility is based on Cobra, which can generate tab completion script.

For zsh, put either
```
autoload -U compinit && compinit
source <(./dt completion zsh)
```
or
```
./dt completion zsh > "${fpath[1]}/_dt"
autoload -U compinit && compinit
```
to ~/.zshrc