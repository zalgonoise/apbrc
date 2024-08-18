## apb:rc

*A configuration files modifier for [APB: Reloaded](https://gamersfirst.com/apb/).*

______________


### Context

From time to time, even on small patches to the game client, it erases (some of) its modified configuration files.

Sure, these configuration files were not modified in-game, but by manually editing the `.ini` files in the game's 
installation folder; but it becomes tiring to edit these files every time, with the same changes.

I don't like keeping a "backup config" folder. These are usually folders with certain configuration files that you 
simply copy over and replace whatever was there on the first place. While these respect the target directories, they
do not consider changes to the game's configuration (so they can break any time).

____________

### Solution

This repository contains a Go application, which can be executed to apply changes to the game's configuration files.

Its target use-case is in a command-line interface environment (executing the binary with different options set as 
command-line flags), or simply by executing it if placed under the game's installation folder (such as next to APB.exe).

The app does not execute the game. It goes through its configured modifiers, applies the changes to the files, and 
exits. If executed from the command-line (cmd, or powershell), you are able to read the output of each change and 
execution.

_____________

### Installation and Usage

This is a drop-in executable that you can run from any directory in the machine where the game is installed, provided 
you point to the right folder. It can also be executed by double-clicking the `.exe` if placed anywhere under the 
`APB Reloaded` folder, where the default application values are applied.

#### Compiling from source

You can compile the binary from source by cloning the repository and using the Go compiler to build the executable. 
For this, you will need both `git` and `go` installed on your computer. Then, execute:

```shell
git clone https://github.com/zalgonoise/apbrc \
&& cd apbrc \
&& go build -o .\build .\cmd\apbrc
```

These commands will clone the repository into the folder that you're in, in the terminal, open this project, and build 
the `.exe` under the `\build` directory. References use Windows paths (with backslashes \) as this is a Windows 
executable.

#### Direct download

Releases of this binary are published on Github, which you can download and use directly. These will be exactly like the 
ones generated in the previous step.


_____________

### [Modifiers](./processor/modifiers)

The modifiers are configurable and the application can be further extended to have more features and functionalities. 
Below is a list of the supported modifiers

#### [Frame Rate Modifier](./processor/modifiers/engine/fps.go)

This modifier changes the limits for the frame rate limiters. Below are the configuration values that change this
behavior:

__Target__: `\Engine\Config\BaseEngine.ini`

|          Key           | Default |                                       Description                                        |
|:----------------------:|:-------:|:----------------------------------------------------------------------------------------:|
| `MinSmoothedFrameRate` |   22    | Changes the minimum frame rate when the `Smoothed` option is set, in the game's Settings |
| `MaxSmoothedFrameRate` |   100   | Changes the maximum frame rate when the `Smoothed` option is set, in the game's Settings |
|  `MaxClientFrameRate`  |   128   |     Changes the game client's maximum frame rate when the `Smoothed` option is unset     |


Below are the configuration flags that modify these values, alongside their default values applied when you don't use
or specify that configuration option and value.

|  Flag  | Type  | Default |                                  Description                                   |
|:------:|:-----:|:-------:|:------------------------------------------------------------------------------:|
| `-cap` | `int` |    0    |       Frame rate limit value to set when the Smoothed option is disabled       |
| `-min` | `int` |   22    | Minimum frame rate value to set when the Smoothed frame rate option is enabled |
| `-max` | `int` |   128   | Maximum frame rate value to set when the Smoothed frame rate option is enabled |


#### [Input Bindings Modifiers](./processor/modifiers/input)

These modifiers change the behavior of the Sprint and Crouch actions.

[Sprint lock](./processor/modifiers/input/sprint_lock.go) causes it to toggle __off__ when held down, which is the opposite behavior of the key. It allows to always 
sprint without holding the assigned key for it.

[Crouch lock](./processor/modifiers/input/crouch_lock.go) causes it to behave as a press-and-hold key; instead of an on-off switch. It can be still locked in crouch 
mode by combining the crouch key and hitting the jump key, for instance, and unlocked by pressing the jump key again.

__Target__: `\APBGame\Config\DefaultInput.ini`

|            Key             |                                       Default                                        |                                       Description                                       |
|:--------------------------:|:------------------------------------------------------------------------------------:|:---------------------------------------------------------------------------------------:|
| `+Bindings=(Name="Sprint"` | `+Bindings=(Name="Sprint",Command="InputSprinting \| OnRelease InputStopSprinting")` | Changes the behavior of the Sprint action, on what happens on key-press and key-release |
|  `+Bindings=(Name="Duck"`  |     `+Bindings=(Name="Duck",Command="Button m_bDuckButton \| InputToggleDuck")`      | Changes the behavior of the Crouch action, on what happens on key-press and key-release |


Below are the configuration flags that modify these values. If `-lock-sprint` and `-hold-crouch` are both unset, no changes occur. 
If the caller wants to reset the input configuration mods, they should use the `-reset-input` option instead.

|      Flag      |  Type  | Default |                     Description                     |
|:--------------:|:------:|:-------:|:---------------------------------------------------:|
| `-lock-sprint` | `bool` |   N/A   |      Sets the input bindings to always sprint       |
| `-hold-crouch` | `bool` |   N/A   | Sets the input bindings to press-and-hold to crouch |
| `-reset-input` | `bool` |   N/A   |       Resets any input bindings modifications       |

_______

### Issues and Suggestions

If you have suggestions, ideas or contributions, please browse the Issues section. Feel free to open a new one if your 
topic has not yet been mentioned.