#!/bin/bash

PLAYS="
allswell
asyoulikeit
comedy_errors
cymbeline
lll
measure
merry_wives
merchant
midsummer
much_ado
pericles
taming_shrew
tempest
troilus_cressida
twelfth_night
two_gentlemen
winters_tale
1henryiv
2henryiv
henryv
1henryvi
2henryvi
3henryvi
henryviii
john
richardii
richardiii
cleopatra
coriolanus
hamlet
julius_caesar
lear
macbeth
othello
romeo_juliet
timon
titus"

URL=http://shakespeare.mit.edu

for p in $PLAYS; do
  ifile=$URL/$p/full.html
  ofile=$p.html
  wget -O $ofile $ifile
  gzip $ofile
done
