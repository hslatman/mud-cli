# mud-cli

MUD CLI provides multiple utilities to work with [Manufacturer Usage Descriptions](https://datatracker.ietf.org/doc/rfc8520/) (RFC8520).

## Description

Manufacturer Usage Descriptions (MUDs) allow manufacturers of IoT equipment to specify the intended network communication patterns of the devices they manufacture. 
The access control policies described in a MUD File allow network controllers to automatically enforce rules on the device, resulting in devices only being allowed to communicate within the boundaries of the access control policies. 

## Things that can be done

* Fix (most) TODOs ... :-)
* Improve README.md
* Add 'Use' to commands
* Add tests
* Replace calls to fmt with proper logging / output
* Allow the tool to be chained (i.e. use STDIN/STDOUT, pipes, etc.)
* A command for generating MUD files (from pcap or some different way)
* A command for editing MUD files (i.e. metadata)
* A command that initializes a .mud directory inside user HOME, that is used for intermediate storage? If necessary, of course.
...