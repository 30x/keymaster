#!/bin/sh
rm testbundle.zip
rm testsystem.zip
cd template/testbundle
zip ../../testbundle.zip -r .
cd ../../template/testsystem
zip ../../testsystem.zip -r .
cd ../..
