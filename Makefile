PKGNAME=vmpool

UNAME:=$(shell uname)
ifeq ($(UNAME),Darwin)
TMP_PATTERN:=$(shell mktemp -d tmpbuild-XXXXXX)
TMPDIR:=$(shell pwd)/$(TMP_PATTERN)
TAR_TMP_DIR:=$(shell mktemp -d -t tmptarball)
else
TMP_PATTERN:=$(shell mktemp -d -u -p . -t rpmbuild-XXXXXXX)
TMPDIR:=$(shell pwd)/$(TMP_PATTERN)
TAR_TMP_DIR:=$(shell mktemp -d -u -t tarball-XXXXXXX)
endif

SPEC_FILE=$(PKGNAME).spec
PWD:=$(shell pwd)

RPMBUILD := $(shell if test -f /usr/bin/rpmbuild ; then echo /usr/bin/rpmbuild ; else echo "x" ; fi)
RPM_DEFINES = --define "_specdir $(TMPDIR)/SPECS" --define "_rpmdir $(TMPDIR)/RPMS" --define "_sourcedir $(TMPDIR)/SOURCES" --define "_srcrpmdir $(TMPDIR)/SRPMS" --define "_builddir $(TMPDIR)/BUILD"

MAKE_DIRS= $(TMPDIR)/SPECS $(TMPDIR)/SOURCES $(TMPDIR)/BUILD $(TMPDIR)/SRPMS $(TMPDIR)/RPMS

# If there is a VERSION file, use it.
#   else use git describe
VERSION:=$(shell if [ -f "VERSION" ] ; then  cat VERSION ; else   git describe | sed -e 's/-/\./g' ; fi)

TARBALL=$(PKGNAME)-$(VERSION).tar.gz

build: vmpool

fmt:
	go fmt *.go ; rm -rf tmpbuild-*

vmpool:
	go build -o vmpool -ldflags "-X=main.version=$(VERSION)" *.go
	@rm -rf tmp*

install:
	mkdir -p $(DESTDIR)/usr/local/bin
	cp -pr vmpool $(DESTDIR)/usr/local/bin

linux:
# In order to get the cross-compile options on mac, install via
#     brew install go --cross-compile-common
	GOARCH=amd64 GOOS=linux go build

clean:
	rm -rf vmpool *tar.gz rpmbuild-* *.src.rpm tmp* VERSION *.tar

uninstall:
	rm -rf $(DESTDIR)/usr/local/bin/vmpool

tarball:
	git archive --format=tar --prefix=vmpool-$(VERSION)/ $(shell git describe)  > vmpool-$(VERSION).tar
	mkdir vmpool-$(VERSION)
	[ -f "VERSION" ] || echo $(VERSION) > VERSION
	mv VERSION vmpool-$(VERSION)
	tar rf vmpool-$(VERSION).tar vmpool-$(VERSION)/VERSION
	gzip vmpool-$(VERSION).tar
	rm -rf vmpool-$(VERSION) tmpbuild-*

# If you're on a system with rpm, you can build a srpm to throw at mock or something.
srpm: tarball
	@mkdir -p $(MAKE_DIRS)
	@cp -f $(TARBALL) $(TMPDIR)/SOURCES
	@cp -f $(SPEC_FILE) $(TMPDIR)/SPECS
	sed -i 's/==VERSION==/$(VERSION)/g' $(TMPDIR)/SPECS/$(SPEC_FILE)
	@wait
	@$(RPMBUILD) $(RPM_DEFINES) -bs $(TMPDIR)/SPECS/$(SPEC_FILE)
	@mv -f $(TMPDIR)/SRPMS/* .
	@rm -rf $(TMPDIR)
	@echo
	@ls *src.rpm

.PHONY: intall fmt clean tarball uninstall linux all build vmpool
