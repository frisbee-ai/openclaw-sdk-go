# Changelog

All notable changes to this project will be documented in this file.

## Unreleased



### 🚀 Features

- Add RequestRateLimiter, TokenBucketLimiter, and MaxRetries ([f36a470](f36a470ba211c65f24a90fc286ca098c8b0bb0b2)) - (Lin Yang)
- Add ErrTooManyPendingRequests and ErrMaxRetriesExceeded typed errors ([fe2d433](fe2d433d880b0ba1265bfd5c75917c43ca85b5fa)) - (Lin Yang)
- Implement MaxRetries enforcement, TLS/Logger wiring, and disconnect-triggered reconnect ([1189b3e](1189b3ed63e76eac862ee60f64b675bc196e67fa)) - (Lin Yang)
- Implement pending request limits and fix channel ownership ([e8f5252](e8f5252a8031320c214607ae24151bbb7d67aae8)) - (Lin Yang)
- Add rate limiter, reduced mutex scope, and rate-limit options ([d0b2234](d0b22349a7a40325e4ff3596610f837ceed43e89)) - (Lin Yang)
- Add ConnectionMetrics struct and EventPriority type ([94ffe22](94ffe2247b97b8af449a1ed91bd0c2b53dc9e063)) - (Lin Yang)
- Add GetTickIntervalMs and GetStaleMultiplier to TickMonitor ([bc48c40](bc48c4043c30e55e3e695919fa54133f9f746c5f)) - (Lin Yang)
- Add AttemptCount() to ReconnectManager with atomic counter ([243f4a3](243f4a35848de3c1095ccb2c67ee78e51778750e)) - (Lin Yang)
- Add GetMetrics() to OpenClawClient interface and client ([e4921b0](e4921b02dcf7cb78c4e70a934b62873713bd4358)) - (Lin Yang)
- OBS-02 and OBS-03 implementation ([5803615](580361595b5b8a589bb16e34ac21b5d8d07708f4)) - (Lin Yang)
- Add benchstat CI job with regression detection ([ae43a95](ae43a953454bc2c8bd048896c872de9bb10a1a80)) - (Lin Yang)
- Add benchmarks and fuzz testing with round-trip assertions ([f3e980e](f3e980e256f53ea362d060db8d1bfc3674505b0f)) - (Lin Yang)

### 🐛 Bug Fixes

- Rename file to files for codecov-action v5 ([6e6f113](6e6f113c641dac704ef12d6690027475af823066)) - (Lin Yang)
- Stabilize flaky reconnect and channel ownership tests ([ebd76e1](ebd76e1cdca2bbc683eeebe8462487a7d4d02553)) - (Lin Yang)



### ♻️ Refactor

- Align Go SDK API with TypeScript SDK (127 methods) ([087d2d2](087d2d2e2f87aa596876878e7ea879a73e8460b7)) - (Lin Yang)
- Reorganize client struct into logical sub-structs ([a797aac](a797aacd1aa886cca79050cb38d1bb602b67f2ef)) - (Lin Yang)

### ✅ Testing

- Add test coverage for 7 new API modules ([4ee490a](4ee490ad2ade7a7adbfbd3c910344bdd58db490d)) - (Lin Yang)
- Extend TestDefaultReconnectConfig to check MaxRetries=10 ([5773424](5773424f152f08442be1f15607265a3247d15ff7)) - (Lin Yang)
- Add OBS-04 EventBufferSize verification tests ([8ff46c7](8ff46c7156e59d1f2b68a595bbe903579c88cba4)) - (Lin Yang)
- Add 24 fuzz corpus files for all frame types ([597e83c](597e83c70f093d2b1dcc2bd50ec7e904002212fa)) - (Lin Yang)



### 📖 Documentation

