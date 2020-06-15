run using:
git --no-pager diff --name-only HEAD $(git merge-base upstream/master HEAD) | go run github.com/k15r/codeowners
