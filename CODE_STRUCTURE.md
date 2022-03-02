# Code Structure

## Go Directories

### `/cmd`

Main applications for this project.

The directory name for each application should match the name of the executable
you want to have (e.g., `/cmd/myapp`).

Don't put a lot of code in the application directory. If you think the code can
be imported and used in other projects, then it should live in the `/pkg` directory.
If the code is not reusable or if you don't want others to reuse it, put that code
in the `/internal` directory. You'll be surprised what others will do, so be
explicit about your intentions!

It's common to have a small `main` function that imports and invokes the code
from the `/internal` and `/pkg` directories and nothing else.


### `/internal`

Private application and library code. This is the code you don't want others
importing in their applications or libraries. Note that this layout pattern is
enforced by the Go compiler itself. See the Go 1.4 [`release notes`](https://golang.org/doc/go1.4#internalpackages)
for more details. Note that you are not limited to the top level `internal`
directory. You can have more than one `internal` directory at any level of your
project tree.

You can optionally add a bit of extra structure to your internal packages to
separate your shared and non-shared internal code. It's not required (especially
for smaller projects), but it's nice to have visual clues showing the intended
package use. Your actual application code can go in the `/internal/app` directory
(e.g., `/internal/app/myapp`) and the code shared by those apps in the
`/internal/pkg` directory (e.g., `/internal/pkg/myprivlib`).

### `/pkg`

Library code that's ok to use by external applications (e.g., `/pkg/mypubliclib`).
Other projects will import these libraries expecting them to work, so think twice
before you put something here :-) Note that the `internal` directory is a better
way to ensure your private packages are not importable because it's enforced by Go.
The `/pkg` directory is still a good way to explicitly communicate that the code
in that directory is safe for use by others. The [`I'll take pkg over internal`](https://travisjeffery.com/b/2019/11/i-ll-take-pkg-over-internal/)
blog post by Travis Jeffery provides a good overview of the `pkg` and `internal`
directories and when it might make sense to use them.

It's also a way to group Go code in one place when your root directory contains
lots of non-Go components and directories making it easier to run various Go tools
(as mentioned in these talks: [`Best Practices for Industrial Programming`](https://www.youtube.com/watch?v=PTE4VJIdHPg)
from GopherCon EU 2018, [GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)
and [GoLab 2018 - Massimiliano Pippi - Project layout patterns in Go](https://www.youtube.com/watch?v=3gQa1LWwuzk)).

This is a common layout pattern, but it's not universally accepted and some in the
Go community don't recommend it.

It's ok not to use it if your app project is really small and where an extra level
of nesting doesn't add much value (unless you really want to :-)). Think about it
when it's getting big enough and your root directory gets pretty busy (especially
if you have a lot of non-Go app components).

### `/vendor`

Application dependencies (managed manually or by your favorite dependency management
tool like the new built-in [`Go Modules`](https://github.com/golang/go/wiki/Modules)
feature). The `go mod vendor` command will create the `/vendor` directory for you.
Note that you might need to add the `-mod=vendor` flag to your `go build` command if
you are not using Go 1.14 where it's on by default.

Don't commit your application dependencies if you are building a library.

Note that since [`1.13`](https://golang.org/doc/go1.13#modules) Go also enabled the
module proxy feature (using [`https://proxy.golang.org`](https://proxy.golang.org)
as their module proxy server by default). Read more about it [`here`](https://blog.golang.org/module-mirror-launch)
to see if it fits all of your requirements and constraints. If it does, then you
won't need the `vendor` directory at all.

## Service Application Directories

### `/api`

OpenAPI/Swagger specs, JSON schema files, protocol definition files.

## Web Application Directories

### `/web`

Web application specific components: static web assets, server side templates and SPAs.

## Common Application Directories

### `/configs`

Configuration file templates or default configs.

### `/init`

System init (systemd, upstart, sysv) and process manager/supervisor (runit, supervisord)
configs.

### `/scripts`

Scripts to perform various build, install, analysis, etc operations.

These scripts keep the root level Makefile small and simple (e.g.,
[`https://github.com/hashicorp/terraform/blob/master/Makefile`](https://github.com/hashicorp/terraform/blob/master/Makefile)).

### `/build`

Packaging and Continuous Integration.

Put your cloud (AMI), container (Docker), OS (deb, rpm, pkg) package configurations
and scripts in the `/build/package` directory.

Put your CI (travis, circle, drone) configurations and scripts in the `/build/ci`
directory. Note that some of the CI tools (e.g., Travis CI) are very picky about
the location of their config files. Try putting the config files in the `/build/ci`
directory linking them to the location where the CI tools expect them (when possible).

### `/deployments`

IaaS, PaaS, system and container orchestration deployment configurations and
templates (docker-compose, kubernetes/helm, mesos, terraform, bosh). Note that
in some repos (especially apps deployed with kubernetes) this directory is
called `/deploy`.

### `/test`

Additional external test apps and test data. Feel free to structure the `/test`
directory anyway you want. For bigger projects it makes sense to have a data
subdirectory. For example, you can have `/test/data` or `/test/testdata` if you
need Go to ignore what's in that directory. Note that Go will also ignore
directories or files that begin with "." or "\_", so you have more flexibility
in terms of how you name your test data directory.

## Other Directories

### `/docs`

Design and user documents (in addition to your godoc generated documentation).

### `/tools`

Supporting tools for this project. Note that these tools can import code from
the `/pkg` and `/internal` directories.

### `/examples`

Examples for your applications and/or public libraries.

### `/third_party`

External helper tools, forked code and other 3rd party utilities (e.g., Swagger UI).

### `/githooks`

Git hooks.

### `/assets`

Other assets to go along with your repository (images, logos, etc).

### `/website`

This is the place to put your project's website data if you are not using GitHub
pages.

## Directories You Shouldn't Have

### `/src`

Some Go projects do have a `src` folder, but it usually happens when the devs
came from the Java world where it's a common pattern. If you can help yourself
try not to adopt this Java pattern. You really don't want your Go code or Go
projects to look like Java :-)

Don't confuse the project level `/src` directory with the `/src` directory Go
uses for its workspaces as described in [`How to Write Go Code`](https://golang.org/doc/code.html).
The `$GOPATH` environment variable points to your (current) workspace (by default
it points to `$HOME/go` on non-windows systems). This workspace includes the top
level `/pkg`, `/bin` and `/src` directories. Your actual project ends up being
a sub-directory under `/src`, so if you have the `/src` directory in your project
the project path will look like this: `/some/path/to/workspace/src/your_project/src/your_code.go`.
Note that with Go 1.11 it's possible to have your project outside of your `GOPATH`,
but it still doesn't mean it's a good idea to use this layout pattern.