class Openwar < Formula
    homepage "http://www.openwar.io"
    url "https://github.com/andreas-jonsson/openwar/archive/master.zip"
    sha1 ""

    depends_on 'go' => '1.6'
    depends_on 'gtk+'
    depends_on 'sdl2'
    depends_on 'sdl2_mixer'

    def install
        system "go get"
        system "go build openwar.go"

        bin.install "openwar"
        share.install "data"
    end
end
