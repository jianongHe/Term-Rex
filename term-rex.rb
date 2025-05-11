class TermRex < Formula
  desc "Terminal-based dinosaur runner game"
  homepage "https://github.com/jianongHe/Term-Rex"
  url "https://github.com/jianongHe/Term-Rex/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "YOUR_SHA256_HASH" # 需要替换为实际的哈希值
  license "MIT"

  depends_on "go" => :build

  def install
    system "go", "build", "-o", bin/"term-rex", "."
    prefix.install "assets"
  end

  test do
    system "#{bin}/term-rex", "--version"
  end
end
