#!/bin/bash
for s in `find assets/ icons/ -iname *.png`; do
    echo $s
    B64D=`cat $s | base64 -w0`
    sed -i "s/url(\"${s/\//\\\/}\")/url(\"data:image\/png;base64,${B64D//\//\\\/}\")/g" gtk.css
done

for s in `find assets/ icons/ -iname *.svg`; do
    echo $s
    B64D=`cat $s | base64 -w0`
    sed -i "s/url(\"${s/\//\\\/}\")/url(\"data:image\/svg+xml;base64,${B64D//\//\\\/}\")/g" gtk.css
done
