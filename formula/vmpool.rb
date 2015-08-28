require "formula"

# We use a custom download strategy to properly configure
# vmpool's version information when built against HEAD.
# This is populated from git information using git describe.
class VmpoolHeadDownloadStrategy < GitDownloadStrategy
  def stage
    @clone.cd {reset}
    safe_system 'git', 'clone', @clone, '.'
  end
end

class Vmpool < Formula
  desc "Interact with a vmpooler instance easily."

  homepage "https://github.com/stahnma/vmpool-cli"
  url "http://yum.stahnkage.com/sources/vmpool-0.2.0.tar.gz"

  head "https://github.com/stahnma/vmpool-cli.git", :shallow => false, :using => VmpoolHeadDownloadStrategy

  depends_on "go" => :build
  depends_on "pandoc" => :build

  def install
    system "make vmpool"
    bin.install 'vmpool'
    system "pandoc -s -t man manpage.md -o vmpool.1"
    man1.install "vmpool.1"
  end

  def caveats
    "You're going to want to build head on this, since I rarely update the source tarball."
  end
end
