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
cd ${VENDOR_TMP_DIR:?}
mkdir output

# ===========================================================================

# ===================================
# =            Bootstrap            =
# ===================================
# Bootstrap
wget -O package.tgz https://registry.npmjs.org/bootstrap/-/bootstrap-5.1.3.tgz
tar -xzf package.tgz
mkdir output/bootstrap
mv -t output/bootstrap package/dist package/LICENSE package/README.md
rm -rf package*

# Bootstrap Icons
wget -O package.tgz https://registry.npmjs.org/bootstrap-icons/-/bootstrap-icons-1.9.1.tgz
tar -xzf package.tgz
mkdir output/bootstrap-icons
mv -t output/bootstrap-icons package/font package/LICENSE.md package/README.md
rm -rf package*

# Bootstrap Dark 5
wget -O package.tgz https://registry.npmjs.org/bootstrap-dark-5/-/bootstrap-dark-5-1.1.3.tgz
tar -xzf package.tgz
mkdir output/bootstrap-dark-5
mv -t output/bootstrap-dark-5 package/dist package/LICENSE.md package/README.md
rm -rf package*

# ===================================
# =              Luxon              =
# = https://moment.github.io/luxon/ =
# ===================================
wget -O package.tgz https://registry.npmjs.org/luxon/-/luxon-3.1.0.tgz
tar -xzf package.tgz
mkdir output/luxon
mv -t output/luxon package/build/global/luxon.min.js package/README.md
rm -rf package*

# ===================================
# =         Simple Keyboard         =
# ===================================
wget -O package.tgz https://github.com/hodgef/simple-keyboard/archive/refs/tags/3.4.0.tar.gz
tar -xzf package.tgz
mkdir output/simple-keyboard
mv -t output/simple-keyboard simple-keyboard-*/*
rm -rf package.tgz simple-keyboard*

# ===================================
# =           Time picker           =
# =    https://getdatepicker.com    =
# ===================================
# Popper JS
wget -O package.tgz https://registry.npmjs.org/@popperjs/core/-/core-2.11.6.tgz
tar -xzf package.tgz
mkdir output/popperjs-core
mv -t output/popperjs-core package/dist package/LICENSE.md package/README.md
rm -rf package*

# Tempus Dominus (time picker)
wget -O package.tgz https://github.com/Eonasdan/tempus-dominus/archive/refs/tags/v6.2.6.tar.gz
tar -xzf package.tgz
mkdir output/tempus-dominus
mv -t output/tempus-dominus tempus-dominus-*/dist tempus-dominus-*/LICENSE tempus-dominus-*/README.md
rm -rf package.tgz tempus-dominus*

# ===========================================================================

# Move output into place and cleanup
cd -
rm -rf ${OUTPUT_DIR:?}
mv ${VENDOR_TMP_DIR:?}/output ${OUTPUT_DIR:?}
touch ${OUTPUT_DIR:?}/.gitkeep
rm -rf ${VENDOR_TMP_DIR:?}
