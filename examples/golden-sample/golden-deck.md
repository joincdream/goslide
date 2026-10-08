---
title: "Modern Large Language Model Architecture"
author: "AI Core Architecture Lab"
theme: "corporate"
size: "16:9"
paginate: true
header: "2026 AI Architecture Tech Seminar"
footer: "© 2026 Cloit Tech Architecture. All rights reserved."
---

<!-- _layout: cover -->
<!-- _paginate: false -->
<!-- _header: "" -->

# Modern Large Language Model Architecture
### From Transformer Core Mechanisms to Scalable Serving

**Speaker**: Cloit AI Architecture Lab  
**Date**: October 2026

<!-- note:
- Welcome attendees and introduce the session.
- Emphasize that this talk covers production-grade Go serving infrastructure beyond pure theory.
-->

---

<!-- _layout: section -->

# 01. Paradigm Shift: From RNN to Transformers
> *"General methods that leverage computation are ultimately the most effective."*

<!-- note:
- Section transition slide.
- Briefly remind the audience of Rich Sutton's Bitter Lesson and empirical scaling laws.
-->

---

<!-- _layout: two-cols -->

## Sequential vs Transformer Architecture

### Sequential Models (RNN / LSTM)
- **Sequential Bottleneck**: Hidden state dependency prevents GPU parallelization
- **Vanishing Context**: Gradient decay degrades long-range attention
- **Training Complexity**: $O(N)$ sequential steps bottleneck large corpora

<!-- split -->

### Transformer Decoder
- **Full Parallelization**: Simultaneous matrix operations across tensor cores
- **Self-Attention Mechanism**: $O(1)$ direct correlation path across all tokens
- **Predictable Scaling Law**: Capability grows reliably with compute budget

<!-- note:
- Guide audience through the left-to-right comparison.
- Highlight the parallelization advantage on the right column.
-->

---

<!-- _backgroundDim: 0.7 -->

## High-Throughput Go Streaming Proxy

- **Ultra-lightweight Goroutines**: Under 4KB baseline memory per active client
- **Native SSE Token Streaming**: Low-latency token delivery without heavy WebSocket overhead
- **Deterministic Single Binary**: Zero CGO dependencies for seamless cross-platform deployment

> [!TIP] Serving Optimization
> Leveraging Go's standard `net/http` and `http.Flusher` allows serving tens of thousands of concurrent streams with zero buffer bloat.

<!-- note:
- Showcase production architectural wins on this highlighted visual slide.
- Highlight the memory efficiency of Go over heavier alternatives.
-->

---

<!-- _fragmentStyle: dim -->

## Go Streaming Core Implementation

- **Low-Latency Streaming**: Immediate token flushing using standard `http.Flusher`
<!-- pause -->
- **Graceful Context Cancellation**: Client disconnects immediately release server resources

```go {3-5}
func StreamHandler(w http.ResponseWriter, r *http.Request) {
    flusher, _ := w.(http.Flusher)
    for token := range stream.Tokens() {
        fmt.Fprintf(w, "data: %s\n\n", token)
        flusher.Flush()
    }
}
```

<!-- note:
- Walk through lines 3 to 5 where the token flush occurs.
- Show how code line highlighting and incremental reveal direct audience attention.
-->
