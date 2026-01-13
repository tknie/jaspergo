#!/bin/sh

TMPDIR=/tmp

for i in `find tests7 -type f -name '*.jrxml'`; do
    echo "Comparing $i with output/$(basename $i)"
    xmllint --format $i >$TMPDIR/1.xml
    xmllint --format output/$(basename $i) >$TMPDIR/2.xml
    diff -w $TMPDIR/1.xml $TMPDIR/2.xml
    if [ $? -ne 0 ]; then
        echo "Mismatch found for $i"
       # exit 1
    fi
done


