class Lazycontainer < Formula
  desc "CLI for listing containers and viewing logs"
  homepage "https://github.com/andreybleme/lazycontainer"
  version "1.0.0"

  if Hardware::CPU.intel?
    url "https://github.com/andreybleme/lazycontainer/releases/download/v1.0.0/lazycontainer_darwin_amd64"
    sha256 "3d9e5762c68f3f03f1f42cc9880ed4b6dcdfc0bc0fa03c7a1c5d1c7d2a76e38c"
  else
    url "https://github.com/andreybleme/lazycontainer/releases/download/v1.0.0/lazycontainer_darwin_arm64"
    sha256 "ec9b4186c8a3c7cc9d2cd0fc6b00226d887da69800562caa57b4593a4a8142d8"
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
