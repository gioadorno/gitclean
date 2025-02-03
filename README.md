# GitClean

An easy to use tool for rebasing and cleaning up branches

### Commands
- Rebase (Fetches the input branch (or default origin/master), rebases with HEAD, force push to HEAD)
- Reset

#### Flags
(Optional) -b, -branch - When doing a rebase, you can add the -b flag to tell git which branch you want to rebase. Default value is origin/master

## Rebase
This will fetch the requested branch (default is origin/master) to retrieve the updated changes. Then it will log out any changes that have been made. If there are no changes GitClean will exit, since there is no rebase to be had. Otherwise it will continue to perform a rebase. If there are conflicts, GitClean should log out the conflicts and exit, so the user can manually resolve any conflicts. If there are no conflicts or all conflicts are resolved (Rerun rebase if an exit was performed) a force push will take place was the rebase is successful.

##### Example
```
gitclean rebase
```
or
```
gitclean rebase -b dashboard-feature
```

## Reset
This will grab the latest commit hash of the parent branch and do a reset --soft to that commit hash. User will then be prompted to type a commit message and select git files to add (default is .). Then it will push to your branch, which should now have one commit.

##### Example
```
gitclean reset

### Enter a commit message:
Initiate dashboard feature

### Git file(s) or path to add. Default is '.':
(left blank for default)
```
