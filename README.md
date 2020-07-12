# Porter

Welcome to Porter: A simple garage door controller written in Go.

This project is a ground-up rebuild of my previous [doorMan](https://github.com/ctrezevant/doorMan) project, with a few goals in mind:

- **Easy configuration:** Porter includes hardware, HTTP, and security configuration into an easy-to-read JSON file. The hardware component is highly
configurable, so it's very easy to wire in one or several garage doors and should be highly compatible with different systems. You can even work around
wiring mistakes.

- **Performance and Stability:** Eschewing JavaScript for Go (my favorite language) has realized significant gains in both of these areas.

- **Security:** Apart from the security benefits of Go itself, porter allows users to define scoped API keys with varying permissions for different clients.
By specifying certificate file locations in the configuration, the REST API can be served over HTTPS.

- **Clean and Simple Interface:** Everything is exposed behind an easy-to-use HTTP REST API, and a fully-functional client library (`porter/client`) is included for use
in building integrations (such as the [HomeKit integration](https://github.com/ctrezevant/hkporter), [web interface](https://github.com/ctrezevant/porter-web), and [Twilio-based status monitor](https://github.com/ctrezevant/reporter)).

- **Easy Portability:** The base Porter package supports the Raspberry Pi by way of the `go-rpio` package. However, I've done some work to improve the ease with
which support can be added for additional boards. The `porter/hw` package serves as a wrapper around platform-specific libraries, so you could (in theory) update those
bindings to your library of choice and be on your way with your favorite board.


# Installation

## Building Porter

Before you begin, you'll need to build porter. Fetch the dependencies using `go get`, and if necessary [cross compile](https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5) porter for the correct `GOOS`/`GOARCH` for your target platform.

For the Raspberry Pi, your build command would look like this:

`GOOS=linux GOARCH=arm go build -o porter main.go`

Next, you'll need to copy the compiled porter binary to the appropriate location on your host. I recommend using `/usr/local/bin/porter`.

## Installation

The following instructions walk you through the process of setting up a Porter instance on a Raspberry Pi. These instructions assume that you're using either Raspbian or Ubuntu Server.

This brief guide is intended for advanced users who require only minimal guidance to configure hardware and system services.

At some point in time, I'll probably include releases of Porter in this repository in .deb/.rpm format. This will make installation significantly easier.


### Hardware Setup

Hardware configuration is straightforward, although hardware fundamentals are outside of the scope of this document. Instead, an example wiring scheme that will work with a provided default configuration is presented here.  

Note: The Raspberry Pi unfortunately uses several pin numbering schemes simultaneously, which can be confusing. For physical wiring, I'll specify the physical pin number on the GPIO header. In the software configuration later on, we'll have to use the BCM pin numbers. Resources like [pinout.xyz](https://pinout.xyz/) are an excellent reference if you're trying to follow along.

Additional Note: The `gpio` command (install with `apt -y install rpi.gpio`) is extremely useful for troubleshooting. Specifically, `gpio readall` will print out the current state of all pins in a very readable format.

On the hardware side, you're going to need the following components:

- A 5v Relay module, such as [this one](https://www.sainsmart.com/products/2-channel-5v-relay-module) to control your garage door opener.
- A magnetic reed switch (such as [this](https://www.amazon.com/Magnetic-Switch-Normally-Closed-Security/dp/B0735BP1K4/)) to determine whether the door is open or closed.
- Some length of wiring to reach from your garage door lift and the reed switch to the pi (this will vary depending on your particular setup)

Your final pinout should look like this:

```
Pi Pin | Device
   2   | Relay Vcc In
   6   | Relay Gnd
  11   | Relay In1
  14   | Reed Switch
  16   | Reed Switch
```

Note the active and inactive state of your relays and switches- these will need to be reflected in the Porter configuration file if state information and commands are to be read/sent properly. For example, a magnetic reed switch in a normally-closed configuration would read 1 when inactive (meaning that no magnets are nearby). Similarly, a relay module (such as the Sainsmart) is active when its input signal is low.

Another important thing to keep in mind is that some pins on the Pi might be in a high or low state when the system starts up. Keep this in mind, as choosing the wrong pin could, for example, cause your door to open each time the system reboots. I cannot guarantee that the example pins provided above will be compatible with the hardware you are using, so be sure to perform some testing on your own.

For the lift, make sure that your relay is in a normally open state on the control lines. These lines are the same wires leading to physical switches you might have on the wall that allow you to open and close the door. Generally, you can tap directly into a wiring block on the back of your lift in order to attach the relay as a controller. Consult the manual for your lift and perform some testing before proceeding.

### Creating the Porter Service Account

For security reasons, I highly recommend running the Porter service as its own dedicated, unprivileged user account. You can create this account using the commands below:

Note: Adding the `porter` user to the GPIO group is an important step, as it provides the necessary permissions for Porter to modify the Pi's GPIO pins. At the time of writing Porter, certain elements of this were still under development. If you're still getting permissions errors about accessing `/dev/gpiomem` after adding the porter user to the `gpio group`, you may need to copy the udev rules in this repository (in `scripts/udev.rules`) to `/etc/udev/rules.d/gpio.rules` and reboot for the below to work.

```
useradd -M -N -r -s /bin/false -c "porter service account" porter
groupadd porter
adduser porter porter
```

### Adding Porter to Systemd

Like many services on a Linux system, the Porter service will be managed by Systemd. To enable this, we'll need to create a couple of small configuration files so that Systemd knows how to start and manage the service.

First, we'll need to create an environment file, which contains the command line arguments that Systemd should start the Porter service with. Porter takes exactly one command-line argument, `-config`, which specifies the location of your configuration file. Create the file `/etc/default/porter` with the following contents:

```
-config /etc/porter/config.json
```

Next, copy the contents of `porter.service` (located in the root of this repository) to `/etc/systemd/system/porter.service`.

Finally, we can tell Systemd to load our new Porter unit file, and to start it at boot time:

```
systemctl daemon-reload
systemctl enable porter
```

### Configuring the Porter Service

Now it's time to configure Porter.

First, create a directory in `/etc/` to hold Porter's configuration file. We'll need to set the appropriate permissions to make sure that the Porter service can read (but not overwrite) the file's contents:

```
mkdir /etc/porter
chown root:porter -R /etc/porter
chmod 644 -R /etc/porter
```

Here's what a basic Porter configuration file would look like for the wiring guide provided earlier in this tutorial.

This configuration provides the following:

- One door definition.

- Two API keys, one with full administrative privileges and another that can only read (but not modify) the door state.

- The API server is configured to listen on all addresses at port 8080 with SSL disabled.

```
{
	"http": {
		"listen_addr": "0.0.0.0:8080"
	},
	"doors": [{
		"name": "Door",
		"lift_ctl_pin": 17,
		"lift_ctl_inactive_state": 1,
		"lift_ctl_trip_time_ms": 100,
		"door_sensor_pin": 23,
		"door_sensor_closed_state": 1
	}],
	"keys": [{
			"name": "Admin",
			"secret": "some_secret_string",
			"allow_methods": ["*"]
		},
		{
			"name": "Read-Only",
			"secret": "readonly",
			"allow_methods": ["list"]
		}
	]
}
```

Edit the example to taste and place the completed file in `/etc/porter/config.json`.

### Managing Porter

To start Porter, run:

```
service porter start
```

To stop the server, run:

```
service porter stop
```

To view logs and error messages, view the service's status like so:

```
service porter status
```

### Testing Porter

Now that you've installed Porter, its Systemd unit, and configuration file, start the service by running `service porter start`.

If you encounter errors these will be logged. Use `service porter status` to diagnose and correct them.

You are now ready to test the API. Here are some example cURL commands.

##### Set variables for testing

Run this command first. If you've changed your Porter API token or door name from the example provided above, make sure you reflect that in the below commands.

Assuming you're running these commands on the Pi that's hosting porter, you can reach the API server locally at `127.0.0.1:8080`. However, if you are running these commands on a different system, be sure to update the IP/port information as well.


```
export PORTER_TOKEN="some_secret_string"
export PORTER_IP="127.0.0.1:8080"
export PORTER_DOOR="Door"
```

##### API Endpoints

- **GET** `/api/v1/list`: Lists the status of all doors Porter is configured to manage.

- **PUT** `/api/v1/open/<door name>`: Opens the specified door. The command will succeed only if the door is currently closed.

- **PUT** `/api/v1/close/<door name>`: Closes the specified door. The command will succeed only if the door is currently open.

- **PUT** `/api/v1/trip/<door name>`: Activates the lift controller for the specified door. The command will succeed regardless of the door's current state. This is useful in cases where the door might be stuck, or its state is not being read correctly.

- **PUT**  `/api/v1/lock/<door name>`: This enables a soft lock on the specified door, preventing any open, close, or trip commands from acting on it. This is a useful way to disable Porter if necessary.

- **PUT**  `/api/v1/unlock/<door name>`: Removes a lock from the specified door


##### Read door state

```
curl -H "Authorization: Bearer $PORTER_TOKEN" http://$PORTER_IP/api/v1/list
```

##### Open the door

```
curl -H "Authorization: Bearer $PORTER_TOKEN" http://$PORTER_IP/api/v1/open/$PORTER_DOOR
```


##### Close the door

```
curl -H "Authorization: Bearer $PORTER_TOKEN" http://$PORTER_IP/api/v1/close/$PORTER_DOOR
```

As you can see, the API is quite simple to use and should be easy to integrate into other projects and solutions like Homeassistant.

Included in this repository is a fully-functional API client library, `porter/client`. You can use this library to add Porter support to your other Go projects, as I've done with the [HomeKit integration](https://github.com/ctrezevant/hkporter), [web interface](https://github.com/ctrezevant/porter-web), and [Twilio-based status monitor](https://github.com/ctrezevant/reporter). Be sure to check out those projects if you're interested in their functionality, or are seeking more examples of how to use the API and client library.