- Synchronize README with current implementation and add missing config options ([f06077d](f06077d726b2c1ba2a6ae537c32718f5f1b65679)) - (Lin Yang)
- Map existing codebase ([f8050da](f8050da17e397ae49c185fef1706043d38c75c99)) - (Lin Yang)
- Initialize project (brownfield) ([2385b74](2385b74c3dff8b6e68e7e7b037cef6330ca7c5fe)) - (Lin Yang)
- Complete project research ([a59d2b0](a59d2b0d36690c1e699845dd9a6c1b473b08916b)) - (Lin Yang)
- Define v1 requirements ([7510a23](7510a23b3795ac68f1e50faab761d7358cf4ad45)) - (Lin Yang)
- Create roadmap (5 phases) ([3c50ecb](3c50ecb01409bd95ae66b2611e179851be37c06e)) - (Lin Yang)
- Research foundation hardening phase ([965258a](965258a4f19246bda82d9a6220fdc29523bec8d7)) - (Lin Yang)
- Create phase plan ([67bae2a](67bae2ac5764fb89ba578ff1472609a25e8c464e)) - (Lin Yang)
- Cross-AI review for phase 1 ([28fae2b](28fae2bd1d06d6004dd9fe34877ff56415e7142a)) - (Lin Yang)
- Revise plans based on cross-AI review feedback ([f7893ad](f7893addcd02c5916a6279c5f126297fae981618)) - (Lin Yang)
- Complete plan 01 - foundation types and typed errors ([f402fe6](f402fe6acfb4e5b8ece4c903b597c0d1aafd2924)) - (Lin Yang)
- Complete plan 03 - foundation hardening ([1542e6d](1542e6d054deff8f97a7a36595a3afde94bd1e97)) - (Lin Yang)
- Complete plan 02 - foundation hardening ([bf7f383](bf7f38344054270ae7cedaada682ccdc58c8a76f)) - (Lin Yang)
- Complete foundation hardening ([15c1170](15c117018a2e3e2c68bfe12d0354ae945f07d9dc)) - (Lin Yang)
- Capture phase 2 observability context ([1f53f37](1f53f371b4e3cf8d8a1cdf0bc1bb28c6ad373bda)) - (Lin Yang)
- Record phase 2 context session ([c9928a9](c9928a937ac80ed45a56f264a24ea65e173be9fa)) - (Lin Yang)
- Research observability phase domain ([65a0edb](65a0edbd690df74b8026c5bdbea5a624768d4c63)) - (Lin Yang)
- Add research and validation strategy ([05b65a3](05b65a319e226e6e0634f785ee02d365263506f1)) - (Lin Yang)
- Create observability phase plan (3 plans) ([1b63e5e](1b63e5e1ef2044130b5d064ce776cc69df8fd0b1)) - (Lin Yang)
- Complete 02-observability-01 plan ([87990b7](87990b7b257545abd2f76cfdc4aec8cee8542a5b)) - (Lin Yang)
- Complete 02-observability-02 plan ([bc0fe0b](bc0fe0b4193022e13da69def362146becaed64fe)) - (Lin Yang)
- Update roadmap progress for 02-observability-02 ([b61ba3e](b61ba3e1b3a6c297fb947a788265e476a9f6ddfe)) - (Lin Yang)
- Complete observability phase ([ccb088d](ccb088d2c19c4a4d5d168271aa122b467a62eb68)) - (Lin Yang)
- Complete observability phase verification ([3916d3b](3916d3b206f366f4e48031222821fbbcf019c263)) - (Lin Yang)
- Sync GSD workflow artifacts from Phase 1 and Phase 2 ([783a325](783a325f7743c8c94a40922922b8eeaab6d1fa93)) - (Lin Yang)
- Capture client struct refactor context ([2f39af0](2f39af0bf587ff056af82ae2bcfcd05077566a2d)) - (Lin Yang)
- Record phase 3 context session ([9948291](9948291a98de8f23b4097ee83662b3a8be4d95e6)) - (Lin Yang)
- Create phase plan ([e87f667](e87f6677117fee833dc6197199c7501bd81f7e70)) - (Lin Yang)
- Create client struct refactor plans ([dab4d1f](dab4d1fd89af5cd22ad39f76f7b88eb5756f6672)) - (Lin Yang)
- Complete client struct refactor plan ([1d93ddd](1d93ddd5bf4b7eb8f352379680c6a9f541e7d8dd)) - (Lin Yang)
- Update Close/Disconnect interface docs per D-03 ([07d96a3](07d96a3e2449af40da03ca698a92f36c2ac3c6a3)) - (Lin Yang)
- Complete 03-02 plan ([812bbd8](812bbd8d91770a6b7584f82e7e08460227b871b8)) - (Lin Yang)
- Complete phase execution ([a1dc360](a1dc360af1a402153e52fecdfdadcfc8cc0297f1)) - (Lin Yang)
- Evolve PROJECT.md after phase completion ([234f386](234f386794c65f6299a8e31a20913fa2c7c92000)) - (Lin Yang)
- Capture phase 04 context ([92dc9ed](92dc9ed198c6713b9101211d3cfa1683ca96a678)) - (Lin Yang)
- Record phase 04 context session ([970aaa2](970aaa2f8bb2d746701f83e34dbb82b981cdd172)) - (Lin Yang)
- Create benchmarking and fuzz testing plans ([24bb026](24bb026b12417fe84c887303bf862409f23ff8fa)) - (Lin Yang)
- Update Phase 04 progress and mark TEST-03 complete ([4dfd5d9](4dfd5d95ba8cdf8056878dbc204c3379389ac724)) - (Lin Yang)
- Complete benchmarking and fuzz testing phase ([c9885a6](c9885a6b0c5f8dc7b389b8d3ca6e9600a329a193)) - (Lin Yang)
- Complete benchmarking and fuzz testing plan ([3795e9a](3795e9a6833413fc46feac2eae0e47d1b17677f6)) - (Lin Yang)
- Complete phase execution ([910966b](910966bbcab34cb36b46d9e2f403292feb4fb916)) - (Lin Yang)
- Evolve PROJECT.md after phase completion ([11c365d](11c365da0d44153845955e04ac04e3104726a954)) - (Lin Yang)
- Capture release infrastructure context ([03d0968](03d0968dd24f5b6ae86c7c0412a894d3116e454b)) - (Lin Yang)
- Record phase 05 context session ([d7e7114](d7e711426748b5731db204373e6074ad28e1af43)) - (Lin Yang)
- Add release infrastructure research ([0cc7b06](0cc7b06a23f620b3061949cc047aa0da34af2ff2)) - (Lin Yang)
- Add research and validation artifacts ([7591900](7591900242850cbaf3caea24b08f77ca315570d5)) - (Lin Yang)
- Create release infrastructure phase plan ([a6495ed](a6495edcd5432c488a98a2244a16282df71e88de)) - (Lin Yang)
- Complete 05-01 plan - GoReleaser v2 configuration ([ee36b93](ee36b9369b83dceee962a9022f20e9789239039d)) - (Lin Yang)
- Complete 05-03 plan - verify git-cliff configuration ([08f929e](08f929efde823034d707f1ed20438b91c7d39cd2)) - (Lin Yang)
- Complete release infrastructure - goreleaser config, git-cliff verified, skip v1.0.1 tag (REL-02 partial) ([cd74799](cd74799525ece96e7d1e7e62636afa3b65ba2678)) - (Lin Yang)
- Evolve PROJECT.md after release infrastructure completion ([e0a7060](e0a706003dcc9e4cdd6f7043fa7a524e59a29286)) - (Lin Yang)




