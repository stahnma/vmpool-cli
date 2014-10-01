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
  # The version inside the tarball is nil. This suck. #FIXME
  url "http://yum.stahnkage.com/sources/vmpool-0.0.0.10.gcf9ed2a.tar.gz"
  sha1 "1b6b4dd5e08bad5fe61b5d7d705fa3f56556db50"

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
