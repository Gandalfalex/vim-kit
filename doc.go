// Package vimkit is a UI-agnostic toolkit for vim-style key handling in
// terminal apps. It has no dependencies, so apps on any bubbletea version
// (or none) can share it.
//
//   - keys: named actions, configurable bindings, count prefixes and
//     multi-key sequences (gg, zz).
//   - motion: list/pager motions (j, ctrl+d, G, H/M/L, zz, ...) and in-line
//     text motions (h, l, w, b, e, 0, $).
//   - visual: line and character selections and the text they cover.
package vimkit
