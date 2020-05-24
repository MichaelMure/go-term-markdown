    ![]([34mhttps://cdn.rawgit.com/MichaelMu[0m
    [34mre/git-bug/master/misc/logo/logo-alp[0m
    [34m          ha-flat-bg.svg[0m)

    ╔══════════════════════════════════╗
    ║             git-bug              ║
    ╚══════════════════════════════════╝

    [Build Status][1]
    [Backers on Open Collective][2]
    [Sponsors on Open Collective][3]
    [License: GPL v3][4]
    [GoDoc][5]
    [Go Report Card][6]
    [Gitter chat][7]

    [1]: https://travis-ci.org/MichaelMu
    re/git-bug
    [2]: #backers
    [3]: #sponsors
    [4]:
    http://www.gnu.org/licenses/gpl-3.0
    [5]: https://godoc.org/github.com/Mi
    chaelMure/git-bug
    [6]: https://goreportcard.com/report
    /github.com/MichaelMure/git-bug
    [7]:
    https://gitter.im/the-git-bug/Lobby

    [44;3mgit-bug[0m is a bug tracker that:

    [32m• [0m[1mfully embed in git[0m: you only need
      your git repository to have a bug
      tracker
    [32m• [0m[1mis distributed[0m: use your normal
      git remote to collaborate, push
      and pull your bugs !
    [32m• [0m[1mworks offline[0m: in a plane or under
      the sea ? keep reading and writing
      bugs
    [32m• [0m[1mprevent vendor locking[0m: your usual
      service is down or went bad ? you
      already have a full backup
    [32m• [0m[1mis fast[0m: listing bugs or opening
      them is a matter of milliseconds
    [32m• [0m[1mdoesn't pollute your project[0m: no
      files are added in your project
    [32m• [0m[1mintegrate with your tooling[0m: use
      the UI you like (CLI, terminal,
      web) or integrate with your
      existing tools through the CLI or
      the GraphQL API
    [32m• [0m[1mbridge with other bug trackers[0m:
      [bridges][8] exist to import and
      soon export to other trackers.

    [8]: #bridges

    🚧  This is now more than a proof of
    concept, but still not fully stable.
    Expect dragons and unfinished
    business. 🚧

    [32;1m1 Install[0m

    [31m<details><summary>go[0m
    [31mget</summary></details>[0m[31m<summary>go[0m
    [31mget</summary>[0mgo get[31m<summary>go[0m
    [31mget</summary>[0m[31m<details><summary>go[0m
    [31mget</summary></details>[0m

    [32;1m┃ [0mgo get -u
    [32;1m┃ [0mgithub.com/MichaelMure/git-bug

    If it's not done already, add golang
    binary directory in your PATH:

    [32;1m┃ [0m[32mexport[0m [34mPATH[0m[90m=[0m[34m$PATH[0m:[1m[32m$([0mgo env
    [32;1m┃ [0mGOROOT[1m[32m)[0m/bin:[1m[32m$([0mgo env GOPATH[1m[32m)[0m/bin

    [31m<details><summary>Pre-compiled binar[0m
    [31mies</summary></details>[0m[31m<summary>Pre-[0m
    [31mcompiled[0m
    [31mbinaries</summary>[0mPre-compiled
    binaries[31m<summary>Pre-compiled binari[0m
    [31mes</summary>[0m[31m<details><summary>Pre-co[0m
    [31mmpiled binaries</summary></details>[0m

    [32m1. [0mGo to the [release page][9] and
       download the appropriate binary
       for your system.
    [32m2. [0mCopy the binary anywhere in your
       PATH
    [32m3. [0mRename the binary to [44;3mgit-bug[0m (or
       [44;3mgit-bug.exe[0m on windows)

    [9]: https://github.com/MichaelMure/
    git-bug/releases/latest

    That's all !

    [31m<details><summary>Linux packages</su[0m
    [31mmmary></details>[0m[31m<summary>Linux[0m
    [31mpackages</summary>[0mLinux
    packages[31m<summary>Linux packages</sum[0m
    [31mmary>[0m[31m<details><summary>Linux[0m
    [31mpackages</summary></details>[0m

    [32m• [0m[Archlinux (AUR)][10]

    [10]: https://aur.archlinux.org/pack
    ages/?K=git-bug

    [32;1m2 CLI usage[0m

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
    [query][11]:

    [11]: doc/queries.md

    [32;1m┃ [0mgit bug ls "status:open sort:edit"

    You can now use commands like [44;3mshow[0m,
    [44;3mcomment[0m, [44;3mopen[0m or [44;3mclose[0m to display
    and modify bugs. For more details
    about each command, you can run [44;3mgit[0m
    [3;44mbug <command> --help[0m or read the
    [command's documentation][12].

    [12]: doc/md/git-bug.md

    [32;1m3 Interactive terminal UI[0m

    An interactive terminal UI is
    available using the command [44;3mgit bug[0m
    [3;44mtermui[0m to browse and edit bugs.

    Termui recording

    [32;1m4 Web UI (status: WIP)[0m

    You can launch a rich Web UI with
    [44;3mgit bug webui[0m.

    Web UI screenshot 1
    Web UI screenshot 2

    This web UI is entirely packed
    inside the same go binary and serve
    static content through a localhost
    http server.

    The web UI interact with the backend
    through a GraphQL API. The schema is
    available [here][13].

    [13]: graphql/

    [32;1m5 Bridges[0m

    [92m5.1 Importer implementations[0m

    ┌──────────┬───────────┬───────────┐
    │[0m[0m          │[0mGithub[0m     │ [0mLaunchpad[0m │
    ╞══════════╪═══════════╪═══════════╡
    │[0m[1mincrementa[0m│[0m:heavy_chec[0m│    [0m:x:[0m    │
    │[1ml[0m[31m<br/>[0m(can[0m│[0mk_mark:[0m    │           │
    │[0mimport[0m    │           │           │
    │[0mmore than[0m │           │           │
    │[0monce)[0m     │           │           │
    ├──────────┼───────────┼───────────┤
    │[0m[1mwith resum[0m│[0m:x:[0m        │    [0m:x:[0m    │
    │[1me[0m[31m<br/>[0m(dow[0m│           │           │
    │[0mnload only[0m│           │           │
    │[0mnew data)[0m │           │           │
    ├──────────┼───────────┼───────────┤
    │[0m[1midentities[0m[0m│[0m:heavy_chec[0m│[0m:heavy_chec[0m│
    │          │[0mk_mark:[0m    │  [0mk_mark:[0m  │
    ├──────────┼───────────┼───────────┤
    │[0midentities[0m│[0m:x:[0m        │    [0m:x:[0m    │
    │[0mupdate[0m    │           │           │
    ├──────────┼───────────┼───────────┤
    │[0m[1mbug[0m[0m       │[0m:heavy_chec[0m│[0m:heavy_chec[0m│
    │          │[0mk_mark:[0m    │  [0mk_mark:[0m  │
    ├──────────┼───────────┼───────────┤
    │[0mcomments[0m  │[0m:heavy_chec[0m│[0m:heavy_chec[0m│
    │          │[0mk_mark:[0m    │  [0mk_mark:[0m  │
    ├──────────┼───────────┼───────────┤
    │[0mcomment[0m   │[0m:heavy_chec[0m│    [0m:x:[0m    │
    │[0meditions[0m  │[0mk_mark:[0m    │           │
    ├──────────┼───────────┼───────────┤
    │[0mlabels[0m    │[0m:heavy_chec[0m│    [0m:x:[0m    │
    │          │[0mk_mark:[0m    │           │
    ├──────────┼───────────┼───────────┤
    │[0mstatus[0m    │[0m:heavy_chec[0m│    [0m:x:[0m    │
    │          │[0mk_mark:[0m    │           │
    ├──────────┼───────────┼───────────┤
    │[0mtitle[0m     │[0m:heavy_chec[0m│    [0m:x:[0m    │
    │[0medition[0m   │[0mk_mark:[0m    │           │
    ├──────────┼───────────┼───────────┤
    │[0m[1mautomated[0m │[0m:x:[0m        │    [0m:x:[0m    │
    │[1mtest suite[0m[0m│           │           │
    └──────────┴───────────┴───────────┘
    [92m5.2 Exporter implementations[0m

    Todo !

    [32;1m6 Internals[0m

    Interested by how it works ? Have a
    look at the [data model][14] and the
    [internal bird-view][15].

    [14]: doc/model.md
    [15]: doc/architecture.md

    [32;1m7 Misc[0m

    [32m• [0m[Bash completion][16]
    [32m• [0m[Zsh completion][17]
    [32m• [0m[ManPages][18]

    [16]: ../misc/bash_completion
    [17]: ../misc/zsh_completion
    [18]: doc/man

    [32;1m8 Planned features[0m

    [32m• [0mmedia embedding
    [32m• [0mexporter to github issue
    [32m• [0mextendable data model to support
      arbitrary bug tracker
    [32m• [0minflatable raptor

    [32;1m9 Contribute[0m

    PRs accepted. Drop by the [Gitter
    lobby][19] for a chat or browse the
    issues to see what is worked on or
    discussed.

    [19]:
    https://gitter.im/the-git-bug/Lobby

    Developers unfamiliar with Go may
    try to clone the repository using
    "git clone". Instead, one should
    use:

    [32;1m┃ [0mgo get -u
    [32;1m┃ [0mgithub.com/MichaelMure/git-bug

    The git repository will then be
    available:

    [32;1m┃ [0m[3m[36m# Note that $GOPATH defaults to[0m
    [32;1m┃ [0m[3;36m$HOME/go[0m
    [32;1m┃ [0m$ [32mcd[0m [34m$GOPATH[0m/src/github.com/Michae
    [32;1m┃ [0mlMure/git-bug/

    You can now run [44;3mmake[0m to build the
    project, or [44;3mmake install[0m to install
    the binary in [44;3m$GOPATH/bin/[0m.

    To work on the web UI, have a look
    at [the dedicated Readme.][20]

    [20]: webui/Readme.md

    [32;1m10 Contributors :heart:[0m

    This project exists thanks to all
    the people who contribute.
    [31m<a href="https://github.com/MichaelM[0m
    [31mure/git-bug/graphs/contributors">[0m[31m<im[0m
    [31mg src="https://opencollective.com/gi[0m
    [31mt-bug/contributors.svg?width=890&but[0m
    [31mton=false" />[0m[31m</a>[0m

    [32;1m11 Backers[0m

    Thank you to all our backers! 🙏
    [[Become a backer][21]]

    [21]: https://opencollective.com/git
    -bug#backer

    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug#backers"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/backe[0m
    [31mr.svg?width=890">[0m[31m</a>[0m

    [32;1m12 Sponsors[0m

    Support this project by becoming a
    sponsor. Your logo will show up here
    with a link to your website.
    [[Become a sponsor][22]]

    [22]: https://opencollective.com/git
    -bug#sponsor

    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/0/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/0/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/1/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/1/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/2/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/2/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/3/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/3/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/4/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/4/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/5/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/5/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/6/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/6/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/7/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/7/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/8/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/8/avatar.svg">[0m[31m</a>[0m
    [31m<a href="https://opencollective.com/[0m
    [31mgit-bug/sponsor/9/website"[0m
    [31mtarget="_blank">[0m[31m<img src="https://op[0m
    [31mencollective.com/git-bug/tiers/spons[0m
    [31mor/9/avatar.svg">[0m[31m</a>[0m

    [32;1m13 License[0m

    Unless otherwise stated, this
    project is released under the
    [GPLv3][23] or later license ©
    Michael Muré.

    [23]: LICENSE

    The git-bug logo by [Viktor
    Teplov][24] is released under the
    [Creative Commons Attribution 4.0
    International (CC BY 4.0)][25]
    license © Viktor Teplov.

    [24]: https://github.com/vandesign
    [25]: ../misc/logo/LICENSE

