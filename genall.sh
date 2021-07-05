mkdir -p ./packages/base

for f in ~/sauerbraten-code/packages/base/*.ogz; do
    echo "trimming $f"
    gunzip --suffix=ogz --stdout $f | ./genserverogz
    gzip --keep --force --best --no-name trimmed.bin
    mv trimmed.bin.gz "./packages/base/$(basename $f)"
done

echo "packaging into trimmed.tar.gz"
tar -caf trimmed.tar.gz ./packages

rm -r ./packages trimmed.bin