## v1.0.0(2026-03-22)




### 🚀 Features

- Add common types and error type hierarchy ([1396bbb](1396bbbde374deea6939e8ddcec9eecb1c4d3387)) - (Lin Yang)
- Add Logger interface with context support ([d1d28ba](d1d28baf06bd832fe64f791e3a97518ce0076e5e)) - (Lin Yang)
- Add auth module with CredentialsProvider and AuthHandler ([16a2b15](16a2b15c1d6310994adfbf94ea336aebf017947b)) - (Lin Yang)
- Add protocol module with types and validation ([6040de1](6040de1e2888ac0b750e3c627b040dece09807f4)) - (Lin Yang)
- Add WebSocket transport module ([9d04235](9d042352b058fad02d23206c66d013cdacb7ca71)) - (Lin Yang)
- Add connection module with state machine and policies ([f066db5](f066db53751ff14706776be68f66724d791b078d)) - (Lin Yang)
- Add events module with tick monitor and gap detector ([044d7be](044d7bea339428cbc84fd831cda6e70a1f896eee)) - (Lin Yang)
- Add managers module with event, request, connection, and reconnect managers ([672bf8b](672bf8b3358380eb4086bd24b0dc653a990315ba)) - (Lin Yang)
- Add timeout manager ([1e1161b](1e1161b4a971b2823dd4249677bf769e5cd51ca6)) - (Lin Yang)
- Add main client with options and reorganize package structure ([162e528](162e528f92822375e8b87e1bd04b551e700c2b55)) - (Lin Yang)
- Add CLI example ([c95b9dc](c95b9dc9789612172bb28a0631547f4f18101c4b)) - (Lin Yang)
- Add WebSocket echo server example ([758ae6f](758ae6fb02b5d13802eaba0332a46386ae3454fc)) - (Lin Yang)
- Enhance TLS certificate validation with comprehensive security checks ([5fa31b4](5fa31b4023f6ef37d1a23798690e0d90b9c49da4)) - (Lin Yang)
- Implement new features of TypeScript to Go SDK migration ([3d24935](3d24935cc4aac393acf48b8d3240c2231e09ff7c)) - (Lin Yang)
- Add context support to Dial operation ([b53e978](b53e9783defff4a9bdb9868fecd2b5fe4c55bf53)) - (Lin Yang)
- Add backpressure timeout for EventManager Emit operations ([3148bfc](3148bfc6e334caeecfa215e185450acca5cea0fb)) - (Lin Yang)
- Add payload size validation against server policy ([9c93c48](9c93c480fab36a4f38ae0fa115b97cd670080cfc)) - (Lin Yang)
- Integrate git-cliff for structured release notes generation ([1bff317](1bff317d110368f8608154dd09a3d84c29fc6749)) - (Lin Yang)

