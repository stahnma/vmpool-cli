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
  homepage "https://github.com/stahnma/vmpool-cli"
  url "http://yum.stahnkage.com/sources/vmpool-0.0.0.13.g648ce9d.tar.gz"
  sha1 "1a51a0ebf4b6b7aa510dcbec70dc6c3216ea7e34"

  head "https://github.com/stahnma/vmpool-cli.git", :shallow => false, :using => VmpoolHeadDownloadStrategy

  depends_on "go" => :build

  def install
    system "make vmpool"
    bin.install 'vmpool'
  end

  def caveats
    "You're going to want to build head on this, since I'll never update the source tarball"
  end
end
