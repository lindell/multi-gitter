# Changelog

### [0.20.4](https://www.github.com/lindell/multi-gitter/compare/v0.20.3...v0.20.4) (2021-02-19)


### Miscellaneous

* removed brew download strategy ([6c35be5](https://www.github.com/lindell/multi-gitter/commit/6c35be50c14dff82363931bd12ade4de204103c1))

### [0.20.3](https://www.github.com/lindell/multi-gitter/compare/v0.20.2...v0.20.3) (2021-02-19)


### Bug Fixes

* fixed homebrew release ([dece0d8](https://www.github.com/lindell/multi-gitter/commit/dece0d8ad5e40c20be37eeb9db42dcdfd9eaf4d4))

### [0.20.2](https://www.github.com/lindell/multi-gitter/compare/v0.20.1...v0.20.2) (2021-02-19)


### Miscellaneous

* added brew install ([#73](https://www.github.com/lindell/multi-gitter/issues/73)) ([3f56a4a](https://www.github.com/lindell/multi-gitter/commit/3f56a4aefe6a984b781e2e5792b929ea6b6962e3))
* improved the base-url description ([5e7ec24](https://www.github.com/lindell/multi-gitter/commit/5e7ec248b5dee732fc276d7c40a62a8c4fb76c1c))

### [0.20.1](https://www.github.com/lindell/multi-gitter/compare/v0.20.0...v0.20.1) (2021-02-17)


### Miscellaneous

* updated to go 1.16 ([c8fa961](https://www.github.com/lindell/multi-gitter/commit/c8fa96154f2925ea724bf9ff0a74027dfc0a9286))

## [0.20.0](https://www.github.com/lindell/multi-gitter/compare/v0.19.1...v0.20.0) (2021-02-16)


### Features

* **gitlab:** option to change base-url for gitlab ([#69](https://www.github.com/lindell/multi-gitter/issues/69)) ([147ebe6](https://www.github.com/lindell/multi-gitter/commit/147ebe67d2902850f06c7575bbe8e43b0372eccd))


### Dependencies

* update module spf13/cobra to v1.1.3 ([7a32bb6](https://www.github.com/lindell/multi-gitter/commit/7a32bb615e0969aa41618885023d797ea101cf5b))
* update module xanzy/go-gitlab to v0.44.0 ([53f834b](https://www.github.com/lindell/multi-gitter/commit/53f834b29f9e801a9fa5d8416ad18b22635ab058))

### [0.19.1](https://www.github.com/lindell/multi-gitter/compare/v0.19.0...v0.19.1) (2021-02-02)


### Dependencies

* update golang.org/x/oauth2 commit hash to f9ce19e ([#66](https://www.github.com/lindell/multi-gitter/issues/66)) ([64d9095](https://www.github.com/lindell/multi-gitter/commit/64d90952856fdfd0517cf03bb752603c708ff6b9))
* update module xanzy/go-gitlab to v0.43.0 ([1a44511](https://www.github.com/lindell/multi-gitter/commit/1a44511a7cb27aefb3e0e6d8e4309e3fc78f4756))

## [0.19.0](https://www.github.com/lindell/multi-gitter/compare/v0.18.0...v0.19.0) (2021-01-21)


### Features

* added --merge-type flag ([#64](https://www.github.com/lindell/multi-gitter/issues/64)) ([dd18402](https://www.github.com/lindell/multi-gitter/commit/dd18402365c0f41440bd580497cbd12e0738bc7e))

## [0.18.0](https://www.github.com/lindell/multi-gitter/compare/v0.17.0...v0.18.0) (2021-01-20)


### Features

* added --fetch-depth flag ([#62](https://www.github.com/lindell/multi-gitter/issues/62)) ([5cdb723](https://www.github.com/lindell/multi-gitter/commit/5cdb72334f151c4950ffd9763b8ee760dbc3f8a5))

## [0.17.0](https://www.github.com/lindell/multi-gitter/compare/v0.16.4...v0.17.0) (2021-01-20)


### Features

* added links to printed prs ([#58](https://www.github.com/lindell/multi-gitter/issues/58)) ([cd76c61](https://www.github.com/lindell/multi-gitter/commit/cd76c6143a9b008f6be08748b77f7c8acc36aaf9))


### Bug Fixes

* added the number of created pull requests ([#56](https://www.github.com/lindell/multi-gitter/issues/56)) ([d432430](https://www.github.com/lindell/multi-gitter/commit/d4324307441ffc74002e1cb4f5c08b83f45a2781))

### [0.16.4](https://www.github.com/lindell/multi-gitter/compare/v0.16.3...v0.16.4) (2021-01-16)


### Bug Fixes

* multi-gitter does now only fetch the base and head branch ([b272644](https://www.github.com/lindell/multi-gitter/commit/b272644355d9291c23de8f028a3132de5a5eb99e))


### Dependencies

* update module stretchr/testify to v1.7.0 ([0a06a24](https://www.github.com/lindell/multi-gitter/commit/0a06a247d93d34986504608d1ff437aa17869d53))

### [0.16.3](https://www.github.com/lindell/multi-gitter/compare/v0.16.2...v0.16.3) (2021-01-15)


### Bug Fixes

* fixed presentation of repos with existing repo ([ac8027b](https://www.github.com/lindell/multi-gitter/commit/ac8027b3cf6c8df46ae3c4e2b79891c14962f7bc))

### [0.16.2](https://www.github.com/lindell/multi-gitter/compare/v0.16.1...v0.16.2) (2021-01-14)


### Bug Fixes

* fixed bug where base branch was left empty ([64d5e22](https://www.github.com/lindell/multi-gitter/commit/64d5e225e631f8b3a0dac3fc3145f0168dacba70))


### Dependencies

* update golang.org/x/oauth2 commit hash to d3ed898 ([feea168](https://www.github.com/lindell/multi-gitter/commit/feea168f7a2d44d9fe08c8b1a995dfc5b213f7ce))

### [0.16.1](https://www.github.com/lindell/multi-gitter/compare/v0.16.0...v0.16.1) (2021-01-12)


### Dependencies

* update module xanzy/go-gitlab to v0.41.0 ([f713ee2](https://www.github.com/lindell/multi-gitter/commit/f713ee227b1013d3e2f293fa4d50dbdbbf980b17))
* update module xanzy/go-gitlab to v0.42.0 ([fd8e373](https://www.github.com/lindell/multi-gitter/commit/fd8e3737e9db7348cf271adab5d6a958e1a794f4))

## [0.16.0](https://www.github.com/lindell/multi-gitter/compare/v0.15.3...v0.16.0) (2020-12-18)


### Features

* added base-branch flag ([8c04b8d](https://www.github.com/lindell/multi-gitter/commit/8c04b8d241f66ec8def92baf8ae27a39a24abcff))


### Dependencies

* update module xanzy/go-gitlab to v0.40.2 ([ce33ff5](https://www.github.com/lindell/multi-gitter/commit/ce33ff5e69a9d984199c0f00e0b8c57ef6bfbc93))

### [0.15.3](https://www.github.com/lindell/multi-gitter/compare/v0.15.2...v0.15.3) (2020-12-08)


### Dependencies

* update golang.org/x/oauth2 commit hash to 08078c5 ([8b94c50](https://www.github.com/lindell/multi-gitter/commit/8b94c50acf0df6c5af2a9b4e81a7aea296575bc8))
* update golang.org/x/oauth2 commit hash to 9317641 ([242cdd0](https://www.github.com/lindell/multi-gitter/commit/242cdd09b5c62813336059f8ac73cf662bc4f71e))

### [0.15.2](https://www.github.com/lindell/multi-gitter/compare/v0.15.1...v0.15.2) (2020-12-03)


### Miscellaneous

* bump github.com/xanzy/go-gitlab from 0.40.0 to 0.40.1 ([550f302](https://www.github.com/lindell/multi-gitter/commit/550f302b23ac1301e84994846766c95d2011fb6e))


### Dependencies

* update golang.org/x/oauth2 commit hash to 0b49973 ([ee02a57](https://www.github.com/lindell/multi-gitter/commit/ee02a57a4512e9e8b29770ec610974ad5ecdf7d2))

### [0.15.1](https://www.github.com/lindell/multi-gitter/compare/v0.15.0...v0.15.1) (2020-12-02)


### Bug Fixes

* corrected the name of the REPOSITORY env var ([9b87070](https://www.github.com/lindell/multi-gitter/commit/9b8707096a85d1106045fb79d13b87c5fe8b99de))


### Miscellaneous

* bump github.com/xanzy/go-gitlab from 0.39.0 to 0.40.0 ([e65e1a8](https://www.github.com/lindell/multi-gitter/commit/e65e1a8a480202e81a3b2488d0c0d350fa3f265d))
