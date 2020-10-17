---
id: milestones
title: Milestones and Roadmap
---

## [v0.7.0-alpha.1](https://github.com/zzpu/openuser/milestone/9)

_This milestone does not have a description._

### [Bug](https://github.com/zzpu/openuser/labels/bug)

Something is not working.

#### Issues

- [ ] Do not create system errors on duplicate credentials when linking oidc providers ([kratos#694](https://github.com/zzpu/openuser/issues/694))

### [Feat](https://github.com/zzpu/openuser/labels/feat)

New feature or request.

#### Issues

- [ ] Selfservice account deletion ([kratos#596](https://github.com/zzpu/openuser/issues/596))
- [ ] Implement Hydra integration ([kratos#273](https://github.com/zzpu/openuser/issues/273))
- [ ] Self-service GDPR identity export ([kratos#658](https://github.com/zzpu/openuser/issues/658))
- [ ] Admin/Selfservice session management ([kratos#655](https://github.com/zzpu/openuser/issues/655))

### [Blocking](https://github.com/zzpu/openuser/labels/blocking)

Blocks milestones or other issues or pulls.

#### Issues

- [ ] Implement Hydra integration ([kratos#273](https://github.com/zzpu/openuser/issues/273))

## [v0.6.0-alpha.1](https://github.com/zzpu/openuser/milestone/8)

_This milestone does not have a description._

### [Bug](https://github.com/zzpu/openuser/labels/bug)

Something is not working.

#### Issues

- [ ] Sending JSON to complete oidc/password strategy flows causes CSRF issues ([kratos#378](https://github.com/zzpu/openuser/issues/378))
- [ ] Unmable to use Auth0 as a generic OIDC provider ([kratos#609](https://github.com/zzpu/openuser/issues/609))
- [ ] Password reset emails sent twice by each of the two kratos pods in my cluster ([kratos#652](https://github.com/zzpu/openuser/issues/652))
- [ ] Building From Source fails ([kratos#711](https://github.com/zzpu/openuser/issues/711))

### [Feat](https://github.com/zzpu/openuser/labels/feat)

New feature or request.

#### Issues

- [ ] Implement Security Questions MFA ([kratos#469](https://github.com/zzpu/openuser/issues/469))
- [ ] Feature request: adjustable thresholds on how many times a password has been in a breach according to haveibeenpwned ([kratos#450](https://github.com/zzpu/openuser/issues/450))
- [ ] Do not send credentials to hooks ([kratos#77](https://github.com/zzpu/openuser/issues/77)) - [@hackerman](https://github.com/aeneasr)
- [ ] Implement immutable keyword in JSON Schema for Identity Traits ([kratos#117](https://github.com/zzpu/openuser/issues/117))
- [ ] Add filters to admin api ([kratos#249](https://github.com/zzpu/openuser/issues/249))
- [ ] Feature Request: Webhooks ([kratos#271](https://github.com/zzpu/openuser/issues/271))
- [ ] Support email verification paswordless login ([kratos#286](https://github.com/zzpu/openuser/issues/286))
- [ ] Prevent account enumeration for profile updates ([kratos#292](https://github.com/zzpu/openuser/issues/292)) - [@hackerman](https://github.com/aeneasr)
- [ ] Support remote argon2 execution ([kratos#357](https://github.com/zzpu/openuser/issues/357)) - [@hackerman](https://github.com/aeneasr)
- [ ] Implement identity state and administrative deactivation, deletion of identities ([kratos#598](https://github.com/zzpu/openuser/issues/598)) - [@hackerman](https://github.com/aeneasr)
- [ ] SMTP Error spams the server logs ([kratos#402](https://github.com/zzpu/openuser/issues/402))
- [ ] Gracefully handle CSRF errors ([kratos#91](https://github.com/zzpu/openuser/issues/91)) - [@hackerman](https://github.com/aeneasr)
- [ ] How to sign in with Twitter ([kratos#517](https://github.com/zzpu/openuser/issues/517))
- [ ] Add ability to import user credentials ([kratos#605](https://github.com/zzpu/openuser/issues/605)) - [@hackerman](https://github.com/aeneasr)
- [ ] Throttling repeated login requests ([kratos#654](https://github.com/zzpu/openuser/issues/654))
- [ ] Require identity deactivation before administrative deletion ([kratos#657](https://github.com/zzpu/openuser/issues/657))
- [ ] Add return_to after logout ([kratos#702](https://github.com/zzpu/openuser/issues/702)) - [@Patrik](https://github.com/zepatrik)
- [ ] Write CLI helper for recommending Argon2 parameters ([kratos#723](https://github.com/zzpu/openuser/issues/723)) - [@Patrik](https://github.com/zepatrik)
- [ ] Add possibility to configure the "claims" query parameter in the auth_url of OIDC providers to request individial id_token claims ([kratos#735](https://github.com/zzpu/openuser/issues/735))

### [Docs](https://github.com/zzpu/openuser/labels/docs)

Affects documentation.

#### Issues

- [ ] Document that identity information (traits, etc) are available to token holders and backend systems ([kratos#43](https://github.com/zzpu/openuser/issues/43)) - [@hackerman](https://github.com/aeneasr)
- [ ] Config JSON Schema needs example values ([kratos#179](https://github.com/zzpu/openuser/issues/179)) - [@hackerman](https://github.com/aeneasr)
- [ ] Elaborate on security practices against DoS and Brute Force ([kratos#134](https://github.com/zzpu/openuser/issues/134))
- [ ] Building From Source fails ([kratos#711](https://github.com/zzpu/openuser/issues/711))

### [Rfc](https://github.com/zzpu/openuser/labels/rfc)

A request for comments to discuss and share ideas.

#### Issues

- [ ] Introduce prevent extension in Identity JSON schema ([kratos#47](https://github.com/zzpu/openuser/issues/47))

## [v0.5.0-alpha.1](https://github.com/zzpu/openuser/milestone/5)

This release focuses on Admin API capabilities

### [Bug](https://github.com/zzpu/openuser/labels/bug)

Something is not working.

#### Issues

- [ ] Refresh Sessions Without Having to Log In Again ([kratos#615](https://github.com/zzpu/openuser/issues/615)) - [@hackerman](https://github.com/aeneasr)
- [ ] Fetching a settings request after error is missing identity data ([kratos#689](https://github.com/zzpu/openuser/issues/689)) - [@hackerman](https://github.com/aeneasr)
- [x] Generate a new UUID/token after every interaction ([kratos#236](https://github.com/zzpu/openuser/issues/236)) - [@hackerman](https://github.com/aeneasr)
- [x] UNIQUE constraint failure when updating identities via Admin API ([kratos#325](https://github.com/zzpu/openuser/issues/325)) - [@hackerman](https://github.com/aeneasr)
- [x] Can not update an identity using PUT /identities/{id} ([kratos#435](https://github.com/zzpu/openuser/issues/435))
- [x] Verification email is sent after password recovery ([kratos#578](https://github.com/zzpu/openuser/issues/578)) - [@hackerman](https://github.com/aeneasr)
- [x] Do not return expired sessions in `/sessions/whoami` ([kratos#611](https://github.com/zzpu/openuser/issues/611)) - [@hackerman](https://github.com/aeneasr)
- [x] Logout does not use new cookie domain setting ([kratos#645](https://github.com/zzpu/openuser/issues/645))
- [x] Email field type changes on second request for request context during registration flow ([kratos#670](https://github.com/zzpu/openuser/issues/670))
- [x] Segmentation fault when running kratos ([kratos#685](https://github.com/zzpu/openuser/issues/685)) - [@Patrik](https://github.com/zepatrik)
- [x] Endpoint whoami returns valid session after user logout ([kratos#686](https://github.com/zzpu/openuser/issues/686)) - [@hackerman](https://github.com/aeneasr)

#### Pull Requests

- [x] fix: escape jsx characters in api documentation ([kratos#703](https://github.com/zzpu/openuser/pull/703)) - [@hackerman](https://github.com/aeneasr)
- [x] fix: mark flow methods' fields as required ([kratos#708](https://github.com/zzpu/openuser/pull/708)) - [@hackerman](https://github.com/aeneasr)

### [Feat](https://github.com/zzpu/openuser/labels/feat)

New feature or request.

#### Issues

- [ ] Implement React SPA sample app ([kratos#668](https://github.com/zzpu/openuser/issues/668)) - [@hackerman](https://github.com/aeneasr)
- [ ] Implement React Native sample application consuming API ([kratos#667](https://github.com/zzpu/openuser/issues/667)) - [@hackerman](https://github.com/aeneasr)
- [ ] Rename strategy to method in internal APIs and Documentation ([kratos#683](https://github.com/zzpu/openuser/issues/683)) - [@hackerman](https://github.com/aeneasr)
- [ ] Configurable CORS headers ([kratos#712](https://github.com/zzpu/openuser/issues/712)) - [@hackerman](https://github.com/aeneasr)
- [x] Implement JSON capabilities in ErrorHandler ([kratos#61](https://github.com/zzpu/openuser/issues/61)) - [@hackerman](https://github.com/aeneasr)
- [x] Allow attaching credentials to identities in CRUD create ([kratos#200](https://github.com/zzpu/openuser/issues/200))
- [x] Move away from UUID-based challenges and responses ([kratos#241](https://github.com/zzpu/openuser/issues/241)) - [@hackerman](https://github.com/aeneasr)
- [x] Add tests to prevent duplicate migration files ([kratos#282](https://github.com/zzpu/openuser/issues/282)) - [@Patrik](https://github.com/zepatrik)
- [x] Session cookie (ory_kratos_session) expired time should be configurable ([kratos#326](https://github.com/zzpu/openuser/issues/326)) - [@hackerman](https://github.com/aeneasr)
- [x] Can not update an identity using PUT /identities/{id} ([kratos#435](https://github.com/zzpu/openuser/issues/435))
- [x] Make session cookie 'domain' property configurable ([kratos#516](https://github.com/zzpu/openuser/issues/516))
- [x] Remove one of in-memory/on-disk SQLite e2e runners and replace with faster test ([kratos#580](https://github.com/zzpu/openuser/issues/580)) - [@Andreas Bucksteeg](https://github.com/tricky42)
- [x] Password similarity policy is too strict ([kratos#581](https://github.com/zzpu/openuser/issues/581)) - [@Patrik](https://github.com/zepatrik)
- [x] Implement a test-error for implementing the Error UI ([kratos#610](https://github.com/zzpu/openuser/issues/610))
- [x] Design of the client cli ([kratos#663](https://github.com/zzpu/openuser/issues/663)) - [@Patrik](https://github.com/zepatrik)
- [x] Rename `request_lifespan` to `lifespan` ([kratos#666](https://github.com/zzpu/openuser/issues/666)) - [@hackerman](https://github.com/aeneasr)

#### Pull Requests

- [ ] feat: prepare v0.5.0 release ([kratos#736](https://github.com/zzpu/openuser/pull/736)) - [@hackerman](https://github.com/aeneasr)
- [x] fix: resolve identity admin api issues ([kratos#586](https://github.com/zzpu/openuser/pull/586)) - [@hackerman](https://github.com/aeneasr)
- [x] feat: implement API-based self-service flows ([kratos#624](https://github.com/zzpu/openuser/pull/624)) - [@hackerman](https://github.com/aeneasr)

### [Docs](https://github.com/zzpu/openuser/labels/docs)

Affects documentation.

#### Issues

- [x] Document multi-tenant set up ([kratos#370](https://github.com/zzpu/openuser/issues/370))
- [x] Remove reverse proxy from node example and rely on port and the domain parameter ([kratos#661](https://github.com/zzpu/openuser/issues/661)) - [@hackerman](https://github.com/aeneasr)

#### Pull Requests

- [ ] feat: prepare v0.5.0 release ([kratos#736](https://github.com/zzpu/openuser/pull/736)) - [@hackerman](https://github.com/aeneasr)

### [Rfc](https://github.com/zzpu/openuser/labels/rfc)

A request for comments to discuss and share ideas.

#### Issues

- [x] Rename login/registration/recovery/... request to flow ([kratos#635](https://github.com/zzpu/openuser/issues/635)) - [@hackerman](https://github.com/aeneasr)

### [Blocking](https://github.com/zzpu/openuser/labels/blocking)

Blocks milestones or other issues or pulls.

#### Issues

- [x] Remove reverse proxy from node example and rely on port and the domain parameter ([kratos#661](https://github.com/zzpu/openuser/issues/661)) - [@hackerman](https://github.com/aeneasr)
- [x] Rename `request_lifespan` to `lifespan` ([kratos#666](https://github.com/zzpu/openuser/issues/666)) - [@hackerman](https://github.com/aeneasr)

#### Pull Requests

- [x] feat: implement API-based self-service flows ([kratos#624](https://github.com/zzpu/openuser/pull/624)) - [@hackerman](https://github.com/aeneasr)
- [x] fix: escape jsx characters in api documentation ([kratos#703](https://github.com/zzpu/openuser/pull/703)) - [@hackerman](https://github.com/aeneasr)
