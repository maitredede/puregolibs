# Purego

Call native libraries without CGO, using [https://github.com/ebitengine/purego](github.com/ebitengine/purego).

## CEC

https://github.com/Pulse-Eight/libcec

## Fontconfig

https://gitlab.freedesktop.org/fontconfig/fontconfig

## GoUSB

https://github.com/libusb/libusb/

## Imagick

https://imagemagick.org/

## libfreefare

https://github.com/nfc-tools/libfreefare

## libnfc

https://github.com/nfc-tools/libnfc

## Plutobook

https://github.com/plutoprint/plutobook

## sane

- Based on https://github.com/tjgq/sane/
- Added multi-image scan: https://github.com/tjgq/sane/pull/13

## Packages

### Debian

```sh
apt-get install \
  libcec-dev libp8-platform-dev libudev-dev \
  libfontconfig-dev \
  libusb-1.0-0-dev \
  libmagickwand-dev \
  libfreefare-dev \
  libnfc-dev \
  libsane-dev
```

For plutobook : https://github.com/plutoprint/plutobook?tab=readme-ov-file#installation-guide
