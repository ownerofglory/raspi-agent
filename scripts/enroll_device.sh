#!/bin/sh

device_id=$1
country='DE'
state='BW'
location='Stuttgart'
organisation='ownerofglory'
org_unit='Devices'

if [ -z "$device_id" ]; then
  echo "Usage: $0 <device-id>"
  exit 1
fi

openssl req -new -newkey rsa:2048 -nodes -keyout device.key -out device.csr \
  -subj "/C=$country/ST=$state/L=$location/O=$organisation/OU=$org_unit/CN=$device_id" \
  -addext "subjectAltName=DNS:$device_id" \
  -addext "extendedKeyUsage=clientAuth"

