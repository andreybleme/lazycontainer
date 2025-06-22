class Lazycontainer < Formula
  desc "CLI for listing containers and viewing logs"
  homepage "https://github.com/andreybleme/lazycontainer"
  version "0.0.1"

  if Hardware::CPU.intel?
    url "https://github.com/andreybleme/lazycontainer/releases/download/v0.0.1/lazycontainer_darwin_amd64"
    sha256 "a5457ddf47d9f1e714c85afdffcd2fe037f5ca130cc8d7761ff8733ab24146f9"
  else
    url "https://github.com/andreybleme/lazycontainer/releases/download/v0.0.1/lazycontainer_darwin_arm64"
    sha256 "1358878c2be97fa1cced5c0baee41804ec6e7c2c45c17d563c37d79c81e1dcec"
  end

  def install
    binary = Hardware::CPU.intel? ? "lazycontainer_darwin_amd64" : "lazycontainer_darwin_arm64"
    bin.install binary => "lazycontainer"
  end

  test do
    # simple smoke‑test; exit code 0 and print help/version
    assert_match "lazycontainer", shell_output("#{bin}/lazycontainer --help")
  end
end
