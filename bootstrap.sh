export PATH=$coreutils/bin:$PATH

mkdir -p "$out"

rawsrc="$1"
gobin="$2"

src=$(cat "$rawsrc")
src=${src//defcustom pus-path \"pus\"/defcustom pus-path \"${gobin}\"}

echo "$src" >"$out/pus.el"
