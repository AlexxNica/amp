# Release Checklist

This is a draft.

## A maintainer's guide to releasing AMP

So you're in charge of a AMP release? Cool. Here's what to do.

If your experience deviates from this document, please document the changes
to keep it up-to-date.

It is important to note that this document assumes that the git remote in your
repository that corresponds to "https://github.com/appcelerator/amp" is named
"origin".  If yours is not (for example, if you've chosen to name it "upstream"
or something similar instead), be sure to adjust the listed snippets for your
local environment accordingly.  If you are not sure what your upstream remote is
named, use a command like `git remote -v` to find out.

If you don't have an upstream remote, you can add one easily using something
like:

```bash
export GITHUBUSER="YOUR_GITHUB_USER"
git remote add origin https://github.com/appcelerator/amp.git
git remote add $GITHUBUSER git@github.com:$GITHUBUSER/amp.git
```

### 1. Pull from master and create a release branch

All releases version numbers will be of the form: vX.Y.Z  where X is the major
version number, Y is the minor version number and Z is the patch release version number.

#### Major releases

The release branch name is just vX.Y because it's going to be the basis for all .Z releases.

```bash
export BASE=vX.Y
export VERSION=vX.Y.Z
git fetch origin
git checkout --track origin/master
git checkout -b release/$BASE
```

This new branch is going to be the base for the release. We need to push it to origin so we
can track the cherry-picked changes and the version bump:

```bash
git push origin release/$BASE
```

When you have the major release branch in origin, we need to create the bump fork branch
that we'll push to our fork:

```bash
git checkout -b bump_$VERSION
```

#### Patch releases

If we have the release branch in origin, we can create the forked bump branch from it directly:

```bash
export VERSION=vX.Y.Z
export PATCH=vX.Y.Z+1
git fetch origin
git checkout --track origin/release/$BASE
git checkout -b bump_$PATCH
```

We cherry-pick only the commits we want into the bump branch:

```bash
# get the commits ids we want to cherry-pick
git log
# cherry-pick the commits starting from the oldest one, without including merge commits
git cherry-pick -s -x <commit-id>
git cherry-pick -s -x <commit-id>
...
```

### 2. Update the VERSION file on master

We don't want to stop contributions to master just because we are releasing.
So, after the release branch is up, we bump the VERSION to mark the start of
the "next" release.

#### 2.1 Update the VERSION file