### 🐛 Bug Fixes

- Prevent memory leak in ReconnectManager by using time.NewTimer ([da1d24a](da1d24ad4497eb35b8460b3d051f283f53cb80d4)) - (Lin Yang)
- Remove infinite loop in ProtocolNegotiator and simplify version matching ([47f3a73](47f3a73e3f034617e27a5752c1655d76f37c466e)) - (Lin Yang)
- Properly handle errors in connection manager and transport ([d204694](d2046944a29c892995ee345b84384c1c35619fd6)) - (Lin Yang)
- Properly load TLS certificates in transport layer ([373d972](373d97266e0173df6269f451607671bafcb5cb4a)) - (Lin Yang)
- Resolve Timer reset race condition in TickMonitor ([d1612f4](d1612f40292937ede2e40e3391e2dd47020d3e86)) - (Lin Yang)
- Log panics in EventManager handlers instead of silent discard ([d070973](d0709736d2c6eb6b3135b51ee37e9eed1d6562b7)) - (Lin Yang)
- Replace unsafe.Pointer with atomic counter in EventManager ([ef79823](ef7982355bc560049f178151c49e237053255a4c)) - (Lin Yang)
- Validate credential values in auth package ([c22a551](c22a5517814b80a8d74c42876261ae3305ecd90b)) - (Lin Yang)
- Enhance method name validation with regex and comprehensive tests ([29f9eb1](29f9eb1c5a85e1244f6c6e860acdd52c19aeacd1)) - (Lin Yang)
- Replace type assertions with json.Unmarshal to prevent runtime panics ([b5e39db](b5e39dba63a5c38958d6fd8b64cdce7d8c044f3e)) - (Lin Yang)
- Implement performHandshake and fix state transition path ([aaacfcc](aaacfccc345acfeabec9556377eeb87eb41ae5a1)) - (Lin Yang)
- Reconnect uses stored params instead of raw Connect ([5f739dc](5f739dc16be8f542e482906731d7ee254c278088)) - (Lin Yang)
- Use strconv.Itoa for protocol version to prevent overflow ([9524606](9524606ffcd4fdaf16d4f1b80fc6c4fb4642f11e)) - (Lin Yang)
- Add closed flag to prevent RequestManager double-close ([6a9023b](6a9023b0d903d1c5b2bc5c1984e05fee1e173833)) - (Lin Yang)
- Use write lock in GetStaleDuration to prevent data race ([0bae076](0bae076286423d151cef725f2b7ef1c07b8bb9a2)) - (Lin Yang)
- Add background goroutine for automatic stale detection ([8401c39](8401c39df4a6e001689f6bf269ef9d1eedee9e15)) - (Lin Yang)
- Simplify IsRequestError and fix Unwrap chain ([a7c471a](a7c471a50431be9332d526629447988c7120fa78)) - (Lin Yang)
- Use crypto/rand for request ID generation ([d30949c](d30949c8754a2f75dc5bf80eccb3028d38b2204d)) - (Lin Yang)
- Add WithClientID option and fix tests missing ClientID ([8a78b06](8a78b061175c91c190d98031dae46f5f8427f837)) - (Lin Yang)
- Make channel buffer size configurable ([e42c623](e42c6237d19187f4d6e2a8d560e26cc4ac5d54cc)) - (Lin Yang)
- Log warning when events are dropped due to full channel ([baff296](baff296143c75e3b68f57eb6259ad737ea9f7fa3)) - (Lin Yang)
- Resolve race conditions in request and event managers ([5853e8e](5853e8ef81097be40aa035a9197bc08a3a2adc36)) - (Lin Yang)
- Log protocol negotiation errors instead of silently discarding ([0652469](0652469b16cba3d58070967c2a7b841cf9ec48c5)) - (Lin Yang)
- Add security warning when InsecureSkipVerify is enabled ([e438a5f](e438a5ffb93e83d72952dcaf54cfe8bdefce9f94)) - (Lin Yang)
- Migrate golangci-lint config to v2 format ([d980789](d98078918532f5fcde5728fbef73509ffe366f4d)) - (Lin Yang)
- Remove invalid --all flag from git-cliff command ([be83719](be83719db83a15bd904dfc4e1eb262f3f8d05c31)) - (Lin Yang)
- Skip build step for library project ([022a142](022a1426423af3eb0c65789c25deaf69a021a1ef)) - (Lin Yang)
- Rename build to builds in goreleaser config ([2194e31](2194e31bd198e3ca30b6c9a60097929e9c253f3e)) - (Lin Yang)
- Fixup! fix(ci): rename build to builds in goreleaser config ([fd15547](fd155473e890e0a57c9fb70e898686fab4ab72f0)) - (Lin Yang)


