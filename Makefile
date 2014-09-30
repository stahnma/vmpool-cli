
PKGNAME=vmpool-cli
TMP_PATTERN:=$(shell mktemp -d -u -p . -t rpmbuild-XXXXXXX)
TMPDIR=$(shell pwd)/$(TMP_PATTERN)
TAR_TMP_DIR:=$(shell mktemp -d -u -t tarball-XXXXXXX)

SPEC_FILE=$(PKGNAME).spec


RPMBUILD := $(shell if test -f /usr/bin/rpmbuild ; then echo /usr/bin/rpmbuild ; else echo "x" ; fi)
RPM_DEFINES = --define "_specdir $(TMPDIR)/SPECS" --define "_rpmdir $(TMPDIR)/RPMS" --define "_sourcedir $(TMPDIR)/SOURCES" --define "_srcrpmdir $(TMPDIR)/SRPMS" --define "_builddir $(TMPDIR)/BUILD"

MAKE_DIRS= $(TMPDIR)/SPECS $(TMPDIR)/SOURCES $(TMPDIR)/BUILD $(TMPDIR)/SRPMS $(TMPDIR)/RPMS

ifndef VERSION
VERSION:=$(shell git describe | sed -e 's/-/\./g')
endif

TARBALL=$(PKGNAME)-$(VERSION).tar.gz

build: vmpool

fmt:
	go fmt vmpool.go

vmpool:
	go build -ldflags "-X main.version $(VERSION)" vmpool.go


install:
	mkdir -p $(DESTDIR)/usr/local/bin
	cp -pr vmpool $(DESTDIR)/usr/local/bin

linux:
# In order to get the cross-compile options on mac, install via
#     brew install go --cross-compile-common
	GOARCH=amd64 GOOS=linux go build

clean:
	rm -rf vmpool *tar.gz rpmbuild-* *.src.rpm

uninstall:
	rm -rf $(DESTDIR)/usr/local/bin/vmpool

tarball:
	mkdir -p $(TAR_TMP_DIR)/$(PKGNAME)-$(VERSION)
	cd ..; cp -pr $(PKGNAME)/* $(TAR_TMP_DIR)/$(PKGNAME)-$(VERSION); rm -rf $(TAR_TMP_DIR)/$(PKGNAME)-$(VERSION)/{contrib,*.spec}
	cd $(TAR_TMP_DIR); tar pczf $(TARBALL)  $(PKGNAME)-$(VERSION)
	mv $(TAR_TMP_DIR)/$(TARBALL) .
	rm -rf $(TAR_TMP_DIR)

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


.PHONY: intall fmt clean
