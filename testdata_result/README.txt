    [31m<p align="center">[0m[31m<img width="150px"
    src="https://cdn.rawgit.com/MichaelM
    ure/git-bug/master/misc/logo/logo-al
    pha-flat-bg.svg">[0m[31m</p>[0m

    [32;1m1 git-bug[0m
    ────────────────────────────────────

    [![Build Status](https://travis-ci.o
    rg/MichaelMure/git-bug.svg?branch=ma
    ster)](https://travis-ci.org/Michael
    Mure/git-bug)[![Backers on Open Coll
    ective](https://opencollective.com/g
    it-bug/backers/badge.svg)](#backers)
    [![Sponsors on Open Collective](htt
    ps://opencollective.com/git-bug/spon
        sors/badge.svg)](#sponsors)
    [![License: GPL v3](https://img.shie
    lds.io/badge/License-GPLv3+-blue.svg
    )](http://www.gnu.org/licenses/gpl-3
    .0)[![GoDoc](https://godoc.org/githu
    b.com/MichaelMure/git-bug?status.svg
    )](https://godoc.org/github.com/Mich
    aelMure/git-bug)[![Go Report Card](h
    ttps://goreportcard.com/badge/github
    .com/MichaelMure/git-bug)](https://g
    oreportcard.com/report/github.com/Mi
    chaelMure/git-bug)[![Gitter chat](ht
    tps://badges.gitter.im/gitterHQ/gitt
    er.png)](https://gitter.im/the-git-b
                 ug/Lobby)

    [44;3mgit-bug[0m is a bug tracker that:
    [32m• [0m[1mfully embed in git[21m: you only need
      your git repository to have a bug
      tracker
    [32m• [0m[1mis distributed[21m: use your normal
      git remote to collaborate, push
      and pull your bugs !
    [32m• [0m[1mworks offline[21m: in a plane or under
      the sea ? keep reading and writing
      bugs
    [32m• [0m[1mprevent vendor locking[21m: your usual
      service is down or went bad ? you
      already have a full backup
    [32m• [0m[1mis fast[21m: listing bugs or opening
      them is a matter of milliseconds
    [32m• [0m[1mdoesn't pollute your project[21m: no
      files are added in your project
    [32m• [0m[1mintegrate with your tooling[21m: use
      the UI you like (CLI, terminal,
      web) or integrate with your
      existing tools through the CLI or
      the GraphQL API
    [32m• [0m[1mbridge with other bug trackers[21m:
      [bridges]([34m#bridges[0m) exist to
      import and soon export to other
      trackers.

    🚧  This is now more than a proof of
    concept, but still not fully stable.
    Expect dragons and unfinished
    business. 🚧

    [32;1m1.1 Install[0m

    [31m<details>[0m[31m<summary>[0mgo get[31m</summary>[0m

    [32;1m┃ [0mgo get -u
    [32;1m┃ [0mgithub.com/MichaelMure/git-bug

    If it's not done already, add golang
    binary directory in your PATH:

    [32;1m┃ [0m[32mexport[0m [34mPATH[0m[90m=[0m[34m$PATH[0m:[1m[32m$([0mgo env
    [32;1m┃ [0mGOROOT[1m[32m)[0m/bin:[1m[32m$([0mgo env GOPATH[1m[32m)[0m/bin

    [31m</details>[0m

    [31m<details>[0m[31m<summary>[0mPre-compiled
    binaries[31m</summary>[0m
    [32m1. [0mGo to the [release page]([34mhttps://
       github.com/MichaelMure/git-bug/re
       leases/latest[0m) and download the
       appropriate binary for your
       system.
    [32m2. [0mCopy the binary anywhere in your
       PATH
    [32m3. [0mRename the binary to [44;3mgit-bug[0m (or
       [44;3mgit-bug.exe[0m on windows)

    That's all !

    [31m</details>[0m

    [31m<details>[0m[31m<summary>[0mLinux
    packages[31m</summary>[0m
    [32m• [0m[Archlinux (AUR)]([34mhttps://aur.arch
      linux.org/packages/?K=git-bug[0m)

    [31m</details>[0m

    [32;1m1.2 CLI usage[0m

    Create a new identity:

    [32;1m┃ [0mgit bug user create

    Create a new bug:

    [32;1m┃ [0mgit bug add

    Your favorite editor will open to
    write a title and a message.

    You can push your new entry to a
    remote:

    [32;1m┃ [0mgit bug push [<remote>]

    And pull for updates:

    [32;1m┃ [0mgit bug pull [<remote>]

    List existing bugs:

    [32;1m┃ [0mgit bug ls

    Filter and sort bugs using a
    [query]([34mdoc/queries.md[0m):

    [32;1m┃ [0mgit bug ls "status:open sort:edit"

    You can now use commands like [44;3mshow[0m,
    [44;3mcomment[0m, [44;3mopen[0m or [44;3mclose[0m to display
    and modify bugs. For more details
    about each command, you can run [44;3mgit
    bug <command> --help[0m or read the
    [command's
    documentation]([34mdoc/md/git-bug.md[0m).

    [32;1m1.3 Interactive terminal UI[0m

    An interactive terminal UI is
    available using the command [44;3mgit bug
    termui[0m to browse and edit bugs.

    Termui recording

    [32;1m1.4 Web UI (status: WIP)[0m

    You can launch a rich Web UI with
    [44;3mgit bug webui[0m.

    Web UI screenshot 1 Web UI
    screenshot 2

    This web UI is entirely packed
    inside the same go binary and serve
    static content through a localhost
    http server.

    The web UI interact with the backend
    through a GraphQL API. The schema is
    available [here]([34mgraphql/[0m).

    [32;1m1.5 Bridges[0m

    [92m1.5.1 Importer implementations[0m

    ┌─────────────────┬──────┬─────────┐
    │[0m                 │Github[0m│Launchpad[0m│
    ╞═════════════════╪══════╪═════════╡
    │[1mincremental[21m[31m<br/>[0m([0m│✔[0m    │   ❌[0m    │
    │[1m[21m[31m[0mcan import more[0m  │      │         │
    │[1m[21m[31m[0mthan once)[0m       │      │         │
    ├─────────────────┼──────┼─────────┤
    │[1mwith resume[21m[31m<br/>[0m([0m│❌[0m    │   ❌[0m    │
    │[1m[21m[31m[0mdownload only new[0m│      │         │
    │[1m[21m[31m[0mdata)[0m            │      │         │
    ├─────────────────┼──────┼─────────┤
    │[1midentities[21m[0m       │✔[0m    │   ✔[0m    │
    ├─────────────────┼──────┼─────────┤
    │identities update[0m│❌[0m    │   ❌[0m    │
    ├─────────────────┼──────┼─────────┤
    │[1mbug[21m[0m              │✔[0m    │   ✔[0m    │
    ├─────────────────┼──────┼─────────┤
    │comments[0m         │✔[0m    │   ✔[0m    │
    ├─────────────────┼──────┼─────────┤
    │comment editions[0m │✔[0m    │   ❌[0m    │
    ├─────────────────┼──────┼─────────┤
    │labels[0m           │✔[0m    │   ❌[0m    │
    ├─────────────────┼──────┼─────────┤
    │status[0m           │✔[0m    │   ❌[0m    │
    ├─────────────────┼──────┼─────────┤
    │title edition[0m    │✔[0m    │   ❌[0m    │
    ├─────────────────┼──────┼─────────┤
    │[1mautomated test[0m   │❌[0m    │   ❌[0m    │
    │[1msuite[21m[0m            │      │         │
    └─────────────────┴──────┴─────────┘
    [92m1.5.2 Exporter implementations[0m

    Todo !

    [32;1m1.6 Internals[0m

    Interested by how it works ? Have a
    look at the [data
    model]([34mdoc/model.md[0m) and the
    [internal
    bird-view]([34mdoc/architecture.md[0m).

    [32;1m1.7 Misc[0m

    [32m• [0m[Bash
      completion]([34mmisc/bash_completion[0m)
    [32m• [0m[Zsh
      completion]([34mmisc/zsh_completion[0m)
    [32m• [0m[ManPages]([34mdoc/man[0m)

    [32;1m1.8 Planned features[0m

    [32m• [0mmedia embedding
    [32m• [0mexporter to github issue
    [32m• [0mextendable data model to support
      arbitrary bug tracker
    [32m• [0minflatable raptor

    [32;1m1.9 Contribute[0m

    PRs accepted. Drop by the [Gitter lo
    bby]([34mhttps://gitter.im/the-git-bug/L
    obby[0m) for a chat or browse the
    issues to see what is worked on or
    discussed.

    Developers unfamiliar with Go may
    try to clone the repository using
    "git clone". Instead, one should
    use:

    [32;1m┃ [0mgo get -u
    [32;1m┃ [0mgithub.com/MichaelMure/git-bug

    The git repository will then be
    available:

    [32;1m┃ [0m[36m# Note that $GOPATH defaults to
    [32;1m┃ [0m$HOME/go[0m
    [32;1m┃ [0m$ [32mcd[0m [34m$GOPATH[0m/src/github.com/Michae
    [32;1m┃ [0mlMure/git-bug/

    You can now run [44;3mmake[0m to build the
    project, or [44;3mmake install[0m to install
    the binary in [44;3m$GOPATH/bin/[0m.

    To work on the web UI, have a look
    at [the dedicated
    Readme.]([34mwebui/Readme.md[0m)

    [32;1m1.10 Contributors ❤ [0m

    This project exists thanks to all
    the people who contribute. [31m<a href="
    https://github.com/MichaelMure/git-b
    ug/graphs/contributors">[0m[31m<img src="ht
    tps://opencollective.com/git-bug/con
    tributors.svg?width=890&button=false
    " />[0m[31m</a>[0m

    [32;1m1.11 Backers[0m

    Thank you to all our backers! 🙏
    [[Become a backer]([34mhttps://opencolle
    ctive.com/git-bug#backer[0m)]

    [31m<a href="https://opencollective.com/
    git-bug#backers"
    target="_blank">[0m[31m<img src="https://op
    encollective.com/git-bug/tiers/backe
    r.svg?width=890">[0m[31m</a>[0m

    [32;1m1.12 Sponsors[0m

    Support this project by becoming a
    sponsor. Your logo will show up here
    with a link to your website.
    [[Become a sponsor]([34mhttps://opencoll
    ective.com/git-bug#sponsor[0m)]

    [31m<a href="https://opencollective.com/
    git-bug/sponsor/0/website"
    target="_blank">[0m[31m<img src="https://op
    encollective.com/git-bug/tiers/spons
    or/0/avatar.svg">[0m[31m</a>[0m [31m<a href="https
    ://opencollective.com/git-bug/sponso
    r/1/website" target="_blank">[0m[31m<img sr
    c="https://opencollective.com/git-bu
    g/tiers/sponsor/1/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/
    git-bug/sponsor/2/website"
    target="_blank">[0m[31m<img src="https://op
    encollective.com/git-bug/tiers/spons
    or/2/avatar.svg">[0m[31m</a>[0m [31m<a href="https
    ://opencollective.com/git-bug/sponso
    r/3/website" target="_blank">[0m[31m<img sr
    c="https://opencollective.com/git-bu
    g/tiers/sponsor/3/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/
    git-bug/sponsor/4/website"
    target="_blank">[0m[31m<img src="https://op
    encollective.com/git-bug/tiers/spons
    or/4/avatar.svg">[0m[31m</a>[0m [31m<a href="https
    ://opencollective.com/git-bug/sponso
    r/5/website" target="_blank">[0m[31m<img sr
    c="https://opencollective.com/git-bu
    g/tiers/sponsor/5/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/
    git-bug/sponsor/6/website"
    target="_blank">[0m[31m<img src="https://op
    encollective.com/git-bug/tiers/spons
    or/6/avatar.svg">[0m[31m</a>[0m [31m<a href="https
    ://opencollective.com/git-bug/sponso
    r/7/website" target="_blank">[0m[31m<img sr
    c="https://opencollective.com/git-bu
    g/tiers/sponsor/7/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/
    git-bug/sponsor/8/website"
    target="_blank">[0m[31m<img src="https://op
    encollective.com/git-bug/tiers/spons
    or/8/avatar.svg">[0m[31m</a>[0m [31m<a href="https
    ://opencollective.com/git-bug/sponso
    r/9/website" target="_blank">[0m[31m<img sr
    c="https://opencollective.com/git-bu
    g/tiers/sponsor/9/avatar.svg">[0m[31m</a>[0m

    [32;1m1.13 License[0m

    Unless otherwise stated, this
    project is released under the
    [GPLv3]([34mLICENSE[0m) or later license ©
    Michael Muré.

    The git-bug logo by [Viktor Teplov](
    [34mhttps://github.com/vandesign[0m) is
    released under the [Creative Commons
    Attribution 4.0 International (CC
    BY 4.0)]([34mmisc/logo/LICENSE[0m) license
    © Viktor Teplov.
