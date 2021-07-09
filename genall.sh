mkdir -p ./packages/base

for f in ~/sauerbraten-code/packages/base/*.ogz; do
    echo "trimming $f"
    gunzip --suffix=ogz --stdout $f \
    | ./genserverogz \
    | gzip --force --best --no-name > "./packages/base/$(basename $f)"
done

echo "packaging into trimmed.tar.gz"
tar -caf trimmed.tar.gz ./packages

rm -r ./packages
