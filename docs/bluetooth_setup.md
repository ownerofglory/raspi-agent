# Bluetooth setup
It assumes you’re using PipeWire as your audio system (recommended over PulseAudio for Bluetooth).

## Install Required Components

Update your package list and install Bluetooth and PipeWire components:
```shell
sudo apt update
sudo apt install -y \
    bluez bluez-tools \
    pipewire pipewire-pulse wireplumber \
    pipewire-audio-client-libraries \
    libspa-0.2-bluetooth

```

## Restart Services

After installation, restart the Bluetooth and audio services:

```shell
sudo systemctl restart bluetooth
systemctl --user restart pipewire pipewire-pulse wireplumber

```

## Verify Audio Stack

Confirm that the WirePlumber service (PipeWire session manager) is running correctly:
```shell
systemctl --user status wireplumber

```

## Connect Your Bluetooth Speaker

Start the Bluetooth control tool:
```shell
bluetoothctl

```
Inside the interactive shell, run:
```shell
power on
agent on
scan on

```
Wait until your speaker appears, note its MAC address (XX:XX:XX:XX:XX:XX), then run:
```shell
pair XX:XX:XX:XX:XX:XX
trust XX:XX:XX:XX:XX:XX
connect XX:XX:XX:XX:XX:XX

```
Once connected, you should see a confirmation like:
```shell
Connection successful
```

Verify BlueZ Integration

Check that your Bluetooth speaker was registered as a PipeWire audio card:
```shell
pactl list cards short

```
You should see a line similar to:
```shell
bluez_card.XX_XX_XX_XX_XX_XX    PipeWire    ...
```

Test Audio Output

Play a simple stereo test sound:
```shell
speaker-test -t wav -c 2

```