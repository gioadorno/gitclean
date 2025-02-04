# GitClean

An easy to use tool for rebasing and cleaning up branches

### Commands
- Clean (Rebase and Reset)
- Rebase (Fetches the input branch (or default origin/master), rebases with HEAD, force push to HEAD)
- Reset

#### Flags
(Optional) -b, -branch - When doing a rebase, you can add the -b flag to tell git which branch you want to rebase. Default value is origin/master

## Rebase
This will fetch the requested branch (default is origin/master) to retrieve the updated changes. Then it will log out any changes that have been made. If there are no changes GitClean will exit, since there is no rebase to be had. Otherwise it will continue to perform a rebase. If there are conflicts, you will have to manually resolve those. Gitclean will let you know to press Enter to continue if all conflicts or resolved or press q to abort the process.

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
## Clean
```
gitclean clean
```
or
```
gitclean clean -b dashboard-feature
```
Does both the rebase and reset functionality
