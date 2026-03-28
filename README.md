# date-tracker
Track multiple recurring tasks in daily granularity.
Simply run `dt show` to show the past dates' data and `dt check <task_name>` to check today.

## Generation tab completion script
Cobra based command can generate tab completion script.

For zsh, put either
```
source <(./dt completion zsh)
```
or
```
./dt completion zsh > "${fpath[1]}/_dt"
autoload -U compinit && compinit
```
to ~/.zshrc