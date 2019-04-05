# IPMI board

&copy; Andrew Ying 2019.

A web-accessable IPMI / IP KVM system that provides full keyboard control, monitor view, and and the ability to reboot computers with standard motherboards remotely as if you were sitting in front of them with a keyboard and monitor.

![Screenshot](https://github.com/Fmstrat/diy-ipmi/raw/master/Resources/Screenshot.png)

## Requirements
* Rasberry Pi 3 Model B+ (https://www.amazon.co.uk/dp/B07BDR5PDW)
* Micro SD card (https://www.amazon.co.uk/dp/B073K14CVB)
* 2.5Amp power (https://www.amazon.co.uk/dp/B01CCR5P8U)

### Power Switch
* GPIO cables (https://www.amazon.co.uk/dp/B01LVVIOUO)
* Relay board (https://www.amazon.co.uk/dp/B0057OC5WK)
* 3 core electrical cable (https://www.amazon.co.uk/dp/B01N5WVQP2)

### On Air Sign
* On Air LED sign (https://www.aliexpress.com/item/LB480-On-Air-Recording-Studio-NEW-NR-LED-Neon-Light-Sign-home-decor-crafts/1000006552370.html)
* Mains plug (https://uk.rs-online.com/web/p/mains-plus-sockets/0490009)

### IP-KVM
* Rasberry Pi Zero WH (https://www.modmypi.com/raspberry-pi/raspberry-pi-zero-board-379/rpi-zero-board/raspberry-pi-zero-wireless-pre-soldered-header)
* Micro SD card (https://www.amazon.co.uk/dp/B073K14CVB)
* Easycap UTV007 device (https://www.amazon.co.uk/dp/B07B9ZXLN2)
* HDMI to S-Video (https://www.amazon.co.uk/dp/B0088HMN3E)
* USB TTL Serial cable (https://www.amazon.co.uk/dp/B00KAE2EL4)
* Micro USB to USB male cable (https://www.amazon.co.uk/dp/B074VM7SMM)

## Before assembling

Flash https://downloads.raspberrypi.org/raspbian_lite_latest onto the micro SD cards.

Before putting the SD into the Pi0, add this to the end of `/boot/config.txt`:
```
[...]
dtoverlay=dwc2
```
Then add `libcomposite` to the end of `/etc/modules`.

Before putting the SD into the Pi3, create a blank file called `SSH` on the boot drive to enable SSH.

## Setting up the hardware
Here is a diagram of how you connect all of the pieces:

![Diagram](https://repository.camfm.co.uk/tech/ipmi/raw/master/Resources/Diagram.png)

The detailed process is to:
- Connect the Pi3 to the relay board using this method: http://youtu.be/oaf_zQcrg7g
- Connect the Pi0 to the Pi3 using this method: https://www.thepolyglotdeveloper.com/2017/02/connect-raspberry-pi-pi-zero-usb-ttl-serial-cable/. You do not need to supply power to the Pi0, it will get power via the GPIO pins.
- Plug the easycap device and the USB TTL device into the USB ports on the Pi3
- Connect the HDMI out of your computer into the HDMI to S-Video box, and connect it to the EasyCap device via an S-Video cable
- Connect the Pi0 to the server via a micro USB to USB male cable

## Installation

### Setting up the Pi 3

First, let's get all the software we need:
```
sudo apt-get update
sudo apt-get -y install libav-tools screen lighttpd php php-cgi git
cd /opt
sudo git clone git@repository.camfm.co.uk:tech/ipmi.git
sudo chown pi ipmi -R
chmod +x /opt/ipmi/Pi3/*.py
chmod +x /opt/ipmi/Pi3/*.sh
```

To test the ability to reboot the computer via the relay board, run:
```
/opt/ipmi/Pi3/rebootServer.py
```
Test this script to see if it resets the computer. If you are interested, look in the python script to see the numbers associated with which of the 8 relays you could use for multiple computers.

Next, set up the HTTP server.
```

read -s -p "Password for web IPMI console (user 'ipmi'): " IPMIPASS
echo ""
echo "ipmi:${IPMIPASS}" | sudo tee --append /var/www/ipmipasswd

sudo lighty-enable-mod fastcgi-php
echo '' | sudo tee --append /etc/lighttpd/lighttpd.conf
sudo adduser www-data gpio
echo 'server.modules += ( "mod_auth" )' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo 'auth.debug = 2' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo 'auth.backend = "plain"' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo 'auth.backend.plain.userfile = "/var/www/ipmipasswd"' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo 'auth.require = ( "/" =>' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo '        (' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo '                "method" => "basic",' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo '                "realm" => "Password protected area",' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo '                "require" => "user=ipmi"' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo '        )' | sudo tee --append /etc/lighttpd/lighttpd.conf
echo ')' | sudo tee --append /etc/lighttpd/lighttpd.conf

cd /var/www/
sudo mv /var/www/html /var/www/html.orig
sudo ln -s /opt/ipmi/Pi3/html /var/www/html

echo '[Server 1]' | sudo tee --append /etc/ipmi.conf
echo 'TTY=/dev/ttyUSB0' | sudo tee --append /etc/ipmi.conf
echo 'VID=/dev/video0' | sudo tee --append /etc/ipmi.conf
echo 'INP=1' | sudo tee --append /etc/ipmi.conf
echo 'PIN=2' | sudo tee --append /etc/ipmi.conf

sudo service lighttpd force-reload
sudo systemctl restart lighttpd
sudo systemctl enable lighttpd
```

Lastly, get everything running right by completing some tasks and updating `/etc/rc.local` to ensure they happen at boot:
```
sudo chmod a+rw /dev/video0
sudo mkdir -p /mnt/ramdisk
sudo mount -t tmpfs -o size=3m tmps /mnt/ramdisk
sudo chown www-data /mnt/ramdisk
sudo v4l2-ctl -d /dev/video0 --set-input=1
sudo chmod a+rw /dev/ttyUSB0

sudo sed -i 's/exit 0//g' /etc/rc.local
echo "chmod a+rw /dev/video0" | sudo tee --append /etc/rc.local
echo "mkdir -p /mnt/ramdisk" | sudo tee --append /etc/rc.local
echo "mount -t tmpfs -o size=3m tmps /mnt/ramdisk" | sudo tee --append /etc/rc.local
echo "chown www-data /mnt/ramdisk" | sudo tee --append /etc/rc.local
echo "v4l2-ctl -d /dev/video0 --set-input=1" | sudo tee --append /etc/rc.local
echo "chmod a+rw /dev/ttyUSB0" | sudo tee --append /etc/rc.local
echo "exit 0" | sudo tee --append /etc/rc.local
```

### Setting up the Pi 0

Make sure you can access the Pi0 from the Pi3 by running:
```
screen /dev/ttyUSB0 115200
```
Press enter until you see a login prompt. Do not login. Instead, exit the session by hitting `Control-A` then typing `:quit` and pressing enter.

On the Pi3, run:
```
/opt/ipmi/Pi3/checkPi0Login.sh

echo "sudo systemctl enable serial-getty@ttyAMA0.service" >> /dev/ttyUSB0
echo "sudo cp /lib/systemd/system/serial-getty@.service /etc/systemd/system/serial-getty@ttyAMA0.service" >> /dev/ttyUSB0
echo "sudo sed -i 's/agetty --keep-baud 115200/agetty -a pi --keep-baud 115200/g' /etc/systemd/system/serial-getty@ttyAMA0.service" >> /dev/ttyUSB0
echo "sudo systemctl daemon-reload" >> /dev/ttyUSB0

echo "sudo systemctl disable networking" >> /dev/ttyUSB0
echo "sudo apt-get -y remove dhcpcd5 isc-dhcp-client isc-dhcp-common" >> /dev/ttyUSB0
sleep 60

echo "rm -f /tmp/B64" >> /dev/ttyUSB0
for LINE in $(base64 /opt/diy-ipmi/Pi0/enableHID.sh); do echo "echo $LINE >> /tmp/B64" >> /dev/ttyUSB0; done
echo "base64 -d /tmp/B64 > /home/pi/enableHID.sh" >> /dev/ttyUSB0
echo "chmod +x /home/pi/enableHID.sh" >> /dev/ttyUSB0

echo "rm -f /tmp/B64" >> /dev/ttyUSB0
for LINE in $(base64 /opt/diy-ipmi/Pi0/sendkeys.c); do echo "echo $LINE >> /tmp/B64" >> /dev/ttyUSB0; done
echo "base64 -d /tmp/B64 > /home/pi/sendkeys.c" >> /dev/ttyUSB0
echo "gcc -o /home/pi/sendkeys /home/pi/sendkeys.c" >> /dev/ttyUSB0

sudo apt-get -y install libusb-dev
cd /opt/diy-ipmi/Pi0/
gcc -o hub-ctrl hub-ctrl.c -lusb
for LINE in $(base64 hub-ctrl); do echo "echo $LINE >> /tmp/B64" >> /dev/ttyUSB0; done
echo "base64 -d /tmp/B64 > /home/pi/hub-ctrl" >> /dev/ttyUSB0
echo "chmod +x /home/pi/hub-ctrl" >> /dev/ttyUSB0
cd -

echo "sudo /home/pi/enableHID.sh" >> /dev/ttyUSB0
echo "sudo sed -i 's/exit 0//g' /etc/rc.local" >> /dev/ttyUSB0
echo "echo /home/pi/enableHID.sh | sudo tee --append /etc/rc.local" >> /dev/ttyUSB0
echo "echo exit 0 | sudo tee --append /etc/rc.local" >> /dev/ttyUSB0
```

### Access the IPMI
You should now be able to access the IPMI console at `http://<RaspberryPi3IP>/`. From here you can set up SSL and port forwarding to the device as your situation requires.

## Managing multiple servers
Multiple servers can be managed by using multile USB capture devices (one per server), multiple USB-to-Serial adapters each with their own Pi0 (one per server), and using one of the 8 relays on the relay board as the reset button. Once installed, edit `/etc/ipmi.conf`. The format should be:
```
[My first server]
TTY=/dev/ttyUSB0
VID=/dev/video0
INP=1
PIN=2

[My second server]
TTY=/dev/ttyUSB1
VID=/dev/video1
INP=1
PIN=3
```
- When hooking up multiple TV tuners and USB-to-Serial devices to Pi0s, it is recommended to use a powered USB hub.
- The `INP` field is `1` for s-video or `0` for composite. You can actually use one capture device for two servers by using both inputs.
- The `PIN` value is the pin on the relay board. The valid values are `2`, `3`, `4`, `17`, `27`, `22`, `10`, `9`.


## Tips
If you're using this on a Linux headless machine that's not running X11, it's easier to read the console if things are nice and big. Disabling the video driver is an easy way to accomplish this as it will just default to a generic NTSC or PAL signal. For instance, with an nVidia card, you could follow these directoins: https://askubuntu.com/questions/481414/install-nvidia-driver-instead-of-nouveau

To enable SSL, you need a combined PEM file which can be created by:
```
cat myserver.key mycert.crt myca.ca-bundle > mypem.pem
```
Then add this into the `/etc/lighttpd/lighttpd.conf` file (editing as required):
```
$SERVER["socket"] == ":443" {
	ssl.engine = "enable"
	ssl.pemfile = "/opt/certs/mypem.pem"
	server.name = "host.domain.com"
	server.document-root = "/var/www/html"
	ssl.use-sslv2 = "disable"
	ssl.use-sslv3 = "disable"
	ssl.use-compression = "disable"
	ssl.honor-cipher-order = "enable"
	ssl.cipher-list = "ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-AES256-SHA384:ECDHE-RSA-AES256-SHA384:ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES128-SHA256:ECDHE-RSA-AES128-SHA256:ECDHE-ECDSA-RC4-SHA:ECDHE-RSA-RC4-SHA:ECDH-ECDSA-RC4-SHA:ECDH-RSA-RC4-SHA:ECDHE-RSA-AES256-SHA:RC4-SHA"
	server.errorlog = "/var/log/lighttpd/serror.log"
	accesslog.filename = "/var/log/lighttpd/saccess.log"
}
```

## Troubleshooting

If you're not getting video, here are some troubleshooting methods:

First make sure the `INP` field in `/etc/ipmi.conf` is set to the right input. Input 0 is usually Composite, and Input 1 is usually S-Video.

Connect a source and test to see if it's working.
```
apt-get install mencoder
mencoder tv:// -tv driver=v4l2:norm=NTSC:device=/dev/video0:input=1:fps=5 -nosound -ovc copy -o test.avi
v4l2-ctl -d /dev/video0 --set-input=1
avconv -f video4linux2 -i /dev/video0 -vframes 1 -s 720x480 test.jpg
```
Control-C that, and sftp the files to a host for viewing and playback:
```
sftp test.avi root@hostname:/folder/test.avi
sftp test.avi root@hostname:/folder/test.jpg
```

You can test keyboard control from the Pi0 with commands like:
```
/home/pi/sendkeys /dev/hidg0 keyboard
```
You will need to type things like "a" and press ENTER to send the "a." Other ways to send keys include things like:
```
echo 'a' | /home/pi/sendkeys /dev/hidg0 keyboard
echo 'return' | /home/pi/sendkeys /dev/hidg0 keyboard
echo 'shift a' | /home/pi/sendkeys /dev/hidg0 keyboard
echo 'left-meta space' | /home/pi/sendkeys /dev/hidg0 keyboard
```

To reset the keyboard controller run this on the Pi3:
```
/opt/diy-ipmi/Pi3/resetPi0HID.sh
```

If you get wierd USB errors, make sure your power brick is strong enough and that you're not trying to power the HDMI to S-Video converter via the Pi3.
