# Standards and scope

Versions checked against the 3GPP specification archive on **2026-10-15**. The
latest published version of each document is the one listed.

| Document | Version | Clauses used | Archive |
| --- | --- | --- | --- |
| 3GPP TS 23.003, *Numbering, addressing and identification* | V20.1.0 | §2.2 and §2.3 IMSI, §6.2.1 IMEI, §6.2.2 IMEISV, §6.3, Annex B | [23.003](https://www.3gpp.org/ftp/Specs/archive/23_series/23.003/) |
| 3GPP TS 29.002, *Mobile Application Part (MAP) specification* | V19.1.0 | §17.7.8 PLMN-Id coding | [29.002](https://www.3gpp.org/ftp/Specs/archive/29_series/29.002/) |

TS 23.003 defines the structure of each identity: the IMSI as a three-digit MCC,
a two- or three-digit MNC and the MSIN, decimal and at most 15 digits; the IMEI
as TAC, SNR and a check or spare digit, 15 digits; the IMEISV as TAC, SNR and
SVN, 16 digits. TS 29.002 gives the three-octet coding of the PLMN identity.

Recheck the archive for newer versions before a significant revisit of this
package, and update this file with the date of the check.
