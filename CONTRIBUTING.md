# Contributing

:+1::tada: Thanks for taking the time to contribute! :tada::+1:

## Reporting issues

Feel free to [open issues](https://github.com/bpicode/fritzctl/issues/new) if you...

* ...found a problem with `fritzctl` or the documentation,
* ...have a question about `fritzctl` or the documentation,
* ...want a new feature in `fritzctl`.

If the issue should not be disclosed in public, contact the [maintainer](https://github.com/bpicode)
and encrypt the message using PGP with the public key [4096R/8A896560](https://pgp.mit.edu/pks/lookup?op=get&search=0x198D1DA18A896560) 

## How to contribute

Contributions to the codebase of `fritzctl` should be proposed in the form of a pull request using
GitHub. To this end create a fork of this repository.
 
Within git...

* Keep the number of commits small. Ideally, "one commit = one feature" or "one commit = one bugfix".
* Try to follow the [7 rules of a great commit message](https://chris.beams.io/posts/git-commit/#seven-rules).
* If a commit addresses a [known issue](https://github.com/bpicode/fritzctl/issues), include the issue
  number in the commit message.
* Don't worry if something went wrong, most mistakes [can be fixed](http://ohshitgit.com).

The coding conventions...

* Format code according to `gofmt`. 
* Test your code.
* If new 3rd party dependencies arise, reflect if those can be avoided.
* `make codequality` gives a hint on common problems.
* Further reading: [golang/CodeReviewComments](https://github.com/golang/go/wiki/CodeReviewComments).

When changes make a large impact...

* Do the changes break backward compatibility of the cli? The commit message and pull request 
  description should say so. It may be that these changes are delayed until the next major version
  release.
* Do the changes break backward compatibility of the API? The same comments from above apply. 
