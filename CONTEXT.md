# Vocabulary

Words as this project uses them.

- **Board**: the currently serving leadership, who ask for advice on current matters.
- **Advisor**: someone, usually a former Board member, who knows history the written record doesn't explain.
- **Harvest**: the half of Agneta that draws knowledge out of Sources and Advisors. The MVP is its start.
- **Counsel**: the half that answers the Board's questions, citing what the Harvest built. Later a
  Mattermost bot using MCP tools. Not in the MVP.
- **Case**: one historical matter the organization handled, e.g. a collaboration ended or a contract and
  what it left out. The basic unit of knowledge.
- **Source**: existing recorded material: the shared record by default, private messages only with consent.
- **Account**: what one Advisor said about a Case, dated and approved by them. Accounts are never merged;
  conflicting ones are kept side by side.
- **Gap**: an open question about a Case that the Sources don't answer. It belongs to the Case: once filled,
  it is filled for everyone.
- **Visibility**: who may see an Account, set by its Advisor: Board, Restricted or Sealed.
- **Session**: one sitting in which an Advisor is interviewed. The first is in person.
- **Offboarding**: the agreed end of interviewing an Advisor, confirmed by the Advisor and the Board.
- **Interview skill**: the agent skill that runs a Session and writes the report.
- **Report**: one JSON file per interview, a flat list of records.
- **Record**: one atomic question–answer pair with its date. How records relate to Accounts and Cases is
  not decided yet.
- **Author**: who uploaded a report, taken from their sign-in. Stored for traceability, never shown to
  agents or the Board. This conflicts with the Counsel naming who said each Account; still open.
- **Upload service**: this repo's backend. Validates reports and writes records to Postgres; later also
  serves the MCP tools.
