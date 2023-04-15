#!/bin/bash
# See http://redsymbol.net/articles/unofficial-bash-strict-mode/
set -euo pipefail
IFS=$'\n\t'

# Settings
VENDOR_TMP_DIR=.vendor-tmp
OUTPUT_DIR=src/services/gui/html/static/vendor

# Create temporary dir
rm -rf ${VENDOR_TMP_DIR:?} || true
mkdir ${VENDOR_TMP_DIR:?}

# ===========================================================================

# Get all packages
npm install --prefix=src

# Bootstrap
mkdir ${VENDOR_TMP_DIR:?}/bootstrap
cp -rt ${VENDOR_TMP_DIR:?}/bootstrap src/node_modules/bootstrap/dist

# Bootstrap Icons
mkdir ${VENDOR_TMP_DIR:?}/bootstrap-icons
cp -rt ${VENDOR_TMP_DIR:?}/bootstrap-icons src/node_modules/bootstrap-icons/font

# Luxon
mkdir ${VENDOR_TMP_DIR:?}/luxon
cp -rt ${VENDOR_TMP_DIR:?}/luxon src/node_modules/luxon/build/global/luxon.min.js

# Simple Keyboard
mkdir ${VENDOR_TMP_DIR:?}/simple-keyboard
cp -rt ${VENDOR_TMP_DIR:?}/simple-keyboard src/node_modules/simple-keyboard/build

# Popper JS
mkdir ${VENDOR_TMP_DIR:?}/popperjs-core
cp -rt ${VENDOR_TMP_DIR:?}/popperjs-core src/node_modules/@popperjs/core/dist

# Tempus Dominus (time picker)
mkdir ${VENDOR_TMP_DIR:?}/tempus-dominus
cp -rt ${VENDOR_TMP_DIR:?}/tempus-dominus src/node_modules/@eonasdan/tempus-dominus/dist

# ===========================================================================

# Move output into place and cleanup
rm -rf ${OUTPUT_DIR:?}
mkdir -p ${OUTPUT_DIR:?} # Ensure parent folder exists
rmdir ${OUTPUT_DIR:?}
mv ${VENDOR_TMP_DIR:?} ${OUTPUT_DIR:?}
touch ${OUTPUT_DIR:?}/.gitkeep
rm -rf ${VENDOR_TMP_DIR:?}
