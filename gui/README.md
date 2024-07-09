# gui

gui is a simple immediate mode UI library in pure Go. It was inspired by Sean Barrett's [demonstration](http://silverspaceship.com/inner/imgui/) from 2005.

## Limitations

Currently library uses CGO, despite claiming otherwise. It supports only X Window System.

There are numerous bugs and misfeatures. For example, if you allow sliders to autoscale their width, they will try to fit their label. But if you turn on numeric display, the label can change widths depending on the value, causing the scrollbar to change widths as the value changes (causing a feedback loop since scrollbar changes widths out from under the mouse pointer!). Also grep for `BUGS` to see more.

The main drawback is a dearth of widgets. No radio buttons, no scrollbars and scrollable regions, no combo boxes, no menus.
