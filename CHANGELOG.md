v0.2.0 2016-04-10

CHANGELOG
- 2016-04-10, Version 0.2.0:
* Add new feature `gnvm session` Use any Node.js version of the local already exists version by current session.
* Add new feature `gnvm npm`     NPM version management.
* Add new feature `gnvm search`  Search and Print Node.js version detail usage wildcard mode or regexp mode.
* Add new feature `gnvm reg`     Add config property `noderoot` to Environment variable `NODE_HOME`.
* Add new `gnvm` icon, inlcude:  `32*32`, `64*64`, `128*128` size.
* Senior  feature `gnvm ls`      Add detail print, e.g. `gnvm ls -r -d -l`.
* Senror  feature `gnvm config`  Add print all property`gnvm config` and test custom registry `gnvm config registry test`.
* Senior  feature `gnvm config`  built-in `DEFAULT` <http://nodejs.org/dist/> and `TAOBAO` <http://npm.taobao.org/mirrors/node> keywords.
* Senior  feature `gnvm install` Auto support iojs, e.g. `gnvm install 1.0.0` and arch, e.g. `gnvm install 5.1.1-x64 5.1.1-x86`.
* Senior  feature `gnvm version` Add new flag `-r`, print `CHANGELOG`, e.g. `gnvm version -r -d`.

-2016-02-24, Version 0.1.4 beta:
* Fix <http://nodejs.org/dist> changes the structure errors, including: `gnvm install`, `gnvm ls -r`, `gnvm config INIT`.
* Add New registry <http://npm.taobao.org/mirrors/node> use `gnvm config registry TAOBAO`.

-2014-07-23, Version 0.1.3:
* Fixbug of `node.exe` process to take up,  `gnvm use x.xx.xx` not work.
* When usage `gnvm use x.xx.xx`, kill node.exe process automatically.

-2014-07-15, Version 0.1.2:
* Adapter `go version 1.3`.
* Fix bug of usage `gnvm update latest -g` adapter go version error.

-2014-06-06, Version 0.1.1:
* Change `util/p/print.go`   to `github.com/Kenshin/cprint`.
* Change `util/curl/curl.go` to `github.com/Kenshin/curl`.
* Add this project to travis-ci.org.
* Remove `nodehandle.cmd` method.
* Optimize `nodehandle.copy` method logic.
* Fix bug of When not global `node.exe`, need get `gnvm.exe` path.

-2014-05-30, Version 0.1.0:
* `gnvm version`         Print GNVM version number.
* `gnvm install`         Install any Node.js version.
* `gnvm uninstall`       Uninstall local Node.js version and npm.
* `gnvm use`             Use any the local already exists of Node.js version.
* `gnvm update`          Update Node.js latest version.
* `gnvm ls`              Show all [local] [remote] Node.js version.
* `gnvm node-version`    Show [global] [latest] Node.js version.
* `gnvm config`          Setter and getter .gnvmrc file.