### ⚡ Performance

- Eliminate redundant JSON marshal and slice allocation ([81b701f](81b701f447ad7d92ef155db49d4a355d303e1b26)) - (Lin Yang)
- Replace time.After with time.NewTimer+Reset in EventManager.Emit ([bac682c](bac682c2db2fbd2b90eaa338f0a9c94c6476d377)) - (Lin Yang)
- Eliminate heap allocation in generateRequestID ([26b8d90](26b8d90a4cbab91a531dc6650114de327b674884)) - (Lin Yang)

### ♻️ Refactor

- Move source files to pkg/openclaw directory ([be15cff](be15cffb1fd981918c70694f0d4c4bee1a8b04fe)) - (Lin Yang)
- Move utils package from pkg/openclaw/utils to pkg/utils ([b5e5bee](b5e5bee3a4a5de5f56260449f28f1b00d6861c1a)) - (Lin Yang)
- Move openclaw package from pkg/openclaw to pkg ([b748cb7](b748cb737656c9e2d82c5e845376efd1f681a0ac)) - (Lin Yang)
- Consolidate re-exports in client.go, move types to pkg/types/ ([e5fffa6](e5fffa6cbf368662d0d3424ea2629685bed782a2)) - (Lin Yang)
- Clean up unused code and fix misleading comment ([e91520a](e91520a5547015b6155b31c14003e833f0447b70)) - (Lin Yang)
- Remove unused stateMachine field ([e97a82c](e97a82ce5a4a9cb9e0a817a0424086086c2238b8)) - (Lin Yang)
- Split api_params.go into domain-specific files ([26cfc60](26cfc60f4d63211d065c89ae560bc86597cf5921)) - (Lin Yang)
- Remove dead setupConnectionHandlers placeholder ([6093aae](6093aaeac930d0b2c0feb071e23856270147a830)) - (Lin Yang)

### ✅ Testing

- Add comprehensive error type tests ([c04a9a6](c04a9a65fc5e0cd832c8f51727d6fb67bd53864a)) - (Lin Yang)
- Enhance test coverage and fix bugs ([4ebfe4a](4ebfe4ab96b182500a3ad37b83fe20ecda479ffa)) - (Lin Yang)
- Add comprehensive tests for all 8 API namespaces ([f22698e](f22698e1466839d8dbf4f896e208d86ab6da623b)) - (Lin Yang)
- Fix MalformedJSON no-op test to actually parse JSON ([24c930c](24c930c9ec90fe79a755198045a1a91c4a3e511a)) - (Lin Yang)
- Improve coverage for types, api, managers, utils, and client packages ([bd53c2d](bd53c2da8ecff05ccfd3bb07ea7bc5c88ef1c56d)) - (Lin Yang)



