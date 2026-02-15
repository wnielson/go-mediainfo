# Changelog

## v0.1.7 (2026-02-15)
- TS (`.ts`) ParseSpeed=0.5 scan windowing improvements (head/tail + optional mid-window for DTVCC) to better match official `mediainfo`.
- TS jump resets now freeze MPEG-2 GOP info before clearing parser state (avoids GOP inference regressions after seeks).
- Parity notes: `Nickelodeon - Generic Halloween Promo.ts` now diff `0` vs official JSON raw at ParseSpeed=0.5; remaining TS diffs are AC-3 stats window/count edge cases on a small control set.

## v0.1.6
- See git history for details.
