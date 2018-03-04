## Debanize 104
Here will show how to debanize an application implemented in cmd sub directory,
with two dependencies and with multiple source files.

### Go version
Actual host is Debian stretch, but with stretch-backport. 
```bash
go version
go version go1.8.1 linux/amd64
```

### Use dh-make-golang
Set environment
```bash
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$GOBIN
```
Use option type *program|library* since dh-make-golang can easily be confused by the files in the repository.
```bash
mkdir -p ~/PKGDEB/104
cd ~/PKGDEB/104
dh-make-golang -type program github.com/berrak/104
...
2018/03/04 14:42:33 Determining dependencies
2018/03/04 14:42:34 Build-Dependency "github.com/berrak/pkg103greet" is not yet available in Debian, or has not yet been converted to use XS-Go-Import-Path in debian/control
...
```
Notice that our home made package has to be taken care of during the packaging process.
Directory (104) and files after above run:
```bash
104
104_0.0~git20180304.eaf0f5b.orig.tar.xz
itp-104.txt
```
Binary application name is not changed by Debian, thus it remains as *104*.
Ignore the text file, its only for Debian internal usage. Lets look into **104** directory.
It contains three git branches:
```bash
git vbr
* master
  pristine-tar
  upstream
```
Switch to pristine-tar branch and confirm have same base name as the original tarball:
```bash
git sw2 pristine-tar
ls -l
104_0.0~git20180304.eaf0f5b.orig.tar.xz.delta
104_0.0~git20180304.eaf0f5b.orig.tar.xz.id


git sw2 master
```
Edit **changelog** in the debian directory.
```bash
cd debian
vi changelog
```
Replace *UNRELEASED* with stable and remove TODO with (Closes: #123456) or whatever.
Edit the file **control**. Add the Build dependency to the packaged version of *github.com/berrak/pkg103greet*:
Add the local debanized package *golang-github-berrak-pkg103greet-dev* like so:
```bash
...
Build-Depends: debhelper (>= 10),
               dh-golang,
               golang-any,
               golang-github-berrak-pkg103greet-dev,
               golang-github-bgentry-speakeasy-dev
...
```
**rules** does not need the DH_GOPKG stanza anymore, since control has *XS-Go-Import-Path: github.com/berrak/104*.
Thus, the default file should be ready to run for first fakeroot-tests:
```bash
#!/usr/bin/make -f

%:
        dh $@ --buildsystem=golang --with=golang    
```
Now we will try to build it first with fakeroot to identify problems with these files or the repository:
```bash
cd ..
fakeroot debian/rules build
dh build --buildsystem=golang --with=golang
   dh_testdir -O--buildsystem=golang
   dh_update_autotools_config -O--buildsystem=golang
   dh_autoreconf -O--buildsystem=golang
   dh_auto_configure -O--buildsystem=golang
   dh_auto_build -O--buildsystem=golang
	go install -v -p 1 github.com/berrak/104/cmd/104
src/github.com/berrak/104/cmd/104/104.go:25:2: cannot find package "github.com/bgentry/speakeasy" in any of:
	/usr/lib/go-1.8/src/github.com/bgentry/speakeasy (from $GOROOT)
	/home/bekr/PKGDEB/104/104/obj-x86_64-linux-gnu/src/github.com/bgentry/speakeasy (from $GOPATH)
dh_auto_build: go install -v -p 1 github.com/berrak/104/cmd/104 returned exit code 1
debian/rules:4: recipe for target 'build' failed
make: * * * [build] Error 1
```
Clearly missing the installation of the last dependency. Let's install it.
In an previous chapter, *golang-github-berrak-pkg103greet-dev* is already installed. 
```bash
sudo apt-get install golang-github-bgentry-speakeasy-dev
```
Remove *obj-x86_64-linux-gnu* and rerun fakeroot:
Let's build again:
```bash
fakeroot debian/rules build
dh build --buildsystem=golang --with=golang
   dh_testdir -O--buildsystem=golang
   dh_update_autotools_config -O--buildsystem=golang
   dh_autoreconf -O--buildsystem=golang
   dh_auto_configure -O--buildsystem=golang
   dh_auto_build -O--buildsystem=golang
	go install -v -p 1 github.com/berrak/104/cmd/104
github.com/berrak/pkg103greet
github.com/bgentry/speakeasy
github.com/berrak/104/cmd/104
   dh_auto_test -O--buildsystem=golang
	go test -v -p 1 github.com/berrak/104/cmd/104
=== RUN   TestOne
--- PASS: TestOne (0.00s)
PASS
ok  	github.com/berrak/104/cmd/104	0.001s
   create-stamp debian/debhelper-build-stamp
```
Nice, lets try to build the binary. Clean up and then re-run:
```bash
fakeroot debian/rules binary
dh binary --buildsystem=golang --with=golang
   dh_testdir -O--buildsystem=golang
   dh_update_autotools_config -O--buildsystem=golang
   dh_autoreconf -O--buildsystem=golang
   dh_auto_configure -O--buildsystem=golang
   dh_auto_build -O--buildsystem=golang
	go install -v -p 1 github.com/berrak/104/cmd/104
github.com/berrak/pkg103greet
github.com/bgentry/speakeasy
github.com/berrak/104/cmd/104
   dh_auto_test -O--buildsystem=golang
	go test -v -p 1 github.com/berrak/104/cmd/104
=== RUN   TestOne
--- PASS: TestOne (0.00s)
PASS
ok  	github.com/berrak/104/cmd/104	0.001s
   create-stamp debian/debhelper-build-stamp
   dh_testroot -O--buildsystem=golang
   dh_prep -O--buildsystem=golang
   dh_auto_install -O--buildsystem=golang
	mkdir -p /home/bekr/PKGDEB/104/104/debian/104/usr/share/gocode/src/github.com/berrak/104
	cp -r -T src/github.com/berrak/104 /home/bekr/PKGDEB/104/104/debian/104/usr/share/gocode/src/github.com/berrak/104
   dh_installdocs -O--buildsystem=golang
   dh_installchangelogs -O--buildsystem=golang
   dh_perl -O--buildsystem=golang
   dh_link -O--buildsystem=golang
   dh_strip_nondeterminism -O--buildsystem=golang
   dh_compress -O--buildsystem=golang
   dh_fixperms -O--buildsystem=golang
   dh_strip -O--buildsystem=golang
   dh_makeshlibs -O--buildsystem=golang
   dh_shlibdeps -O--buildsystem=golang
   dh_installdeb -O--buildsystem=golang
   dh_golang -O--buildsystem=golang
   dh_gencontrol -O--buildsystem=golang
dpkg-gencontrol: warning: Depends field of package 104: unknown substitution variable ${shlibs:Depends}
   dh_md5sums -O--buildsystem=golang
   dh_builddeb -u-Zxz -O--buildsystem=golang
dpkg-deb: building package '104' in '../104_0.0~git20180304.eaf0f5b-1_amd64.deb'.
```
We have our debian package built, albeit with a warning about ${shlibs:Depends} from the control file.
We will ignore that for now. Next step is to debanize in a chrooted environment. Set it up with:
```bash
DIST=stretch ARCH=amd64 git-pbuilder create
sudo ln -s /var/cache/pbuilder/base-stretch-amd64.cow base.cow
```
Before removing the build directory **104**, note that no binary is packaged.
```bash
cd debian
tree 104
104
├── DEBIAN
│   ├── control
│   └── md5sums
└── usr
    └── share
        ├── doc
        │   └── 104
        │       ├── changelog.Debian.gz
        │       └── copyright
        └── gocode
            └── src
                └── github.com
                    └── berrak
                        └── 104
                            └── cmd
                                └── 104
                                    ├── 104.go
                                    ├── one.go
                                    └── one_test.go

12 directories, 7 files
```

We will have to add a directive to include our binary file in the debin directory to take care of that.
Remove the the directory **104** and all new files in **debian** that above runs have created.
Also remove the directory ../*obj-x86_64-linux-gnu*.

```bash
rm -fr 104
rm -fr ../obj-x86_64-linux-gnu
rm 104.substvars
rm debhelper-build-stamp
rm files
```

Before we can build in the chroot we have to update the **rules**.
The binaries should be end up in /opt/ZUL/bin after our enterprise name **ZUL**:

```bash
#!/usr/bin/make -f

TMP  = $(CURDIR)/debian/tmp

GO_SRC := 104.go one.go

%:
        dh $@ --buildsystem=golang --with=golang
        
override_dh_auto_build:
        go build $(GO_SRC)

override_dh_auto_install:
        mkdir -p $(TMP)/opt/ZUL/bin
        cp 104 $(TMP)/opt/ZUL/bin        
```
Since we want to have the binary to end up in a new directory,
we have to add a new file *104.install* in the debian directory:
```bash
/opt/ZUL/bin
```
Add all debian files to git:
```bash
cd debian
git add .
git com -m 'Initial packaging'
```
Now we can debanize the 104 application in the chrooted stretch (note: golang-1.7-go (1.7.4-2)) environment:
```bash
gbp buildpackage --git-pbuilder --git-compression=xz
...
...
The following packages have unmet dependencies:
 pbuilder-satisfydepends-dummy : Depends: golang-github-berrak-pkg103greet-dev which is a virtual package and is not provided by any available package

Unable to resolve dependencies!  Giving up...
...
...
Abort.
E: pbuilder-satisfydepends failed.
I: Copying back the cached apt archive contents
I: unmounting dev/ptmx filesystem
I: unmounting dev/pts filesystem
I: unmounting dev/shm filesystem
I: unmounting proc filesystem
I: unmounting sys filesystem
I: Cleaning COW directory
I: forking: rm -rf /var/cache/pbuilder/build/cow.32744
gbp:error: 'git-pbuilder' failed: it exited with 1
```
Now this fails, since all repositories known to APT is given in */etc/apt/sources.list* or in *sources.list.d*.
Before we can continue we need to upload our previous built *golang-github-berrak-pkg103greet-dev* to our own apt repository.

Before we can continue, that has to be taken care of.
See *setting-up-an-apt-repository.md* in the docs folder. Then come back here and continue.

A number of new files have been created at parent directory:
```bash
104_0.0~git20180303.53e3bf5-1_amd64.build
104_0.0~git20180303.53e3bf5-1_amd64.buildinfo
104_0.0~git20180303.53e3bf5-1_amd64.changes
104_0.0~git20180303.53e3bf5-1_amd64.deb
104_0.0~git20180303.53e3bf5-1.debian.tar.xz
104_0.0~git20180303.53e3bf5-1.dsc
104_0.0~git20180303.53e3bf5.orig.tar.xz
itp-104.txt
```
Install the 104 application:
```bash
sudo dpkg -i 104_0.0~git20180303.53e3bf5-1_amd64.deb
Selecting previously unselected package 104.
(Reading database ... 158213 files and directories currently installed.)
Preparing to unpack 104_0.0~git20180303.53e3bf5-1_amd64.deb ...
Unpacking 104 (0.0~git20180303.53e3bf5-1) ...
Setting up 104 (0.0~git20180303.53e3bf5-1) ...
```
Run application:
```bash
/opt/ZUL/bin/104
Hello 123
```
To see what is in the package:
```bash
dpkg -L 104
/.
/opt
/opt/ZUL
/opt/ZUL/bin
/opt/ZUL/bin/104
/usr
/usr/share
/usr/share/doc
/usr/share/doc/104
/usr/share/doc/104/changelog.Debian.gz
/usr/share/doc/104/copyright
```

In another chapter in the ZUL enterprise golang journey, usage of all this will be explained.
The short git aliases used here is in user ~/.gitconfig:
```bash
[aliases]
com = commit
sw2 = checkout
vbr = branch -a
```


