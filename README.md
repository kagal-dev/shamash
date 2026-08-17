# Shamash

<!-- cspell:words Akkadian Mailgun Sumerian Utu šamaš -->

Shamash is a mail system that delegates transport to commercial e-mail
providers and keeps everything else — mailboxes, filtering, search — on
infrastructure its operator controls.

Named for the Mesopotamian Sun God and lord of justice — *šamaš* in
Akkadian, *Utu* in Sumerian — it is the e-mail component of the
[kagal.dev][kagal] ecosystem.

The Go module is `kagal.dev/shamash`.

## Why

A conventional mail host does two jobs that have very little to do with
each other. One is moving messages across the public Internet: earning a
sending reputation, staying off blocklists, authenticating outbound mail,
and absorbing the abuse that arrives unsolicited. The other is holding
messages, which is to say writing files to a disk.

The first job has become specialist work, and the providers who do it do
it well. The second has not changed in thirty years. They are sold
together anyway, so the usual price of reliable delivery is an archive
kept in somebody else's account, readable on their terms and portable
only as far as their export tooling allows.

Shamash unbundles them. Providers move the mail; Shamash keeps it.

## Sending

Shamash accepts submission over SMTP-AUTH. It authenticates the user,
checks the `From` header against the identities that user is permitted
to assert, and passes the accepted message to a provider — Mailgun, for
instance — which performs the delivery.

Rejecting a forged `From` at submission is what makes the arrangement
safe: the provider sees only traffic Shamash has already vouched for, so
one careless account cannot spend the reputation the whole domain
depends on.

## Receiving

Inbound mail arrives at the provider rather than at Shamash. Shamash
collects it — polling the provider's API, or accepting a webhook where
the API is built to push — and writes each message into the Maildir++ of
the user it belongs to.

Several providers, and several domains, can feed a single mailbox; from
the user's side they are indistinguishable. Sieve runs as each message is
filed, so rules take effect at delivery rather than as a later sweep.

## Storage

Every user gets a Maildir++ tree. Where a question comes up about how
mail should sit on disc — folder naming, flag conventions, where a
message's state is recorded — the answer is whatever Dovecot does.
Dovecot is the reference implementation, not a requirement: it need not
be the process serving the mailbox for its layout to be the one Shamash
writes.

The point is that the store stands on its own. Anything that understands
Maildir++ can read it, back it up, or take over from Shamash entirely,
which is the difference between holding your mail and merely being
allowed to download it.

## Search

Dovecot settles the on-disc questions; it settles this one too. Its
full-text search speaks to Solr, so Shamash presents a Solr-compatible
endpoint and inherits that interface rather than inventing one.

Nothing behind the endpoint is Solr. Each user has a private index built
from bleve for full-text retrieval, an HNSW graph for nearest-neighbour
queries, and bbolt as the store underneath. Keeping the indexes per user
means one mailbox's size or churn is nobody else's problem. It also means
each index has exactly one owner: an index admits a single reader and
writer at a time, never several at once.

## Providers

A provider plugin is the only code that knows a particular vendor's API:
how to enumerate and fetch received mail, how to hand over an outbound
message, and — where the vendor pushes instead of waiting to be polled —
how to register and verify a webhook. Above that boundary everything is
expressed in messages and mailboxes, so adding a vendor does not disturb
the rest.

## The binary

Shamash is a single executable. Serving mail, owning the indexes and
operating the whole thing from a terminal are the same program wearing
different hats — one artefact to install, one to keep current.

Two jobs are worth calling out. It can reprocess a directory: recompile
a user's Sieve script and replay it over mail already delivered, so a
rule written today files each message where it would have gone when that
message first arrived. It can also query a user's index directly.

Neither job opens an index itself. Since an index takes one reader and
writer at a time, whichever process holds it is the sole authority on
it, and a command that wanted its own handle would have to wrestle the
mail server for it. Instead the binary daemonises itself on demand and
serves the index over RPC; commands are clients of that daemon. When the
daemon has been idle long enough it releases the index and exits, so
nothing is held open for a mailbox nobody is using.

## Licence

MIT. See [LICENCE.txt][licence].

[kagal]: https://github.com/kagal-dev/kagal
[licence]: LICENCE.txt