### 📖 Documentation

- Add TypeScript to Go migration design spec ([2c4de56](2c4de56229b58f4955e36f23e1dd89a1c7371973)) - (Lin Yang)
- Add implementation plan with 10 phases ([68e5644](68e56443f516632e3599723e8b9116862c7d82c3)) - (Lin Yang)
- Improve Phase 1-2 implementation plans ([e864988](e86498889cb022da0bc9a2301e7bb61c0b3752ed)) - (Lin Yang)
- Improve Phase 3 protocol module plan ([add24e1](add24e142a5d9b6082abf0f4d742a49df726bfd6)) - (Lin Yang)
- Improve Phase 4 transport module plan ([3326f1d](3326f1de79febeb3c22ab7fa029545e3bbc81d36)) - (Lin Yang)
- Improve Phase 5 connection module plan ([16f0ff9](16f0ff959dbdcc71d401de0b230e3b0f53e2f865)) - (Lin Yang)
- Improve Phase 6 events module plan ([808a74b](808a74b4190935e0d1aa0a7a17b6f6bc05e050d3)) - (Lin Yang)
- Improve Phase 7 managers module plan ([18fc763](18fc763c141ca07ddc71c6e5f5d31ce29bec8d6a)) - (Lin Yang)
- Improve Phase 7 managers module plan ([9251be1](9251be18d3fde50090790b4091123e4696a47648)) - (Lin Yang)
- Add timeout validation to Phase 8 utils module plan ([f1d2abe](f1d2abebb17c045f321ea6133f312522f8216c05)) - (Lin Yang)
- Implement Phase 9 main client with thread-safe operations ([6deda27](6deda27b1fd3fd85b021426f4719c332e742058e)) - (Lin Yang)
- Fix import paths and cross-phase API compatibility ([cddfd4d](cddfd4d7231dab75257048f4d9da809e77588107)) - (Lin Yang)
- Update all phase plans to use pkg/openclaw directory structure ([8dcac5e](8dcac5e9c57433d5116db8d4b1f00ffc5341c978)) - (Lin Yang)
- Add CLAUDE.md with project overview and development commands ([9f3d6d5](9f3d6d5c0a89cb7be32765b8b3bdd7ec65bb6efc)) - (Lin Yang)
- Add comprehensive README with installation, usage, and API documentation ([f205051](f205051a193bcb5bc9df0bccbebeccc4764c939a)) - (Lin Yang)
- Update README with badges and project description ([4a61d69](4a61d69a55b38a5637fe4b5a4919551417b055e2)) - (Lin Yang)
- Add comprehensive documentation comments to all non-test Go files ([253f510](253f5104c8523e23a4e79dc8eac9066586c63e07)) - (Lin Yang)
- Document dual ErrorShape design intent ([8c8543a](8c8543a4d5b131e072bd0ae63df4bc6f5b353af4)) - (Lin Yang)
- Improve CLAUDE.md with Key Files and Gotchas sections ([6d858d8](6d858d874ff220991a9309a7548e6085542b7a6f)) - (Lin Yang)
- Clarify CRL/OCSP stub with actionable TODO ([66620a7](66620a703108ee829ced44c4086b3b6de5297d07)) - (Lin Yang)
- Add Codecov coverage badge ([803e7ef](803e7ef60d56799c0e52b70b55299355b7247a0e)) - (Lin Yang)
- Update with new features and API additions since v1.0 ([8fb6042](8fb604231bb29f7bc87e19716a0b8d82c2dbefa7)) - (Lin Yang)

### 🔧 Miscellaneous Tasks

- Add GitHub Actions workflows ([88f31fe](88f31fef50fca21765641540e86739d7fe6cb345)) - (Lin Yang)
- Modify CI workflow for Go version and fail-fast ([1abdb9e](1abdb9ea1518365853e58e64e3caeba2e5cf4bee)) - (Lin Yang)



<!-- generated by git-cliff -->
