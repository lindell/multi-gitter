# Changelog

### [0.40.1](https://www.github.com/lindell/multi-gitter/compare/v0.40.0...v0.40.1) (2022-03-01)


### Bug Fixes

* **gitlab:** made entire owner structure a part of logging and $REPOSITORY ([ccac98a](https://www.github.com/lindell/multi-gitter/commit/ccac98a8c6ab7e4dc23e19519de7ccd958e1be43))


### Dependencies

* update module github.com/spf13/viper to v1.10.1 ([6052e73](https://www.github.com/lindell/multi-gitter/commit/6052e73c7e2f1a790fa4c23732aec0c23dcd131c))
* update module github.com/xanzy/go-gitlab to v0.55.1 ([ea0eacc](https://www.github.com/lindell/multi-gitter/commit/ea0eacce1c3a711b9f3c0d51c2b714ebc72ddb23))


### Miscellaneous

* fixed gitlab api pointer change ([5c39b22](https://www.github.com/lindell/multi-gitter/commit/5c39b2208cb3480eaab816a18d029bd9b20b7ddc))

## [0.40.0](https://www.github.com/lindell/multi-gitter/compare/v0.39.0...v0.40.0) (2022-02-08)


### Features

* add ability to create PR as draft ([#232](https://www.github.com/lindell/multi-gitter/issues/232)) ([dbfef2b](https://www.github.com/lindell/multi-gitter/commit/dbfef2b0f8c3692d5f281d2269bc680263ec2406))

## [0.39.0](https://www.github.com/lindell/multi-gitter/compare/v0.38.3...v0.39.0) (2022-02-05)


### Features

* **gitlab:** respect project level squash setting ([#228](https://www.github.com/lindell/multi-gitter/issues/228)) ([b189661](https://www.github.com/lindell/multi-gitter/commit/b1896610c6f4f88099848f99163984fbf4de113f))


### Bug Fixes

* **gitlab:** close MR instead of deleting it ([#230](https://www.github.com/lindell/multi-gitter/issues/230)) ([af2c2d9](https://www.github.com/lindell/multi-gitter/commit/af2c2d9b86a8cd0c6c09e3667a9f7c7689cf915c))
* better logs when repositories are not used due to permissions ([#226](https://www.github.com/lindell/multi-gitter/issues/226)) ([668d0b0](https://www.github.com/lindell/multi-gitter/commit/668d0b0458988dc24de7beb4d087c2cccc7d3167))


### Dependencies

* update module code.gitea.io/sdk/gitea to v0.15.1 ([c8b4ab8](https://www.github.com/lindell/multi-gitter/commit/c8b4ab8cf55ace9e4f4751b540ba2e0dc7fd2807))
* update module github.com/spf13/cobra to v1.3.0 ([afbe211](https://www.github.com/lindell/multi-gitter/commit/afbe21104568c9550e7cbaa48342ed5bcf3de3b8))

### [0.38.3](https://www.github.com/lindell/multi-gitter/compare/v0.38.2...v0.38.3) (2022-01-21)


### Bug Fixes

* fixed multi line commit message ([#222](https://www.github.com/lindell/multi-gitter/issues/222)) ([995a93c](https://www.github.com/lindell/multi-gitter/commit/995a93cf552300e3c4b580a0ffa2aeb3cdfd61a7))


### Dependencies

* update module github.com/mitchellh/mapstructure to v1.4.3 ([ea97fa2](https://www.github.com/lindell/multi-gitter/commit/ea97fa2010378b353c0a43f832aced2054ae182d))
* update module github.com/xanzy/go-gitlab to v0.52.2 ([c88f791](https://www.github.com/lindell/multi-gitter/commit/c88f79152e1b99381fcc8de249149aa49089c77e))

### [0.38.2](https://www.github.com/lindell/multi-gitter/compare/v0.38.1...v0.38.2) (2021-12-27)


### Bug Fixes

* fixed fury.io token in release ([9deb2d2](https://www.github.com/lindell/multi-gitter/commit/9deb2d258a6b84c8f002b4aaab0ffc8ca135540a))

### [0.38.1](https://www.github.com/lindell/multi-gitter/compare/v0.38.0...v0.38.1) (2021-12-27)


### Bug Fixes

* make sure autocompletion works with settings from config files ([#217](https://www.github.com/lindell/multi-gitter/issues/217)) ([45e855c](https://www.github.com/lindell/multi-gitter/commit/45e855c314738f710092f2babb6a213727a63467))

## [0.38.0](https://www.github.com/lindell/multi-gitter/compare/v0.37.0...v0.38.0) (2021-12-16)


### Features

* added --ssh-auth option ([#215](https://www.github.com/lindell/multi-gitter/issues/215)) ([f5767a8](https://www.github.com/lindell/multi-gitter/commit/f5767a86c44562f3191eb8cff0d3084393ed1ac7))

## [0.37.0](https://www.github.com/lindell/multi-gitter/compare/v0.36.1...v0.37.0) (2021-11-08)


### Features

* added --conflict-strategy ([#210](https://www.github.com/lindell/multi-gitter/issues/210)) ([5dfd6d9](https://www.github.com/lindell/multi-gitter/commit/5dfd6d9fc877d06f905cbeb27e39305d16afee65))


### Bug Fixes

* **bitbucket:** ensure username is set ([#212](https://www.github.com/lindell/multi-gitter/issues/212)) ([a463709](https://www.github.com/lindell/multi-gitter/commit/a4637093e71c3b667afb58e67439bb2b3c9fe927))


### Dependencies

* update module code.gitea.io/sdk/gitea to v0.15.0 ([1b0ac09](https://www.github.com/lindell/multi-gitter/commit/1b0ac094015c4b398147c3fc8759e83462b656b7))
* update module github.com/google/go-github/v39 to v39.2.0 ([45f20a0](https://www.github.com/lindell/multi-gitter/commit/45f20a0e070e71231d785d1bb12cc04ee0d0e2e2))
* update module github.com/xanzy/go-gitlab to v0.51.1 ([78fb3dc](https://www.github.com/lindell/multi-gitter/commit/78fb3dca3b3d5e1aff66799814e3c4a92edda0d7))

### [0.36.1](https://www.github.com/lindell/multi-gitter/compare/v0.36.0...v0.36.1) (2021-10-28)


### Bug Fixes

* make sure GitHub's secondary rate limit is not reached ([#207](https://www.github.com/lindell/multi-gitter/issues/207)) ([8a5fabd](https://www.github.com/lindell/multi-gitter/commit/8a5fabdc9e54bdfbba421a466bd323aae8114bdd))

## [0.36.0](https://www.github.com/lindell/multi-gitter/compare/v0.35.0...v0.36.0) (2021-10-24)


### Features

* add the ability to skip repos from the run command ([#197](https://www.github.com/lindell/multi-gitter/issues/197)) ([d4de4dc](https://www.github.com/lindell/multi-gitter/commit/d4de4dc5dc0d05726db2dabeb515303c21d53994))

## [0.35.0](https://www.github.com/lindell/multi-gitter/compare/v0.34.0...v0.35.0) (2021-10-20)


### Features

* add possibility to add assignees to pull request ([#196](https://www.github.com/lindell/multi-gitter/issues/196)) ([6b685ba](https://www.github.com/lindell/multi-gitter/commit/6b685ba18ce7107e92984fd9654c9c1af274bf95))

## [0.34.0](https://www.github.com/lindell/multi-gitter/compare/v0.33.3...v0.34.0) (2021-10-17)


### Features

* added more information to the version command ([#198](https://www.github.com/lindell/multi-gitter/issues/198)) ([ebf4578](https://www.github.com/lindell/multi-gitter/commit/ebf457822693000fe04caf4c36a5db70c9feab6c))

### [0.33.3](https://www.github.com/lindell/multi-gitter/compare/v0.33.2...v0.33.3) (2021-10-11)


### Bug Fixes

* **github:** fixed that the fetching of pullrequests always pull the latest pr ([#195](https://www.github.com/lindell/multi-gitter/issues/195)) ([aa33af8](https://www.github.com/lindell/multi-gitter/commit/aa33af834d71e7122955b5023ab028c2d5fa42f8))


### Dependencies

* update github.com/gfleury/go-bitbucket-v1 commit hash to dff2223 ([f570ee5](https://www.github.com/lindell/multi-gitter/commit/f570ee5086369b50091f5ad21f6762d96d93782c))
* update module github.com/google/go-github/v38 to v39 ([#191](https://www.github.com/lindell/multi-gitter/issues/191)) ([5088532](https://www.github.com/lindell/multi-gitter/commit/508853232485cd4dd4886f46fead14fa71d7ae59))
* update module github.com/spf13/viper to v1.9.0 ([becff1f](https://www.github.com/lindell/multi-gitter/commit/becff1f7d6fd755565601fd6eb4d321cac6d54a2))

### [0.33.2](https://www.github.com/lindell/multi-gitter/compare/v0.33.1...v0.33.2) (2021-09-29)


### Bug Fixes

* **github:** ignore branch deletion error if branch is already deleted ([#189](https://www.github.com/lindell/multi-gitter/issues/189)) ([d63d041](https://www.github.com/lindell/multi-gitter/commit/d63d04184dc10d3c6538676dacdd63d973d06e02))
* censor http authentication header ([#185](https://www.github.com/lindell/multi-gitter/issues/185)) ([633a2cc](https://www.github.com/lindell/multi-gitter/commit/633a2ccc973070790b0cb644aa9029727a220e20))

### [0.33.1](https://www.github.com/lindell/multi-gitter/compare/v0.33.0...v0.33.1) (2021-09-23)


### Bug Fixes

* **gitlab:** only list projects with Merge Requests enabled ([#184](https://www.github.com/lindell/multi-gitter/issues/184)) ([5d45121](https://www.github.com/lindell/multi-gitter/commit/5d4512112715dbe9ce7cba214531ce93c8b1a360))


### Miscellaneous

* added CODEOWNERS file ([7b85777](https://www.github.com/lindell/multi-gitter/commit/7b8577798fbcc3159dfc06d920ee5b33183f0ce9))

## [0.33.0](https://www.github.com/lindell/multi-gitter/compare/v0.32.0...v0.33.0) (2021-09-10)


### Features

* **bitbucketserver:** added support for bitbucket server ([#178](https://www.github.com/lindell/multi-gitter/issues/178)) ([2f7a1b6](https://www.github.com/lindell/multi-gitter/commit/2f7a1b6e313355a8aa4176cc216bd2d9ad6494a7))


### Dependencies

* update golang.org/x/oauth2 commit hash to 2bc19b1 ([858441a](https://www.github.com/lindell/multi-gitter/commit/858441a9822b6f86d9e68216742f550eb80f7e05))
* update module github.com/google/go-github/v37 to v38 ([#176](https://www.github.com/lindell/multi-gitter/issues/176)) ([f15aaad](https://www.github.com/lindell/multi-gitter/commit/f15aaad21ba92a4d2c05c039f0b7f8963f245e75))

## [0.32.0](https://www.github.com/lindell/multi-gitter/compare/v0.31.1...v0.32.0) (2021-08-12)


### Features

* added --config to status command ([#174](https://www.github.com/lindell/multi-gitter/issues/174)) ([8c52c93](https://www.github.com/lindell/multi-gitter/commit/8c52c931df5fe786a9b9c26e77aebe50241f8391))

### [0.31.1](https://www.github.com/lindell/multi-gitter/compare/v0.31.0...v0.31.1) (2021-08-12)


### Bug Fixes

* added support for GitLab subgroups with --project ([#171](https://www.github.com/lindell/multi-gitter/issues/171)) ([25b5d54](https://www.github.com/lindell/multi-gitter/commit/25b5d543056909fdb1a937118989f06dd4902f80))

## [0.31.0](https://www.github.com/lindell/multi-gitter/compare/v0.30.0...v0.31.0) (2021-08-08)


### Features

* interactive mode  ([#167](https://www.github.com/lindell/multi-gitter/issues/167)) ([7351520](https://www.github.com/lindell/multi-gitter/commit/73515206bc7201b28e0e1faef7e1009b3e5a34f9))

## [0.30.0](https://www.github.com/lindell/multi-gitter/compare/v0.29.2...v0.30.0) (2021-08-01)


### Features

* moved to built in completion command in cobra 1.2.x ([#163](https://www.github.com/lindell/multi-gitter/issues/163)) ([81a7187](https://www.github.com/lindell/multi-gitter/commit/81a7187fce1ab76e6d87bdeee02b268fdb21320b))


### Dependencies

* update module github.com/google/go-github/v36 to v37 ([213a1c6](https://www.github.com/lindell/multi-gitter/commit/213a1c6cc603cec49f889ffe52dd50d22f33ab44))
* update module github.com/xanzy/go-gitlab to v0.50.1 ([cac5518](https://www.github.com/lindell/multi-gitter/commit/cac5518094e4cc82bf2ad6d47c42a593f3031034))

### [0.29.2](https://www.github.com/lindell/multi-gitter/compare/v0.29.1...v0.29.2) (2021-07-01)


### Bug Fixes

* push hooks no longer run with cmd-git implementation ([#159](https://www.github.com/lindell/multi-gitter/issues/159)) ([7360c0d](https://www.github.com/lindell/multi-gitter/commit/7360c0d14b83be627325d0b4ea95177e71c2a565))


### Dependencies

* update golang.org/x/oauth2 commit hash to a41e5a7 ([234ce36](https://www.github.com/lindell/multi-gitter/commit/234ce36753e5eec8d73700a4b65e4ee8ad0773a7))
* update module github.com/go-git/go-git/v5 to v5.4.2 ([016f54d](https://www.github.com/lindell/multi-gitter/commit/016f54d39a8df80558b7c46880a7dfabd16c7e28))
* update module github.com/google/go-github/v35 to v36 ([#162](https://www.github.com/lindell/multi-gitter/issues/162)) ([893d9ea](https://www.github.com/lindell/multi-gitter/commit/893d9eae5dd5f8abcf6c00cb233957aea532d1c2))

### [0.29.1](https://www.github.com/lindell/multi-gitter/compare/v0.29.0...v0.29.1) (2021-06-28)


### Bug Fixes

* commit hooks no longer run with cmd-git implementation ([#157](https://www.github.com/lindell/multi-gitter/issues/157)) ([ba12d08](https://www.github.com/lindell/multi-gitter/commit/ba12d08fee2e8cc0ef8015a1761afde747a2622c))
* downgraded go-diff to fix diff formating ([#156](https://www.github.com/lindell/multi-gitter/issues/156)) ([6ef43a8](https://www.github.com/lindell/multi-gitter/commit/6ef43a847f14d5b81745e9978732eebda5bf8ca9))

## [0.29.0](https://www.github.com/lindell/multi-gitter/compare/v0.28.0...v0.29.0) (2021-06-20)


### Features

* added configuration options through config files ([#150](https://www.github.com/lindell/multi-gitter/issues/150)) ([f38a7ad](https://www.github.com/lindell/multi-gitter/commit/f38a7ad3ffc9f6aaef60913a6a08006b5b672a93))


### Bug Fixes

* made sure any tokens output in the logs are now censored ([#143](https://www.github.com/lindell/multi-gitter/issues/143)) ([0e5cee7](https://www.github.com/lindell/multi-gitter/commit/0e5cee7ecd6dde23d21869058cc383e83b232703))

## [0.28.0](https://www.github.com/lindell/multi-gitter/compare/v0.27.0...v0.28.0) (2021-06-16)


### Features

* added --git-type flag ([cb4701e](https://www.github.com/lindell/multi-gitter/commit/cb4701eb90b98bf585b1a8835368c4cd8f0e0095))

## [0.27.0](https://www.github.com/lindell/multi-gitter/compare/v0.26.1...v0.27.0) (2021-06-14)


### Features

* added fork mode ([#128](https://www.github.com/lindell/multi-gitter/issues/128)) ([f9e7827](https://www.github.com/lindell/multi-gitter/commit/f9e78273440642be662686912b89ff38123bacf7))


### Miscellaneous

* improved logging and added stack trace if --log-level=trace is used ([#138](https://www.github.com/lindell/multi-gitter/issues/138)) ([abccc5f](https://www.github.com/lindell/multi-gitter/commit/abccc5f28ba22e3b2b99d6d3ee1c513a213caea7))

### [0.26.1](https://www.github.com/lindell/multi-gitter/compare/v0.26.0...v0.26.1) (2021-06-09)


### Bug Fixes

* made remove branch on merge the default behaviour for GitLab merge ([#135](https://www.github.com/lindell/multi-gitter/issues/135)) ([9cc5983](https://www.github.com/lindell/multi-gitter/commit/9cc5983407c3b5be4a42c55dbd7c4b03f54d3f23))

## [0.26.0](https://www.github.com/lindell/multi-gitter/compare/v0.25.6...v0.26.0) (2021-06-08)


### Features

* added --include-subgroups flag ([#131](https://www.github.com/lindell/multi-gitter/issues/131)) ([eff19a4](https://www.github.com/lindell/multi-gitter/commit/eff19a4b23030487fa9a3e64553443d2a8fb3133))


### Bug Fixes

* improved error messages for common problems with the script ([de9e525](https://www.github.com/lindell/multi-gitter/commit/de9e5259d2bd900abf72c56f40a76f223cbfffd0))

### [0.25.6](https://www.github.com/lindell/multi-gitter/compare/v0.25.5...v0.25.6) (2021-06-05)


### Bug Fixes

* fixed skip-pr flag description ([#127](https://www.github.com/lindell/multi-gitter/issues/127)) ([1c4e2ac](https://www.github.com/lindell/multi-gitter/commit/1c4e2acec3fee563eb3cfa7391f63ffd5fc1d61e))
* typo where archived should be achieved ([#125](https://www.github.com/lindell/multi-gitter/issues/125)) ([5373ea8](https://www.github.com/lindell/multi-gitter/commit/5373ea8fd37e39ce1eb8edbb860af85faa47e370))

### [0.25.5](https://www.github.com/lindell/multi-gitter/compare/v0.25.4...v0.25.5) (2021-06-01)


### Dependencies

* update golang.org/x/oauth2 commit hash to f6687ab ([cab768a](https://www.github.com/lindell/multi-gitter/commit/cab768a1a6bf93b8f113b0b7221db7a4bab375cd))
* update module github.com/go-git/go-git/v5 to v5.4.1 ([fe45f2e](https://www.github.com/lindell/multi-gitter/commit/fe45f2e9ad2031ae4f436271e4a072101ba80805))

### [0.25.4](https://www.github.com/lindell/multi-gitter/compare/v0.25.3...v0.25.4) (2021-05-16)


### Bug Fixes

* make sure gitignore is used ([#119](https://www.github.com/lindell/multi-gitter/issues/119)) ([f33dee9](https://www.github.com/lindell/multi-gitter/commit/f33dee9a7acd798ab6ad0a7255351c50c9bd456e))

### [0.25.3](https://www.github.com/lindell/multi-gitter/compare/v0.25.2...v0.25.3) (2021-05-11)


### Bug Fixes

* added panic recover on a run repo basis ([#114](https://www.github.com/lindell/multi-gitter/issues/114)) ([6d44adf](https://www.github.com/lindell/multi-gitter/commit/6d44adf5ddbf3783bc4a2224c35a923ab599e7c6))


### Dependencies

* update module github.com/sergi/go-diff to v1.2.0 ([#116](https://www.github.com/lindell/multi-gitter/issues/116)) ([0273abe](https://www.github.com/lindell/multi-gitter/commit/0273abeba104e2ce522b4a97b1498341ad9e41d6))

### [0.25.2](https://www.github.com/lindell/multi-gitter/compare/v0.25.1...v0.25.2) (2021-05-11)


### Bug Fixes

* skip running git diff if debug or lower is not set ([#113](https://www.github.com/lindell/multi-gitter/issues/113)) ([5189374](https://www.github.com/lindell/multi-gitter/commit/51893745153e7825339f7398e844bf6d53404cc8))


### Dependencies

* update module github.com/google/go-github/v33 to v35 ([#110](https://www.github.com/lindell/multi-gitter/issues/110)) ([b6c8667](https://www.github.com/lindell/multi-gitter/commit/b6c8667f1ca48c62b1ec1703f8afa1664dfeca95))

### [0.25.1](https://www.github.com/lindell/multi-gitter/compare/v0.25.0...v0.25.1) (2021-05-01)


### Dependencies

* update golang.org/x/oauth2 commit hash to 81ed05c ([#107](https://www.github.com/lindell/multi-gitter/issues/107)) ([b529c3f](https://www.github.com/lindell/multi-gitter/commit/b529c3f113ccda92ee0981c48b0c26c74facb142))
* update module github.com/go-git/go-git/v5 to v5.3.0 ([905dbdb](https://www.github.com/lindell/multi-gitter/commit/905dbdbfa5b420ee985bed2ff58cfb2399b051b7))
* update module github.com/xanzy/go-gitlab to v0.49.0 ([#109](https://www.github.com/lindell/multi-gitter/issues/109)) ([597d8b4](https://www.github.com/lindell/multi-gitter/commit/597d8b41751b0cf90bc4743fc367ed72487ae35f))

## [0.25.0](https://www.github.com/lindell/multi-gitter/compare/v0.24.2...v0.25.0) (2021-04-25)


### Features

* added Gitea support ([#105](https://www.github.com/lindell/multi-gitter/issues/105)) ([0f89791](https://www.github.com/lindell/multi-gitter/commit/0f89791d62fe32f0d2a98f0b735782898976e3f7))

### [0.24.2](https://www.github.com/lindell/multi-gitter/compare/v0.24.1...v0.24.2) (2021-04-01)


### Dependencies

* update golang.org/x/oauth2 commit hash to 22b0ada ([#92](https://www.github.com/lindell/multi-gitter/issues/92)) ([335eee3](https://www.github.com/lindell/multi-gitter/commit/335eee37c02c54fa7d006ff0aab837b424f7d514))
* update module github.com/google/go-github/v33 to v34 ([#93](https://www.github.com/lindell/multi-gitter/issues/93)) ([03d3278](https://www.github.com/lindell/multi-gitter/commit/03d327835bb7a99362e0b13224200a41d068a642))

### [0.24.1](https://www.github.com/lindell/multi-gitter/compare/v0.24.0...v0.24.1) (2021-03-31)


### Bug Fixes

* fixed windows filepaths ([#89](https://www.github.com/lindell/multi-gitter/issues/89)) ([cb38fc0](https://www.github.com/lindell/multi-gitter/commit/cb38fc08a084dd7b5b05717b852e9804d52e1720))

## [0.24.0](https://www.github.com/lindell/multi-gitter/compare/v0.23.1...v0.24.0) (2021-03-30)


### Features

* added static flag completion for enums ([#87](https://www.github.com/lindell/multi-gitter/issues/87)) ([586dd61](https://www.github.com/lindell/multi-gitter/commit/586dd616418affe1838b4ecfb5714458ffcafd0b))

### [0.23.1](https://www.github.com/lindell/multi-gitter/compare/v0.23.0...v0.23.1) (2021-03-30)


### Bug Fixes

* fixed brew test command ([fc243e8](https://www.github.com/lindell/multi-gitter/commit/fc243e8d7d94c9b1793eb7299a893ba2ba14794c))

## [0.23.0](https://www.github.com/lindell/multi-gitter/compare/v0.22.1...v0.23.0) (2021-03-30)


### Features

* added GitHub autocompletion ([#84](https://www.github.com/lindell/multi-gitter/issues/84)) ([5fee0c4](https://www.github.com/lindell/multi-gitter/commit/5fee0c4b88e802a8be4168f802b79a1701afd3a6))

### [0.22.1](https://www.github.com/lindell/multi-gitter/compare/v0.22.0...v0.22.1) (2021-03-12)


### Dependencies

* update module github.com/google/go-github/v32 to v33 ([#82](https://www.github.com/lindell/multi-gitter/issues/82)) ([1c48de3](https://www.github.com/lindell/multi-gitter/commit/1c48de3a81a64cbac6481b3260bdc3512e98a34f))
* update module github.com/sirupsen/logrus to v1.8.1 ([31dad70](https://www.github.com/lindell/multi-gitter/commit/31dad70383ab5c6d742e393ad97c32c610f85c2b))
* update module github.com/xanzy/go-gitlab to v0.46.0 ([f0a3503](https://www.github.com/lindell/multi-gitter/commit/f0a350323dcee32acf20265b5ebcedce2f1531b9))
* update module github.com/xanzy/go-gitlab to v0.47.0 ([92a18a3](https://www.github.com/lindell/multi-gitter/commit/92a18a3fd27136b5a46c2de423c109d28ad9da71))

## [0.22.0](https://www.github.com/lindell/multi-gitter/compare/v0.21.1...v0.22.0) (2021-03-03)


### Features

* added skip-pr flag ([#80](https://www.github.com/lindell/multi-gitter/issues/80)) ([c4b85ea](https://www.github.com/lindell/multi-gitter/commit/c4b85ea5606a361b13b0a6308f3cfea776f954ad))


### Dependencies

* update module github.com/sirupsen/logrus to v1.8.0 ([8c132b4](https://www.github.com/lindell/multi-gitter/commit/8c132b410baef2812a5525727147eb11f939a870))
* update module github.com/xanzy/go-gitlab to v0.45.0 ([9e1bc9a](https://www.github.com/lindell/multi-gitter/commit/9e1bc9a163f151fddba8f24f44934111d7e9b810))

### [0.21.1](https://www.github.com/lindell/multi-gitter/compare/v0.21.0...v0.21.1) (2021-02-19)


### Bug Fixes

* fixed license file in release ([506084f](https://www.github.com/lindell/multi-gitter/commit/506084fd8b17f42a3311524bc0dbcc29ce39c50b))

## [0.21.0](https://www.github.com/lindell/multi-gitter/compare/v0.20.5...v0.21.0) (2021-02-19)


### Features

* added shell completion command ([c5782a2](https://www.github.com/lindell/multi-gitter/commit/c5782a2e377ecfc071c82c1db0a775e45215a0cc))

### [0.20.5](https://www.github.com/lindell/multi-gitter/compare/v0.20.4...v0.20.5) (2021-02-19)


### Miscellaneous

* moved brew release to Formula folder ([d1ae864](https://www.github.com/lindell/multi-gitter/commit/d1ae8644cd2a4e6138b3f23d23e792509ff7b3ef))

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
