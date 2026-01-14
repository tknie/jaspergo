#!/bin/sh

TMPDIR=/tmp

JASPER_SAMPLE_TEST=${JASPER_SAMPLE_TEST:-$1}

for i in `find tests7 -type f -name '*.jrxml'`; do
    if [ "$JASPER_SAMPLE_TEST" = "" -o "$JASPER_SAMPLE_TEST" = "$(basename $i)" ]; then
    JRXML=$i
    echo "Comparing $i with output/$(basename $JRXML)"
    xmllint --format $i >$TMPDIR/1.xml
    xmllint --format output/$(basename $i) >$TMPDIR/2.xml
    diff -w $TMPDIR/1.xml $TMPDIR/2.xml
    if [ $? -ne 0 ]; then
        echo "Mismatch found for $JRXML"
       # exit 1
    fi
    fi
done


