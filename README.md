# mud-cli

MUD CLI provides multiple utilities to work with [Manufacturer Usage Descriptions](https://datatracker.ietf.org/doc/rfc8520/) (RFC8520).

## Description

Manufacturer Usage Descriptions (MUDs) allow manufacturers of IoT equipment to specify the intended network communication patterns of the devices they manufacture. 
The access control policies described in a MUD file allow network controllers to automatically enforce rules on the device, resulting in devices only being allowed to communicate within the boundaries of the access control policies. 

## Usage

The mud-cli contains the following commands:

* read - reads (and validates) a MUD file and prints the contents
* validate - validates a MUD file
* sign - sign a MUD file, creating a new signature file
* verify - verifies a MUD file signature
* view - shows a graphical representation of the MUD using [MUD Visualizer](https://github.com/iot-onboarding/mud-visualizer)

### Examples

```bash
$ ./mud
mud provides several utilities for working with MUD files

Usage:
  mud [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  help        Help about any command
  read        Reads and prints MUD file contents
  sign        Signs a MUD file
  validate    Validates a MUD file to be formatted correctly
  verify      Verifies the signature for a MUD file
  view        Provides a graphical view of a MUD file
```

### MUD Visualizer

This project embeds [MUD Visualizer](https://github.com/iot-onboarding/mud-visualizer) for visualization of MUD files.

## Things that can be done

* Fix (most) TODOs ... :-)
* Improve README.md
* Add 'Use' texts to commands
* Builds for other platforms (currently only Darwin, Docker would be great; might need some changes or guidance in terms of files)
* Add tests
* Customize / improve the [MUD Visualizer](https://github.com/iot-onboarding/mud-visualizer)? It needs proper attribution, at least.
* Add some more logging (with levels)
* Replace calls to fmt with proper logging / output
* Allow the tool to be chained (i.e. use STDIN/STDOUT, pipes, etc.)
* A command for generating MUD files (from pcap or some different way)
* A command for editing MUD files (i.e. metadata)
* A command that initializes a .mud directory inside user HOME, that is used for intermediate storage? If necessary, of course.
...