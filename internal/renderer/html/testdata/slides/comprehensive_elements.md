---
title: Comprehensive Elements Test
theme: clean
paginate: true
header: Goslide Presentation Test
footer: Confidential - Cloit Inc.
---

<!-- _class: lead -->
# Goslide Expression Showcase
High performance pure Go presentation engine

---

## 1. Table Alignment & Headings

| Feature (Left) | Status (Center) | Priority (Right) |
| :--- | :---: | ---: |
| Table Align | Supported | High |
| Task Lists | Active | Medium |
| Kbd & Del | Complete | Low |

##### H5 Sub-heading Level 5
###### H6 Sub-heading Level 6

---

## 2. GFM Task Lists, Del & Kbd

- [x] Phase 1: Core CSS alignment fixes
- [x] Phase 2: GFM task lists and keycaps
- [ ] Phase 3: Future enhancements

Text formatting test: ~~deprecated approach~~ and press <kbd>Ctrl</kbd> + <kbd>K</kbd> to search.

* Level 1 Bullet
  * Level 2 Nested Bullet
    * Level 3 Deeply Nested Bullet

---

<!--
_class: invert
-->

# Inverted Slide (Dark Mode Accent)

This slide demonstrates the safe `invert` class on Clean theme.

```go
package main

import "fmt"

func main() {
    fmt.Println("Safe color inversion without contrast bugs!")
}
```

---

<!--
_backgroundImage: ./sample.jpg
_backgroundDim: 0.6
-->

# Slide with Background Image and Dimming

This slide has a background image with a 60% black dimming layer.
The text remains perfectly legible!

![width:300px center](./sample.jpg)

---

<!--
_class: compact
-->

## 3. Compact Preset & GFM Alerts

> [!NOTE] Design Principle
> Goslide maintains zero raw HTML requirements for slide authors.

> [!WARNING] Production Advice
> Always verify slide contrast when presenting in lit rooms.

* Compact mode scales down fonts and paddings to fit large amounts of dense content safely within the 960x540 canvas without overflow or clipping.