Update the content of the `VERSION` file to be the next minor (incrementing Y)
and add the `-dev` suffix. For example, after the release branch for 1.5.0 is
created, the `VERSION` file gets updated to `1.6.0-dev` (as in "1.6.0 in the
making").

### 3. Update CHANGELOG.md

If releasing a major version (X or Y increased in vX.Y.Z), simply listing notable user-facing features is sufficient.
```markdown
#### Notable features since <last major version>
* New docker command to do something useful
* Remote API change (deprecating old version)
* Performance improvements in some usecases
* ...
```

For minor releases (only Z increases in vX.Y.Z), provide a list of user-facing changes.
Each change should be listed under a category heading formatted as `#### CATEGORY`.

`CATEGORY` should describe which part of the project is affected.
  Valid categories are:
  * Documentation
  * CLI
  * API
  * Platform
  * Other (please use this category sparingly)

Each change should be formatted as `BULLET DESCRIPTION`, given:

* BULLET: either `-`, `+` or `*`, to indicate a bugfix, new feature or
  upgrade, respectively.

* DESCRIPTION: a concise description of the change that is relevant to the
  end-user, using the present tense. Changes should be described in terms
  of how they affect the user, for example "Add new feature X which allows Y",
  "Fix bug which caused X", "Increase performance of Y".

EXAMPLES:

```markdown
## 0.3.6 (1995-12-25)

#### CLI

+ 'amp stack up' creates and starts a stack

#### Remote API

- Added new monitoring APIs

#### Platform

* Added reverse proxy service
```

If you need a list of contributors between the last major release and the
current bump branch, use something like:
```bash
git log --format='%aN <%aE>' v0.7.0...bump_v0.8.0 | sort -uf
```
Obviously, you'll need to adjust version numbers as necessary.  If you just need
a count, add a simple `| wc -l`.

### 4. Change the contents of the VERSION file

Before the big thing, you'll want to make successive release candidates and get
people to test. The release candidate number `N` should be part of the version:

```bash
export RC_VERSION=${VERSION}-rcN
echo ${RC_VERSION#v} > VERSION
```

### 5. Commit and create a pull request to the "release" branch

```bash
git add VERSION CHANGELOG.md
git commit -m "Bump version to $VERSION"
git push $GITHUBUSER bump_$VERSION
echo "https://github.com/$GITHUBUSER/amp/compare/amp:release/$BASE...$GITHUBUSER:bump_$VERSION?expand=1"
```

That last command will give you the proper link to visit to ensure that you
open the PR against the "release" branch instead of accidentally against
"master" (like so many brave souls before you already have).

### 6. Create a PR to update the AUTHORS file for the release

Update the AUTHORS file.

### 7. Publish artifacts

```bash
aws --profile $PROFILE s3 --region=$REGION cp --acl public-read $BINDIR/$sos/$sarch/$ARCHIVE.tgz s3://$BUCKET/$S3DIR/$dos/$darch/${ARCHIVE}-v$VERSION.tgz
```


If everything looks good, it's time to create a git tag for this candidate:

```bash
git tag -a $RC_VERSION -m $RC_VERSION bump_$VERSION
git push origin $RC_VERSION
```

Announcing on multiple medias is the best way to get some help testing! An easy
way to get some useful links for sharing:

```bash
echo "Linux 64bit binary: https://test.appcelerator.io/builds/Linux/x86_64/amp-${VERSION#v}"
echo "Darwin/OSX 64bit client binary: https://test.appcelerator.io/builds/Darwin/x86_64/amp-${VERSION#v}"
```

We recommend announcing the release candidate on:

- IRC on #amp-cloud, #amp-cloud-dev, #amp-cloud-maintainers
- In a comment on the pull request to notify subscribed people on GitHub
- Any social media that can bring some attention to the release candidate

### 13. Iterate on successive release candidates

Spend several days along with the community explicitly investing time and
resources to try and break AMP in every possible way, documenting any
findings pertinent to the release.  This time should be spent testing and
finding ways in which the release might have caused various features or upgrade
environments to have issues, not coding.  During this time, the release is in
code freeze, and any additional code changes will be pushed out to the next
release.

It should include various levels of breaking AMP, beyond just using AMP
by the book.

Any issues found may still remain issues for this release, but they should be
documented and give appropriate warnings.

During this phase, the `bump_$VERSION` branch will keep evolving as you will
produce new release candidates. The frequency of new candidates is up to the
release manager: use your best judgement taking into account the severity of
reported issues, testers availability, and time to scheduled release date.

Each time you'll want to produce a new release candidate, you will start by
adding commits to the branch, usually by cherry-picking from master:

```bash
git cherry-pick -s -x -m0 <commit_id>
```

You want your "bump commit" (the one that updates the CHANGELOG and VERSION
files) to remain on top, so you'll have to `git rebase -i` to bring it back up.

Now that your bump commit is back on top, you will need to update the CHANGELOG
file (if appropriate for this particular release candidate), and update the
VERSION file to increment the RC number:

```bash
export RC_VERSION=$VERSION-rcN
echo $RC_VERSION > VERSION
```

You can now amend your last commit and update the bump branch:

```bash
git commit --amend
git push -f $GITHUBUSER bump_$VERSION
```

Repeat step 6 to tag the code, publish new binaries, announce availability, and
get help testing.

### 14. Finalize the bump branch

When you're happy with the quality of a release candidate, you can move on and
create the real thing.

You will first have to amend the "bump commit" to drop the release candidate
suffix in the VERSION file:

```bash
echo $VERSION > VERSION
git add VERSION
git commit --amend
```

You will then repeat step 6 to publish the binaries to test

### 15. Get 2 other maintainers to validate the pull request


### 19. Publish final binaries

Once they're tested and reasonably believed to be working, run against
get.appcelerator.io:

```bash
aws --profile $PROFILE s3 --region=$REGION cp --acl public-read $BINDIR/$sos/$sarch/$ARCHIVE.tgz s3://$BUCKET/$S3DIR/$dos/$darch/${ARCHIVE}-v$VERSION.tgz
```

### 21. Apply tag and create release

It's very important that we don't make the tag until after the official
release is uploaded to get.appcelerator.io!

```bash
git tag -a $VERSION -m $VERSION bump_$VERSION
git push origin $VERSION
```

Once the tag is pushed, go to GitHub and create a [new release](https://github.com/appcelerator/amp/releases/new).
If the tag is for an RC make sure you check `This is a pre-release` at the bottom of the form.

Select the tag that you just pushed as the version and paste the changelog in the description of the release.
You can see examples in this two links:

https://github.com/appcelerator/amp/releases/tag/v1.8.0
https://github.com/appcelerator/amp/releases/tag/v1.8.0-rc3

### 22. Go to github to merge the `bump_$VERSION` branch into release

Don't forget to push that pretty blue button to delete the leftover
branch afterwards!

### 23. Update the docs branch

You will need to point the docs branch to the newly created release tag:

```bash
git checkout origin/docs
git reset --hard origin/$VERSION
git push -f origin docs
```

The docs will appear on https://docs.appcelerator.io/ (though there may be cached
versions, so its worth checking http://docs.appcelerator.io.s3-website-us-east-1.amazonaws.com/).
For more information about documentation releases, see `docs/README.md`.

Note that the new docs will not appear live on the site until the cache (a complex,
distributed CDN system) is flushed. The `make docs-release` command will do this
_if_ the `DISTRIBUTION_ID` is set correctly - this will take at least 15 minutes to run
and you can check its progress with the CDN Cloudfront Chrome addon.

### 24. Create a new pull request to merge your bump commit back into master

```bash
git checkout master
git fetch
git reset --hard origin/master
git cherry-pick -s -x $VERSION
git push $GITHUBUSER merge_release_$VERSION
echo "https://github.com/$GITHUBUSER/appcelerator/compare/amp:master...$GITHUBUSER:merge_release_$VERSION?expand=1"
```

Again, get two maintainers to validate, then merge, then push that pretty
blue button to delete your branch.

### 25. Rejoice and Evangelize!

Congratulations! You're done.

Go forth and announce the glad tidings of the new release in `#amp-cloud`,
`#amp-cloud-dev`, and on Twitter!

