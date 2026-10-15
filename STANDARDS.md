# Standards and scope

Versions checked on **2026-10-15**. The 3GPP archive directories and the
portal's version records agree, and the versions below are the latest
published:

- TS 23.003: the [archive](https://www.3gpp.org/ftp/Specs/archive/23_series/23.003/)
  has nothing above `23003-k10.zip` (2026-09-28), and the
  [portal](https://portal.3gpp.org/desktopmodules/Specifications/SpecificationDetails.aspx?specificationId=729)
  lists V20.1.0 uploaded on 2026-09-28.
- TS 29.002: the [archive](https://www.3gpp.org/ftp/Specs/archive/29_series/29.002/)
  has nothing above `29002-j10.zip` (2026-01-07), and the
  [portal](https://portal.3gpp.org/desktopmodules/Specifications/SpecificationDetails.aspx?specificationId=1585)
  lists V19.1.0 uploaded on 2026-01-07.

The public change-request lists of both specifications showed no CR approved
for a version that is not published yet.

| Document | Version | Clauses used |
| --- | --- | --- |
| 3GPP TS 23.003, *Numbering, addressing and identification* | V20.1.0 (2026-09) | §§2.2, 2.3 IMSI, MCC and MNC; §6.2.1 IMEI; §6.2.2 IMEISV; §6.3 reserved SVN; §12.1 PLMN identifier; Annex B check digit |
| 3GPP TS 29.002, *Mobile Application Part (MAP) specification* | V19.1.0 (2025-12) | §17.7.8 TBCD-STRING, IMSI, IMEI and PLMN-Id coding; §17.7.13 PLMN-Id reuse in LCS data types |
| ISO/IEC 7812-1, *Identification cards — Identification of issuers — Part 1: Numbering system* | 2017 (edition 5) | the Luhn method that TS 23.003 Annex B refers to |

The clause text was read from these archive files:

| Archive | Date | SHA-256 |
| --- | --- | --- |
| `23003-k10.zip` | 2026-09-28 | `384e8627ab5bd47740c7f6fa19836f1c5c2014acb5b396b7c2b750c233a02453` |
| `29002-j10.zip` | 2026-01-07 | `ce8369c9fd6d8e8d5d2a226b89a0d059e02f72b7e8e85fd00b4328edaa262e64` |

TS 23.003 V20.1.0 §2.2 defines the IMSI as a three-digit MCC, a two- or
three-digit MNC and an MSIN. The six-digit minimum follows from that structure
with at least one MSIN digit; §2.3 states decimal digits and the fifteen-digit
maximum. The MNC length requires configuration (§2.2 NOTE 1).

Section 6.2.1 defines an eight-digit TAC, six-digit SNR and check or spare
digit. The spare digit is zero when transmitted by the MS; the check digit
guards against manual transcription errors and is not transmitted. Section
6.2.2 defines IMEISV as TAC, SNR and a two-digit SVN; §6.3 reserves SVN 99.
The package offers separate sender checks and permits receivers to retain
SVN 99.

Annex B.1 computes the check digit from the fourteen TAC/SNR digits, excluding
SVN. Annex B.2 doubles odd-labelled digits D1, D3, …, D13, sums decimal digits
of the products with the even-labelled digits, and chooses the check digit
to reach a multiple of ten. Annex B.3 gives `260531793113837` (check digit 7).
The implementation follows the explicit algorithm of Annex B.2.

Annex B.1 cites ISO/IEC 7812 for the Luhn method as reference **[3]**, but
§1.1 lists [3] as TS 23.060; ISO/IEC 7812 is reference **[40]**. [40] is a
non-specific reference, for which §1.1 says the latest version applies:
[ISO/IEC 7812-1:2017](https://www.iso.org/standard/70484.html), edition 5,
confirmed in 2022 and current. The title that [40] gives differs from the
title of that edition.

Section 12.1 defines PLMN-Id as MCC || MNC. TS 29.002 V19.1.0 §17.7.8
specifies its three-octet coding: low/high nibbles hold MCC1/MCC2, MCC3/MNC3,
and MNC1/MNC2. Only MNC3 can be filler `1111`, indicating a two-digit MNC.
Section 17.7.13 imports this PLMN-Id for LCS data types. MCC 262, MNC 01
therefore encodes as `62 f2 10`.

Recheck the archive and the portal's version records before a significant
revisit of this package, and record the date and any disagreement between
them here.
