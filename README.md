# I Don't Need This

![Demo](doc/screenshot.gif)

idnt is a tool for batch removal of applications

## Supported operating systems
* Windows
    * Chocolatey
    * Scoop
* macOS
    * Homebrew (needs [rmtree](https://github.com/beeftornado/homebrew-rmtree))
    * Homebrew Cask
* Arch Linux

## Installation
1. Install [fzf](https://github.com/junegunn/fzf)
2. Install idnt: `go get github.com/r-darwish/idnt`

## Usage
Run `idnt`. You'll be presented with an FZF screen containing your installed applications. Press tab to select an application and enter to approve your selection.

**TIP:** You should check out fzf's documentation