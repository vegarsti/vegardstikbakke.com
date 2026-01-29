---
title: Contributing to DataFusion
slug: datafusion-contributions
date: 2026-01-01
draft: true
---

In the fall of 2025 I made some contributions to [arrow-rs](https://github.com/apache/arrow-rs/pulls?q=is%3Apr+is%3Amerged+author%3Avegarsti+), the official Rust implementation of arrow, and to [DataFusion](https://github.com/apache/datafusion/pulls?q=is%3Apr+is%3Aclosed+author%3Avegarsti), an extensible query engine also in Rust.
I found the community to be active and incredibly friendly and helpful.

DataFusion:

- [fix: UnnestExec preserves relevant equivalence properties of input](https://github.com/apache/datafusion/pull/16985)
- [Respect null elements in `convert_array_to_scalar_vec`](https://github.com/apache/datafusion/pull/17891)
- [Making `array_reverse` 70% faster](https://github.com/apache/datafusion/pull/18500)
- [Support reverse for ListView](https://github.com/apache/datafusion/pull/18424)
- [Add benchmark for array_reverse](https://github.com/apache/datafusion/pull/18425)
- [Add RunEndEncoded type coercion](https://github.com/apache/datafusion/pull/18561)

arrow-rs:

- [Cast support for RunEndEncoded arrays](https://github.com/apache/arrow-rs/pull/8589)
- [Add benchmark for casting to RunEndEncoded](https://github.com/apache/arrow-rs/pull/8710)
- [perf: Use Vec::with_capacity in cast_to_run_end_encoded](https://github.com/apache/arrow-rs/pull/8726)
- [Add cast support for (Large)ListView <-> (Large)List](https://github.com/apache/arrow-rs/pull/8735)

I was the 27th contributor out of 127 for DataFusion release 51.
TODO: Insert screenshot
