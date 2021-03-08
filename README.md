# WoST Gateway Logger

Simple logger of messages on the gateway internal message bus, intended for testing of plugins.


## Objective

Facilitate the development of plugins by logging messages on the gateway internal message bus.


## Status 

The status of this plugin is Alpha.

Basic logging of channel messages to file is functional.


## Audience

This project is aimed at software developers, system implementors and people with a keen interest in IoT. 


## Summary

Gateway plugins communicate TD's, events and actions over the internal message bus. This plugin writes those message channels to file.

This simple plugin also serves as an example on writing plugins.


## Build and Installation

### System Requirements

This plugin runs as part of the WoST gateway. It has no additional requirements other than a working gateway.


### Manual Installation

See the gateway README on plugin installation.


### Build From Source

Build with:
```
make all
```

The plugin can be found in dist/bin for 64bit intel or amd processors, or dist/arm for 64 bit ARM processors. Copy this to the gateway bin or arm directory.
An example configuration file is provided in config/logger.yaml. Copy this to the gateway config directory.

## Usage

Configure the logger.yaml configuration file with the channels to log. The default channels are td, event, action, and plugin channels. See the provided example configuration file for documentation.
