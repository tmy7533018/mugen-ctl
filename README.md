# mugen-ctl

CLI control interface for [mugen-shell](https://github.com/tmy7533018/mugen-shell), a Hyprland desktop shell built with QuickShell.

## Requirements

- [mugen-shell](https://github.com/tmy7533018/mugen-shell) configured and running
- Hyprland
- `swww`, `mpvpaper` — wallpaper backends
- `grim`, `slurp` — screenshot tools
- `matugen` — color palette generator (optional)

## Install

```sh
go install github.com/tmy7533018/mugen-ctl@latest
```

## Usage

```sh
mugen-ctl ipc <mode>           # send command to mugen-shell
mugen-ctl wallpaper set <path> # set image or video wallpaper
mugen-ctl wallpaper get        # print current wallpaper path
mugen-ctl blur list            # list blur presets
mugen-ctl blur current         # show active preset
mugen-ctl blur set <name>      # apply a preset
mugen-ctl idle toggle          # toggle hypridle on/off
mugen-ctl screenshot           # take a region screenshot
mugen-ctl sddm randomize       # set random SDDM login background
```

### IPC modes

`launcher` `calendar` `wallpaper` `music` `notification` `powermenu` `volume`
`window-switcher` `window-switcher-next` `window-switcher-prev` `close`
