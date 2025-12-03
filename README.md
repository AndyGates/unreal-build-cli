Utility to build arguments to run `BuildCookRun` for unreal projects

# Installation

Binaries are available at https://github.com/AndyGates/unreal-build-cli/releases

If you have the go runtime installed you can build and install from source with something like
```
go install github.com/AndyGates/unreal-build-cli@latest
```

# Running
The tool should be run from the terminal, with the working directory being the one with your `.uproject`.
I like to alias `unreal-build-cli` to `ub` for easier running.


# Configuration

## Per Project Config
Create a file named "unreal-build-cli.json" next to your `.uproject` 
This can be used to override/add step options or change defaults etc.

## Default Config
These are the default option sets that could be overriden in a config file

``` json
{
    "ClientOptions": {
        "Options": ["Win64", "PS5", "XSX"],
        "Defaults": [0]
    },
    "ServerOptions": {
        "Options": ["Win64", "Linux"],
        "Defaults": []
    },
    "ConfigurationOptions": {
        "Options": ["Development", "Shipping", "Test", "Debug"],
        "Defaults": [0]
    }
}
```

### Example Config
Here is an example of overriding some of the defaults, and adding an additional cooker option

``` json
{
   "ConfigurationOptions": {
      "Defaults": [ 2 ]
   },
   "StepOptions": {
      "Defaults": [0]
   }
   "AdditionalCookOptions": {
      "Options": ["-MyAdditionalCookArg"]
   } 
}
```
