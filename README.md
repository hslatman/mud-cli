# mud-cli

`mud-cli` provides utilities for working with [Manufacturer Usage Descriptions](https://datatracker.ietf.org/doc/rfc8520/) (RFC8520).

## Description

Manufacturer Usage Descriptions (MUDs) allow manufacturers of IoT equipment to specify the intended network communication patterns of the devices they manufacture. 
The access control policies described in a MUD file allow network controllers to automatically enforce rules on the device, resulting in devices only being allowed to communicate within the boundaries of the access control policies. 

This CLI application provides several utilities for working with MUD files.

## Usage

`mud-cli` contains the following commands:

* read - reads (and validates) a MUD file and prints the contents
* validate - validates a MUD file
* sign - sign a MUD file, creating a new signature file
* verify - verifies a MUD file signature
* view - shows a graphical representation of the MUD using [MUD Visualizer](https://github.com/iot-onboarding/mud-visualizer)

### Examples

```console
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

### Binary Verification

`mud-cli` is signed using [Cosign](https://github.com/sigstore/cosign).
This means that binaries can be verified as follows:

```console
$ cosign verify-blob -key cosign.pub -signature mud-<version>-<arch>.sig mud-<version>-<arch>
Verified OK
```

The public key (`cosign.pub`) is available in the repository. 
Signature files and binaries are available from the [Releases](https://github.com/hslatman/mud-cli/releases) page.

### MUD Visualizer

This project embeds [MUD Visualizer](https://github.com/iot-onboarding/mud-visualizer) for visualization of MUD files.

## Things that can be done ...

... in  no particular order.

* Improve README.md
* Add 'Use' texts with examples to commands
* CI improvements (run on tag, build Docker image incl. Cosign signing, optimize binary size using UPX, embed version information)
* Add self-updater and verifier (using Cosign)
* Add tests
* Fix (most, highest priority) TODOs ... :-)
* Customize / improve the [MUD Visualizer](https://github.com/iot-onboarding/mud-visualizer)? It needs proper attribution to be shown on the page, at least. Can be made to look nicer on the page and also include some more metadata about the MUD file.
* Improve logging / output
* Allow the tool to be chained (i.e. use STDIN/STDOUT, pipes, etc.)?
* A command for generating MUD files (from pcap or some different way)
* A command for editing MUD files (i.e. metadata)
* Allow setting a different location than the user home directory
* ...