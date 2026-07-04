# naver-privacy — API discovery notes

Source: the Naver Privacy Center Next.js app (https://privacy.naver.com).

The app prefetches page content through a CMS API discovered from the
react-query hydration state embedded in the server-rendered flight payload
(queryKey `["page","TRANSPARENCY_REPORT_STATISTICS"]`):

- `GET /api/pages/{pageCode}` → page content + `specificAreaJson`.
  For `TRANSPARENCY_REPORT_STATISTICS` the payload's `statistics` array holds
  28 half-year entries (2012–2025 2H) with per-category request/processing/
  provided/average counts and compliance rates.
  For `NAVER_REPORT` it holds `privacyWhitepapers` + `personalInfoReports`
  (the downloadable publication library).
- `GET /api/pages/intro/{introCode}` → section intros (REPORT, ACTIVITY
  verified live; other codes 400).
- `GET /api/notices?page=0&size=N` (0-indexed) and `GET /api/notices/{id}`.

Live evidence gathered 2026-07-02 (curl): all endpoints unauthenticated,
Korean-language; unknown codes → HTTP 400 `{"errorCode":"BAD_REQUEST"}`.
Counts are strings; "-" marks categories Naver no longer complies with (no
communication user information provided since October 2012).

No existing CLI/SDK wraps this surface (searched GitHub, npm, PyPI for
naver privacy / transparency clients — only search/maps SDKs exist).
