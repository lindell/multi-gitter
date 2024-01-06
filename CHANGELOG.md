# Changelog

## [1.0.0](https://github.com/lindell/multi-gitter/compare/v0.49.1...v1.0.0) (2024-01-06)


### ⚠ BREAKING CHANGES

* added output option to differentiate log and other output

### Features

* add ability to create PR as draft ([#232](https://github.com/lindell/multi-gitter/issues/232)) ([dbfef2b](https://github.com/lindell/multi-gitter/commit/dbfef2b0f8c3692d5f281d2269bc680263ec2406))
* add possibility to add assignees to pull request ([#196](https://github.com/lindell/multi-gitter/issues/196)) ([6b685ba](https://github.com/lindell/multi-gitter/commit/6b685ba18ce7107e92984fd9654c9c1af274bf95))
* add repo-include and repo-exclude options ([#426](https://github.com/lindell/multi-gitter/issues/426)) ([faf0092](https://github.com/lindell/multi-gitter/commit/faf0092a4c5411420f721973595588f6303253da))
* add the ability to skip repos from the run command ([#197](https://github.com/lindell/multi-gitter/issues/197)) ([d4de4dc](https://github.com/lindell/multi-gitter/commit/d4de4dc5dc0d05726db2dabeb515303c21d53994))
* added --config to status command ([#174](https://github.com/lindell/multi-gitter/issues/174)) ([8c52c93](https://github.com/lindell/multi-gitter/commit/8c52c931df5fe786a9b9c26e77aebe50241f8391))
* added --conflict-strategy ([#210](https://github.com/lindell/multi-gitter/issues/210)) ([5dfd6d9](https://github.com/lindell/multi-gitter/commit/5dfd6d9fc877d06f905cbeb27e39305d16afee65))
* added --fetch-depth flag ([#62](https://github.com/lindell/multi-gitter/issues/62)) ([5cdb723](https://github.com/lindell/multi-gitter/commit/5cdb72334f151c4950ffd9763b8ee760dbc3f8a5))
* added --git-type flag ([cb4701e](https://github.com/lindell/multi-gitter/commit/cb4701eb90b98bf585b1a8835368c4cd8f0e0095))
* added --include-subgroups flag ([#131](https://github.com/lindell/multi-gitter/issues/131)) ([eff19a4](https://github.com/lindell/multi-gitter/commit/eff19a4b23030487fa9a3e64553443d2a8fb3133))
* added --merge-type flag ([#64](https://github.com/lindell/multi-gitter/issues/64)) ([dd18402](https://github.com/lindell/multi-gitter/commit/dd18402365c0f41440bd580497cbd12e0738bc7e))
* added --plain-output flag ([#378](https://github.com/lindell/multi-gitter/issues/378)) ([356f471](https://github.com/lindell/multi-gitter/commit/356f471f720e6efe0146bfabcb99b2c214d9c4db))
* added --ssh-auth option ([#215](https://github.com/lindell/multi-gitter/issues/215)) ([f5767a8](https://github.com/lindell/multi-gitter/commit/f5767a86c44562f3191eb8cff0d3084393ed1ac7))
* added base-branch flag ([8c04b8d](https://github.com/lindell/multi-gitter/commit/8c04b8d241f66ec8def92baf8ae27a39a24abcff))
* added close command ([dee1a10](https://github.com/lindell/multi-gitter/commit/dee1a108293d9620d7090cc96598a970f75e3d33))
* added configuration options through config files ([#150](https://github.com/lindell/multi-gitter/issues/150)) ([f38a7ad](https://github.com/lindell/multi-gitter/commit/f38a7ad3ffc9f6aaef60913a6a08006b5b672a93))
* added fork mode ([#128](https://github.com/lindell/multi-gitter/issues/128)) ([f9e7827](https://github.com/lindell/multi-gitter/commit/f9e78273440642be662686912b89ff38123bacf7))
* added Gitea support ([#105](https://github.com/lindell/multi-gitter/issues/105)) ([0f89791](https://github.com/lindell/multi-gitter/commit/0f89791d62fe32f0d2a98f0b735782898976e3f7))
* added GitHub autocompletion ([#84](https://github.com/lindell/multi-gitter/issues/84)) ([5fee0c4](https://github.com/lindell/multi-gitter/commit/5fee0c4b88e802a8be4168f802b79a1701afd3a6))
* added links to printed prs ([#58](https://github.com/lindell/multi-gitter/issues/58)) ([cd76c61](https://github.com/lindell/multi-gitter/commit/cd76c6143a9b008f6be08748b77f7c8acc36aaf9))
* added log-format flag ([4fab0eb](https://github.com/lindell/multi-gitter/commit/4fab0eb5c6e4478b23fac09fc429e0ed3b38b82f))
* added more information to the version command ([#198](https://github.com/lindell/multi-gitter/issues/198)) ([ebf4578](https://github.com/lindell/multi-gitter/commit/ebf457822693000fe04caf4c36a5db70c9feab6c))
* added option to add labels to creates pull requests ([#292](https://github.com/lindell/multi-gitter/issues/292)) ([fccf678](https://github.com/lindell/multi-gitter/commit/fccf678c384ff01cd0247c35860acf0d257e41a7))
* added output option to differentiate log and other output ([c4f52fd](https://github.com/lindell/multi-gitter/commit/c4f52fda0153e4c3898f5b3f59ad1001a505cb52))
* added shell completion command ([c5782a2](https://github.com/lindell/multi-gitter/commit/c5782a2e377ecfc071c82c1db0a775e45215a0cc))
* added skip-pr flag ([#80](https://github.com/lindell/multi-gitter/issues/80)) ([c4b85ea](https://github.com/lindell/multi-gitter/commit/c4b85ea5606a361b13b0a6308f3cfea776f954ad))
* added static flag completion for enums ([#87](https://github.com/lindell/multi-gitter/issues/87)) ([586dd61](https://github.com/lindell/multi-gitter/commit/586dd616418affe1838b4ecfb5714458ffcafd0b))
* **bitbucketserver:** added support for bitbucket server ([#178](https://github.com/lindell/multi-gitter/issues/178)) ([2f7a1b6](https://github.com/lindell/multi-gitter/commit/2f7a1b6e313355a8aa4176cc216bd2d9ad6494a7))
* filter repositories using topics ([#320](https://github.com/lindell/multi-gitter/issues/320)) ([d3c5403](https://github.com/lindell/multi-gitter/commit/d3c54034d56f7826dfa7dbc13851b152334fbf58))
* **github:** added option to add team reviewers ([#351](https://github.com/lindell/multi-gitter/issues/351)) ([bfe05b9](https://github.com/lindell/multi-gitter/commit/bfe05b9b5c53307a8c429278b6491bafe57a2f26))
* **github:** use code search to select repositories ([#399](https://github.com/lindell/multi-gitter/issues/399)) ([d8b75f4](https://github.com/lindell/multi-gitter/commit/d8b75f4ed173c2389dcb3dd89a4213fd8ce78f42))
* **github:** use graphql endpoint to get pull request status ([#242](https://github.com/lindell/multi-gitter/issues/242)) ([60bbbdf](https://github.com/lindell/multi-gitter/commit/60bbbdf526ce6e87d0f952fdae4858e4c6954952))
* **gitlab:** option to change base-url for gitlab ([#69](https://github.com/lindell/multi-gitter/issues/69)) ([147ebe6](https://github.com/lindell/multi-gitter/commit/147ebe67d2902850f06c7575bbe8e43b0372eccd))
* **gitlab:** respect project level squash setting ([#228](https://github.com/lindell/multi-gitter/issues/228)) ([b189661](https://github.com/lindell/multi-gitter/commit/b1896610c6f4f88099848f99163984fbf4de113f))
* **gitlab:** skip archived repos ([#240](https://github.com/lindell/multi-gitter/issues/240)) ([10df83a](https://github.com/lindell/multi-gitter/commit/10df83a80e486dc7f2c8ec085c1111d8fab42cfb))
* interactive mode  ([#167](https://github.com/lindell/multi-gitter/issues/167)) ([7351520](https://github.com/lindell/multi-gitter/commit/73515206bc7201b28e0e1faef7e1009b3e5a34f9))
* moved to built in completion command in cobra 1.2.x ([#163](https://github.com/lindell/multi-gitter/issues/163)) ([81a7187](https://github.com/lindell/multi-gitter/commit/81a7187fce1ab76e6d87bdeee02b268fdb21320b))
* option to skip repositories that are forks ([#341](https://github.com/lindell/multi-gitter/issues/341)) ([941731b](https://github.com/lindell/multi-gitter/commit/941731bfc0a9a89b2abef18286a7a3b06ab5d1db))
* OS and Arch info added to the version command ([#348](https://github.com/lindell/multi-gitter/issues/348)) ([12c0422](https://github.com/lindell/multi-gitter/commit/12c04221fa678bf230a50f4a4386aebe4cfa593f))
* set DRY_RUN when --dry-run is used ([#337](https://github.com/lindell/multi-gitter/issues/337)) ([e4390ee](https://github.com/lindell/multi-gitter/commit/e4390ee8eddc8a37ea73fd5f29d3e9151221900e))
* use repo search to select repositories ([#397](https://github.com/lindell/multi-gitter/issues/397)) ([0f8c2dc](https://github.com/lindell/multi-gitter/commit/0f8c2dcceeb9b635c7405e878df608318b518de4))


### Bug Fixes

* added panic recover on a run repo basis ([#114](https://github.com/lindell/multi-gitter/issues/114)) ([6d44adf](https://github.com/lindell/multi-gitter/commit/6d44adf5ddbf3783bc4a2224c35a923ab599e7c6))
* added support for cancellation of git remote commands ([#296](https://github.com/lindell/multi-gitter/issues/296)) ([83bfbc7](https://github.com/lindell/multi-gitter/commit/83bfbc764a920ffbbedc5a80028b889391360dfd))
* added support for GitLab subgroups with --project ([#171](https://github.com/lindell/multi-gitter/issues/171)) ([25b5d54](https://github.com/lindell/multi-gitter/commit/25b5d543056909fdb1a937118989f06dd4902f80))
* added the number of created pull requests ([#56](https://github.com/lindell/multi-gitter/issues/56)) ([d432430](https://github.com/lindell/multi-gitter/commit/d4324307441ffc74002e1cb4f5c08b83f45a2781))
* allow individual merge failures without aborting all merges ([#310](https://github.com/lindell/multi-gitter/issues/310)) ([e067502](https://github.com/lindell/multi-gitter/commit/e067502f2366e86ce6b979130233d41017f40caf))
* better error message when the same feature and base branch is used ([#281](https://github.com/lindell/multi-gitter/issues/281)) ([ff98e8f](https://github.com/lindell/multi-gitter/commit/ff98e8fff43a3ccc14feb2f24f6c54c134dbc21c))
* better logs when repositories are not used due to permissions ([#226](https://github.com/lindell/multi-gitter/issues/226)) ([668d0b0](https://github.com/lindell/multi-gitter/commit/668d0b0458988dc24de7beb4d087c2cccc7d3167))
* **bitbucket:** ensure username is set ([#212](https://github.com/lindell/multi-gitter/issues/212)) ([a463709](https://github.com/lindell/multi-gitter/commit/a4637093e71c3b667afb58e67439bb2b3c9fe927))
* censor http authentication header ([#185](https://github.com/lindell/multi-gitter/issues/185)) ([633a2cc](https://github.com/lindell/multi-gitter/commit/633a2ccc973070790b0cb644aa9029727a220e20))
* changed the order of version and os in the released archives ([656e9de](https://github.com/lindell/multi-gitter/commit/656e9de09c0bdb9064d47bd2b6adeb23f971eea3))
* cmdgit BranchExist should match full branch name ([#433](https://github.com/lindell/multi-gitter/issues/433)) ([c7a6e69](https://github.com/lindell/multi-gitter/commit/c7a6e69aafb87efafa25248c73d7afa1a8ee76ac))
* commit hooks no longer run with cmd-git implementation ([#157](https://github.com/lindell/multi-gitter/issues/157)) ([ba12d08](https://github.com/lindell/multi-gitter/commit/ba12d08fee2e8cc0ef8015a1761afde747a2622c))
* corrected the name of the REPOSITORY env var ([9b87070](https://github.com/lindell/multi-gitter/commit/9b8707096a85d1106045fb79d13b87c5fe8b99de))
* deleted files are now detected with `git-type: go` ([#273](https://github.com/lindell/multi-gitter/issues/273)) ([b21509c](https://github.com/lindell/multi-gitter/commit/b21509c3b392fbe681c1cc5315ccdc63942abbc2))
* don't panic when a PR branch is deleted ([#385](https://github.com/lindell/multi-gitter/issues/385)) ([fca98e1](https://github.com/lindell/multi-gitter/commit/fca98e134a6640d79aae6bffbe154d42b3393d40))
* downgraded go-diff to fix diff formating ([#156](https://github.com/lindell/multi-gitter/issues/156)) ([6ef43a8](https://github.com/lindell/multi-gitter/commit/6ef43a847f14d5b81745e9978732eebda5bf8ca9))
* ensure pull request exists during conflict resolution ([#369](https://github.com/lindell/multi-gitter/issues/369)) ([2b7166a](https://github.com/lindell/multi-gitter/commit/2b7166a203ddd7d939bc28a90ee26e08adb5a771))
* fixed archive release names ([1462b90](https://github.com/lindell/multi-gitter/commit/1462b901c373c0a1f04c73f38784fac01cb65189))
* fixed brew test command ([fc243e8](https://github.com/lindell/multi-gitter/commit/fc243e8d7d94c9b1793eb7299a893ba2ba14794c))
* fixed bug where base branch was left empty ([64d5e22](https://github.com/lindell/multi-gitter/commit/64d5e225e631f8b3a0dac3fc3145f0168dacba70))
* fixed fury.io token in release ([9deb2d2](https://github.com/lindell/multi-gitter/commit/9deb2d258a6b84c8f002b4aaab0ffc8ca135540a))
* fixed homebrew release ([dece0d8](https://github.com/lindell/multi-gitter/commit/dece0d8ad5e40c20be37eeb9db42dcdfd9eaf4d4))
* fixed license file in release ([506084f](https://github.com/lindell/multi-gitter/commit/506084fd8b17f42a3311524bc0dbcc29ce39c50b))
* fixed multi line commit message ([#222](https://github.com/lindell/multi-gitter/issues/222)) ([995a93c](https://github.com/lindell/multi-gitter/commit/995a93cf552300e3c4b580a0ffa2aeb3cdfd61a7))
* fixed presentation of repos with existing repo ([ac8027b](https://github.com/lindell/multi-gitter/commit/ac8027b3cf6c8df46ae3c4e2b79891c14962f7bc))
* fixed skip-pr flag description ([#127](https://github.com/lindell/multi-gitter/issues/127)) ([1c4e2ac](https://github.com/lindell/multi-gitter/commit/1c4e2acec3fee563eb3cfa7391f63ffd5fc1d61e))
* fixed windows filepaths ([#89](https://github.com/lindell/multi-gitter/issues/89)) ([cb38fc0](https://github.com/lindell/multi-gitter/commit/cb38fc08a084dd7b5b05717b852e9804d52e1720))
* **github:** add retry mechanism to all requests ([#289](https://github.com/lindell/multi-gitter/issues/289)) ([89a642c](https://github.com/lindell/multi-gitter/commit/89a642c551fa032d4389883acd26146109ebfa99))
* **github:** added missing closed pr status ([#290](https://github.com/lindell/multi-gitter/issues/290)) ([9e3644a](https://github.com/lindell/multi-gitter/commit/9e3644a899197470985707112238f9990df99564))
* **github:** allow GitHub app tokens to be used ([#334](https://github.com/lindell/multi-gitter/issues/334)) ([8d86544](https://github.com/lindell/multi-gitter/commit/8d865447765d70c6bdde393ef3ee450355b61ae0))
* **github:** allow only pull permission for print command ([#262](https://github.com/lindell/multi-gitter/issues/262)) ([582c706](https://github.com/lindell/multi-gitter/commit/582c70633e0dcef04699c62b160a57fedb50c00d))
* **github:** chunk get pull requests ([#330](https://github.com/lindell/multi-gitter/issues/330)) ([488cd63](https://github.com/lindell/multi-gitter/commit/488cd6339bd659e88cd89d207ae1d5eb5e40b988))
* **github:** correctly map merged prs ([#268](https://github.com/lindell/multi-gitter/issues/268)) ([0474040](https://github.com/lindell/multi-gitter/commit/0474040c86d7c213987fb5dfc503d6c46288cb3a))
* **github:** fixed that the fetching of pullrequests always pull the latest pr ([#195](https://github.com/lindell/multi-gitter/issues/195)) ([aa33af8](https://github.com/lindell/multi-gitter/commit/aa33af834d71e7122955b5023ab028c2d5fa42f8))
* **github:** ignore branch deletion error if branch is already deleted ([#189](https://github.com/lindell/multi-gitter/issues/189)) ([d63d041](https://github.com/lindell/multi-gitter/commit/d63d04184dc10d3c6538676dacdd63d973d06e02))
* **github:** made sure GraphQL requests, with non-GraphQL errors is reported properly ([#301](https://github.com/lindell/multi-gitter/issues/301)) ([d7e1fda](https://github.com/lindell/multi-gitter/commit/d7e1fda392d1046fa93c2f0304295c0fd7872292))
* **github:** use other format when cloning with token to support more token formats ([#302](https://github.com/lindell/multi-gitter/issues/302)) ([a74cc60](https://github.com/lindell/multi-gitter/commit/a74cc60d8257e69a99f224150f45400223b5911b))
* **gitlab:** close MR instead of deleting it ([#230](https://github.com/lindell/multi-gitter/issues/230)) ([af2c2d9](https://github.com/lindell/multi-gitter/commit/af2c2d9b86a8cd0c6c09e3667a9f7c7689cf915c))
* **gitlab:** made entire owner structure a part of logging and $REPOSITORY ([ccac98a](https://github.com/lindell/multi-gitter/commit/ccac98a8c6ab7e4dc23e19519de7ccd958e1be43))
* **gitlab:** make sure entire project path is shown ([#392](https://github.com/lindell/multi-gitter/issues/392)) ([280dfc7](https://github.com/lindell/multi-gitter/commit/280dfc7b6ed5c8a5345f2207e58e81202a451e33))
* **gitlab:** only list projects with Merge Requests enabled ([#184](https://github.com/lindell/multi-gitter/issues/184)) ([5d45121](https://github.com/lindell/multi-gitter/commit/5d4512112715dbe9ce7cba214531ce93c8b1a360))
* improved error messages for common problems with the script ([de9e525](https://github.com/lindell/multi-gitter/commit/de9e5259d2bd900abf72c56f40a76f223cbfffd0))
* input description for platform args ([#321](https://github.com/lindell/multi-gitter/issues/321)) ([6816c86](https://github.com/lindell/multi-gitter/commit/6816c86ea534112ef55c24036940f2c44b50e1ab))
* made remove branch on merge the default behaviour for GitLab merge ([#135](https://github.com/lindell/multi-gitter/issues/135)) ([9cc5983](https://github.com/lindell/multi-gitter/commit/9cc5983407c3b5be4a42c55dbd7c4b03f54d3f23))
* made sure any tokens output in the logs are now censored ([#143](https://github.com/lindell/multi-gitter/issues/143)) ([0e5cee7](https://github.com/lindell/multi-gitter/commit/0e5cee7ecd6dde23d21869058cc383e83b232703))
* make sure autocompletion works with settings from config files ([#217](https://github.com/lindell/multi-gitter/issues/217)) ([45e855c](https://github.com/lindell/multi-gitter/commit/45e855c314738f710092f2babb6a213727a63467))
* make sure GitHub's secondary rate limit is not reached ([#207](https://github.com/lindell/multi-gitter/issues/207)) ([8a5fabd](https://github.com/lindell/multi-gitter/commit/8a5fabdc9e54bdfbba421a466bd323aae8114bdd))
* make sure gitignore is used ([#119](https://github.com/lindell/multi-gitter/issues/119)) ([f33dee9](https://github.com/lindell/multi-gitter/commit/f33dee9a7acd798ab6ad0a7255351c50c9bd456e))
* mitigated script injection in test install action ([#404](https://github.com/lindell/multi-gitter/issues/404)) ([f77ef91](https://github.com/lindell/multi-gitter/commit/f77ef91416f3b13080372bb00831824eea37c130))
* multi-gitter does now only fetch the base and head branch ([b272644](https://github.com/lindell/multi-gitter/commit/b272644355d9291c23de8f028a3132de5a5eb99e))
* push hooks no longer run with cmd-git implementation ([#159](https://github.com/lindell/multi-gitter/issues/159)) ([7360c0d](https://github.com/lindell/multi-gitter/commit/7360c0d14b83be627325d0b4ea95177e71c2a565))
* remove returning of archived and disabled Github repos ([da2d635](https://github.com/lindell/multi-gitter/commit/da2d635c202004a855217aa941af623b3c5966d7))
* retry when encounting the GitHub rate limit ([#280](https://github.com/lindell/multi-gitter/issues/280)) ([008a26a](https://github.com/lindell/multi-gitter/commit/008a26ae0182c8e69cc0679a0d7ee776415c1d62))
* skip running git diff if debug or lower is not set ([#113](https://github.com/lindell/multi-gitter/issues/113)) ([5189374](https://github.com/lindell/multi-gitter/commit/51893745153e7825339f7398e844bf6d53404cc8))
* typo where archived should be achieved ([#125](https://github.com/lindell/multi-gitter/issues/125)) ([5373ea8](https://github.com/lindell/multi-gitter/commit/5373ea8fd37e39ce1eb8edbb860af85faa47e370))
* **typo:** fixed panicked typo ([#423](https://github.com/lindell/multi-gitter/issues/423)) ([9b5774a](https://github.com/lindell/multi-gitter/commit/9b5774ad8e95f849d9af0b76a58183c7bbdf0462))
* typos in docs and code ([#303](https://github.com/lindell/multi-gitter/issues/303)) ([45ddb60](https://github.com/lindell/multi-gitter/commit/45ddb60cd438e9251b670c0884cc851c29f09d8f))
* update error message to fix grammar ([#345](https://github.com/lindell/multi-gitter/issues/345)) ([5a4c990](https://github.com/lindell/multi-gitter/commit/5a4c990b7f9325d1a5c5b4ee619e792478517754))


### Dependencies

* bump golang.org/x/net from 0.3.0 to 0.7.0 ([#325](https://github.com/lindell/multi-gitter/issues/325)) ([9d0ea43](https://github.com/lindell/multi-gitter/commit/9d0ea43c5ad46a76d8a596e0d33082ce51579aa6))
* update github.com/eiannone/keyboard digest to 0d22619 ([f38b2f5](https://github.com/lindell/multi-gitter/commit/f38b2f5b2bf17afab06ed709c6b5ba4c9c9a6075))
* update github.com/gfleury/go-bitbucket-v1 commit hash to dff2223 ([f570ee5](https://github.com/lindell/multi-gitter/commit/f570ee5086369b50091f5ad21f6762d96d93782c))
* update github.com/gfleury/go-bitbucket-v1 digest to 6e30c57 ([#382](https://github.com/lindell/multi-gitter/issues/382)) ([c201f42](https://github.com/lindell/multi-gitter/commit/c201f425b4fbde46c3d27fa380ab099e0a6f8608))
* update github.com/gfleury/go-bitbucket-v1 digest to 711d7d5 ([0ecf9ee](https://github.com/lindell/multi-gitter/commit/0ecf9ee2e17b162218b251db63b5e9946369df39))
* update github.com/gfleury/go-bitbucket-v1 digest to 8d7be58 ([1d8a19f](https://github.com/lindell/multi-gitter/commit/1d8a19f8316dd8874a0755189902dce24acd09c5))
* update golang.org/x/exp digest to 02704c9 ([f7e273e](https://github.com/lindell/multi-gitter/commit/f7e273e6ebf4ee52d16e76e8472de87eb9c3a52e))
* update golang.org/x/exp digest to 6522937 ([cab859f](https://github.com/lindell/multi-gitter/commit/cab859fdb40c181e9b1dc58448290eec800f6aae))
* update golang.org/x/oauth2 commit hash to 08078c5 ([8b94c50](https://github.com/lindell/multi-gitter/commit/8b94c50acf0df6c5af2a9b4e81a7aea296575bc8))
* update golang.org/x/oauth2 commit hash to 0b49973 ([ee02a57](https://github.com/lindell/multi-gitter/commit/ee02a57a4512e9e8b29770ec610974ad5ecdf7d2))
* update golang.org/x/oauth2 commit hash to 22b0ada ([#92](https://github.com/lindell/multi-gitter/issues/92)) ([335eee3](https://github.com/lindell/multi-gitter/commit/335eee37c02c54fa7d006ff0aab837b424f7d514))
* update golang.org/x/oauth2 commit hash to 2bc19b1 ([858441a](https://github.com/lindell/multi-gitter/commit/858441a9822b6f86d9e68216742f550eb80f7e05))
* update golang.org/x/oauth2 commit hash to 81ed05c ([#107](https://github.com/lindell/multi-gitter/issues/107)) ([b529c3f](https://github.com/lindell/multi-gitter/commit/b529c3f113ccda92ee0981c48b0c26c74facb142))
* update golang.org/x/oauth2 commit hash to 9317641 ([242cdd0](https://github.com/lindell/multi-gitter/commit/242cdd09b5c62813336059f8ac73cf662bc4f71e))
* update golang.org/x/oauth2 commit hash to a41e5a7 ([234ce36](https://github.com/lindell/multi-gitter/commit/234ce36753e5eec8d73700a4b65e4ee8ad0773a7))
* update golang.org/x/oauth2 commit hash to d3ed898 ([feea168](https://github.com/lindell/multi-gitter/commit/feea168f7a2d44d9fe08c8b1a995dfc5b213f7ce))
* update golang.org/x/oauth2 commit hash to f6687ab ([cab768a](https://github.com/lindell/multi-gitter/commit/cab768a1a6bf93b8f113b0b7221db7a4bab375cd))
* update golang.org/x/oauth2 commit hash to f9ce19e ([#66](https://github.com/lindell/multi-gitter/issues/66)) ([64d9095](https://github.com/lindell/multi-gitter/commit/64d90952856fdfd0517cf03bb752603c708ff6b9))
* update golang.org/x/oauth2 digest to 128564f ([63eab95](https://github.com/lindell/multi-gitter/commit/63eab95d69715ca1794f05447e3a588b70d99e62))
* update golang.org/x/oauth2 digest to 2104d58 ([c4605e2](https://github.com/lindell/multi-gitter/commit/c4605e2c8db5440fc68b5b779ee46b360c4459db))
* update golang.org/x/oauth2 digest to f213421 ([dfd1837](https://github.com/lindell/multi-gitter/commit/dfd18376c44f4b9a36628142fc52277b0aa48ec0))
* update module code.gitea.io/sdk/gitea to v0.15.0 ([1b0ac09](https://github.com/lindell/multi-gitter/commit/1b0ac094015c4b398147c3fc8759e83462b656b7))
* update module code.gitea.io/sdk/gitea to v0.15.1 ([c8b4ab8](https://github.com/lindell/multi-gitter/commit/c8b4ab8cf55ace9e4f4751b540ba2e0dc7fd2807))
* update module code.gitea.io/sdk/gitea to v0.16.0 ([ff04448](https://github.com/lindell/multi-gitter/commit/ff0444885bfbf657b4813efdb53986af67b0ac58))
* update module code.gitea.io/sdk/gitea to v0.17.0 ([3949d2f](https://github.com/lindell/multi-gitter/commit/3949d2f4a3d7f315ca48a6c8c18910a1bcce3443))
* update module github.com/go-git/go-git/v5 to v5.10.0 ([077d637](https://github.com/lindell/multi-gitter/commit/077d637789ce7943212ca01827dd06fe314ad75d))
* update module github.com/go-git/go-git/v5 to v5.10.1 ([40f62a1](https://github.com/lindell/multi-gitter/commit/40f62a11cacfd429e232451aeab9e11825de9ce3))
* update module github.com/go-git/go-git/v5 to v5.11.0 [security] ([5fef1cc](https://github.com/lindell/multi-gitter/commit/5fef1ccc936c2d6f31a8d4ae7ac0ebed6409f2c1))
* update module github.com/go-git/go-git/v5 to v5.3.0 ([905dbdb](https://github.com/lindell/multi-gitter/commit/905dbdbfa5b420ee985bed2ff58cfb2399b051b7))
* update module github.com/go-git/go-git/v5 to v5.4.1 ([fe45f2e](https://github.com/lindell/multi-gitter/commit/fe45f2e9ad2031ae4f436271e4a072101ba80805))
* update module github.com/go-git/go-git/v5 to v5.4.2 ([016f54d](https://github.com/lindell/multi-gitter/commit/016f54d39a8df80558b7c46880a7dfabd16c7e28))
* update module github.com/go-git/go-git/v5 to v5.5.1 ([#312](https://github.com/lindell/multi-gitter/issues/312)) ([7c9136c](https://github.com/lindell/multi-gitter/commit/7c9136c7427fcb1d103c0f1e004059f7a793f25d))
* update module github.com/go-git/go-git/v5 to v5.5.2 ([83032fb](https://github.com/lindell/multi-gitter/commit/83032fb7e4cd197c703c7eb8e272938a315c5ef5))
* update module github.com/go-git/go-git/v5 to v5.6.1 ([1fddf2e](https://github.com/lindell/multi-gitter/commit/1fddf2e26cd3cb9f201325decc356d26a1f1444e))
* update module github.com/go-git/go-git/v5 to v5.7.0 ([d407eb9](https://github.com/lindell/multi-gitter/commit/d407eb9ef7bb06ca54b0d28ad992c8832ecb6c41))
* update module github.com/go-git/go-git/v5 to v5.8.1 ([8de2001](https://github.com/lindell/multi-gitter/commit/8de20010a3c8c7ac5eecb6c366975d164e9a8409))
* update module github.com/go-git/go-git/v5 to v5.9.0 ([a954c7c](https://github.com/lindell/multi-gitter/commit/a954c7c0c20cb0a930203c92143456c1571d2e2f))
* update module github.com/google/go-github/v32 to v33 ([#82](https://github.com/lindell/multi-gitter/issues/82)) ([1c48de3](https://github.com/lindell/multi-gitter/commit/1c48de3a81a64cbac6481b3260bdc3512e98a34f))
* update module github.com/google/go-github/v33 to v34 ([#93](https://github.com/lindell/multi-gitter/issues/93)) ([03d3278](https://github.com/lindell/multi-gitter/commit/03d327835bb7a99362e0b13224200a41d068a642))
* update module github.com/google/go-github/v33 to v35 ([#110](https://github.com/lindell/multi-gitter/issues/110)) ([b6c8667](https://github.com/lindell/multi-gitter/commit/b6c8667f1ca48c62b1ec1703f8afa1664dfeca95))
* update module github.com/google/go-github/v35 to v36 ([#162](https://github.com/lindell/multi-gitter/issues/162)) ([893d9ea](https://github.com/lindell/multi-gitter/commit/893d9eae5dd5f8abcf6c00cb233957aea532d1c2))
* update module github.com/google/go-github/v36 to v37 ([213a1c6](https://github.com/lindell/multi-gitter/commit/213a1c6cc603cec49f889ffe52dd50d22f33ab44))
* update module github.com/google/go-github/v37 to v38 ([#176](https://github.com/lindell/multi-gitter/issues/176)) ([f15aaad](https://github.com/lindell/multi-gitter/commit/f15aaad21ba92a4d2c05c039f0b7f8963f245e75))
* update module github.com/google/go-github/v38 to v39 ([#191](https://github.com/lindell/multi-gitter/issues/191)) ([5088532](https://github.com/lindell/multi-gitter/commit/508853232485cd4dd4886f46fead14fa71d7ae59))
* update module github.com/google/go-github/v39 to v39.2.0 ([45f20a0](https://github.com/lindell/multi-gitter/commit/45f20a0e070e71231d785d1bb12cc04ee0d0e2e2))
* update module github.com/google/go-github/v47 to v48 ([#311](https://github.com/lindell/multi-gitter/issues/311)) ([c3ffe09](https://github.com/lindell/multi-gitter/commit/c3ffe095a56d450ea2d5aa81384a4578cc069a12))
* update module github.com/google/go-github/v48 to v50 ([#319](https://github.com/lindell/multi-gitter/issues/319)) ([6f1fe90](https://github.com/lindell/multi-gitter/commit/6f1fe9009131d1d0326c7f63cf6d07b2cd1983f3))
* update module github.com/google/go-github/v50 to v50.1.0 ([366f201](https://github.com/lindell/multi-gitter/commit/366f201f77bc46aaeedfa34fd6ed40dc99a007f6))
* update module github.com/google/go-github/v50 to v50.2.0 ([#359](https://github.com/lindell/multi-gitter/issues/359)) ([481c5da](https://github.com/lindell/multi-gitter/commit/481c5da7223eb798bd04b21eb839aa6aa563a763))
* update module github.com/google/go-github/v50 to v53 ([#362](https://github.com/lindell/multi-gitter/issues/362)) ([ff85919](https://github.com/lindell/multi-gitter/commit/ff85919b1a264ce042e9e2bbca24e356be15e623))
* update module github.com/google/go-github/v53 to v55 ([#390](https://github.com/lindell/multi-gitter/issues/390)) ([d92cb29](https://github.com/lindell/multi-gitter/commit/d92cb29cc96e8a78d19ccdd1816643510169f6c9))
* update module github.com/google/go-github/v55 to v56 ([#412](https://github.com/lindell/multi-gitter/issues/412)) ([946afac](https://github.com/lindell/multi-gitter/commit/946afac76f4a1884823da591d04236daaf3fb795))
* update module github.com/google/go-github/v56 to v57 ([#445](https://github.com/lindell/multi-gitter/issues/445)) ([28a3279](https://github.com/lindell/multi-gitter/commit/28a32793edd1c63bb66b2008ef333f2f7b31c07b))
* update module github.com/mitchellh/mapstructure to v1.4.3 ([ea97fa2](https://github.com/lindell/multi-gitter/commit/ea97fa2010378b353c0a43f832aced2054ae182d))
* update module github.com/mitchellh/mapstructure to v1.5.0 ([152c3b1](https://github.com/lindell/multi-gitter/commit/152c3b1aeebeb996221d28ca6072342e9188493c))
* update module github.com/sergi/go-diff to v1.2.0 ([#116](https://github.com/lindell/multi-gitter/issues/116)) ([0273abe](https://github.com/lindell/multi-gitter/commit/0273abeba104e2ce522b4a97b1498341ad9e41d6))
* update module github.com/sirupsen/logrus to v1.8.0 ([8c132b4](https://github.com/lindell/multi-gitter/commit/8c132b410baef2812a5525727147eb11f939a870))
* update module github.com/sirupsen/logrus to v1.8.1 ([31dad70](https://github.com/lindell/multi-gitter/commit/31dad70383ab5c6d742e393ad97c32c610f85c2b))
* update module github.com/sirupsen/logrus to v1.9.0 ([1c8201a](https://github.com/lindell/multi-gitter/commit/1c8201a3b748b511effeb02b98ae1485f0e97b31))
* update module github.com/sirupsen/logrus to v1.9.3 ([3f0d4c7](https://github.com/lindell/multi-gitter/commit/3f0d4c707636a595404839deecc6be9204cf0407))
* update module github.com/spf13/cobra to v1.3.0 ([afbe211](https://github.com/lindell/multi-gitter/commit/afbe21104568c9550e7cbaa48342ed5bcf3de3b8))
* update module github.com/spf13/cobra to v1.4.0 ([d378aa0](https://github.com/lindell/multi-gitter/commit/d378aa00a1b54f8482ed4d18d94ef47651b44273))
* update module github.com/spf13/cobra to v1.5.0 ([eb15db4](https://github.com/lindell/multi-gitter/commit/eb15db4fc137d5a23f3d137a778e4235b26b4e74))
* update module github.com/spf13/cobra to v1.6.1 ([cfc2861](https://github.com/lindell/multi-gitter/commit/cfc2861089ac4e8697202c6bcf5daf01021fd8bf))
* update module github.com/spf13/cobra to v1.7.0 ([#360](https://github.com/lindell/multi-gitter/issues/360)) ([212a9b6](https://github.com/lindell/multi-gitter/commit/212a9b686db37c58205d7f2e20778ec592e15388))
* update module github.com/spf13/cobra to v1.8.0 ([f976ac2](https://github.com/lindell/multi-gitter/commit/f976ac2d91254e92f5e937f9e0fad99802b8ce79))
* update module github.com/spf13/viper to v1.10.1 ([6052e73](https://github.com/lindell/multi-gitter/commit/6052e73c7e2f1a790fa4c23732aec0c23dcd131c))
* update module github.com/spf13/viper to v1.11.0 ([a8e01ef](https://github.com/lindell/multi-gitter/commit/a8e01ef104ee09a38fe3d3551caa97fb5d00dc0e))
* update module github.com/spf13/viper to v1.12.0 ([ffdf7ae](https://github.com/lindell/multi-gitter/commit/ffdf7ae901b1f98fdd279fd5d1aeb72095bf2ad2))
* update module github.com/spf13/viper to v1.13.0 ([099a9c9](https://github.com/lindell/multi-gitter/commit/099a9c9f3e5341b207ec298356319303df67a31f))
* update module github.com/spf13/viper to v1.14.0 ([294eced](https://github.com/lindell/multi-gitter/commit/294eced3db7de848f4c721a651e4833007bb3281))
* update module github.com/spf13/viper to v1.16.0 ([#364](https://github.com/lindell/multi-gitter/issues/364)) ([af41308](https://github.com/lindell/multi-gitter/commit/af4130863a62217bcf2c5f4fda163b521112a7a0))
* update module github.com/spf13/viper to v1.17.0 ([cfe3409](https://github.com/lindell/multi-gitter/commit/cfe3409acee757b78c542e106436e7ca177cf82b))
* update module github.com/spf13/viper to v1.18.2 ([29a205e](https://github.com/lindell/multi-gitter/commit/29a205e9db876a0baf6091415f47085bb2876b5e))
* update module github.com/spf13/viper to v1.9.0 ([becff1f](https://github.com/lindell/multi-gitter/commit/becff1f7d6fd755565601fd6eb4d321cac6d54a2))
* update module github.com/stretchr/testify to v1.7.1 ([7506930](https://github.com/lindell/multi-gitter/commit/750693020fffae4aef811558fcdd65abf3662935))
* update module github.com/stretchr/testify to v1.8.0 ([81a757c](https://github.com/lindell/multi-gitter/commit/81a757c88d3942ff5295a8edfcec6ad99ff9e669))
* update module github.com/stretchr/testify to v1.8.1 ([e5158dc](https://github.com/lindell/multi-gitter/commit/e5158dc7b3536ec681d090ae18d5db88feb19fe0))
* update module github.com/stretchr/testify to v1.8.2 ([4f3ce98](https://github.com/lindell/multi-gitter/commit/4f3ce98307815934149fbef6da211280bbd581dd))
* update module github.com/stretchr/testify to v1.8.4 ([#361](https://github.com/lindell/multi-gitter/issues/361)) ([2686055](https://github.com/lindell/multi-gitter/commit/2686055d74f2f2af416ff6ed05eb4f47647648be))
* update module github.com/xanzy/go-gitlab to v0.45.0 ([9e1bc9a](https://github.com/lindell/multi-gitter/commit/9e1bc9a163f151fddba8f24f44934111d7e9b810))
* update module github.com/xanzy/go-gitlab to v0.46.0 ([f0a3503](https://github.com/lindell/multi-gitter/commit/f0a350323dcee32acf20265b5ebcedce2f1531b9))
* update module github.com/xanzy/go-gitlab to v0.47.0 ([92a18a3](https://github.com/lindell/multi-gitter/commit/92a18a3fd27136b5a46c2de423c109d28ad9da71))
* update module github.com/xanzy/go-gitlab to v0.49.0 ([#109](https://github.com/lindell/multi-gitter/issues/109)) ([597d8b4](https://github.com/lindell/multi-gitter/commit/597d8b41751b0cf90bc4743fc367ed72487ae35f))
* update module github.com/xanzy/go-gitlab to v0.50.1 ([cac5518](https://github.com/lindell/multi-gitter/commit/cac5518094e4cc82bf2ad6d47c42a593f3031034))
* update module github.com/xanzy/go-gitlab to v0.51.1 ([78fb3dc](https://github.com/lindell/multi-gitter/commit/78fb3dca3b3d5e1aff66799814e3c4a92edda0d7))
* update module github.com/xanzy/go-gitlab to v0.52.2 ([c88f791](https://github.com/lindell/multi-gitter/commit/c88f79152e1b99381fcc8de249149aa49089c77e))
* update module github.com/xanzy/go-gitlab to v0.55.1 ([ea0eacc](https://github.com/lindell/multi-gitter/commit/ea0eacce1c3a711b9f3c0d51c2b714ebc72ddb23))
* update module github.com/xanzy/go-gitlab to v0.68.0 ([78ca7f3](https://github.com/lindell/multi-gitter/commit/78ca7f343b846c56a7fdc76f59ec6a7591d7c3da))
* update module github.com/xanzy/go-gitlab to v0.68.2 ([dab0985](https://github.com/lindell/multi-gitter/commit/dab09851cf9b694096ec11084568685efcd860ff))
* update module github.com/xanzy/go-gitlab to v0.73.1 ([badc233](https://github.com/lindell/multi-gitter/commit/badc23353c3f611031fa70483b6892cf8f2e0b99))
* update module github.com/xanzy/go-gitlab to v0.76.0 ([0ce73dc](https://github.com/lindell/multi-gitter/commit/0ce73dc36d03ea6fa94a7203c0825028e1a97c0b))
* update module github.com/xanzy/go-gitlab to v0.77.0 ([5e39817](https://github.com/lindell/multi-gitter/commit/5e3981748c08e5ee4d63b0a1de1b8b4df50b84fc))
* update module github.com/xanzy/go-gitlab to v0.86.0 ([#365](https://github.com/lindell/multi-gitter/issues/365)) ([79b9a0e](https://github.com/lindell/multi-gitter/commit/79b9a0e57a45ca5c55af310d1a619018166d2362))
* update module github.com/xanzy/go-gitlab to v0.88.0 ([c0f0bb4](https://github.com/lindell/multi-gitter/commit/c0f0bb4ba45c9834fe67a44774a96a984a16913c))
* update module github.com/xanzy/go-gitlab to v0.90.0 ([fd66cd2](https://github.com/lindell/multi-gitter/commit/fd66cd200c08b911fb19dae40dafa6c91170a269))
* update module github.com/xanzy/go-gitlab to v0.92.1 ([19c2925](https://github.com/lindell/multi-gitter/commit/19c292523f3442766f09d5811ec9b38ecb26c9bf))
* update module github.com/xanzy/go-gitlab to v0.92.3 ([3e6cd75](https://github.com/lindell/multi-gitter/commit/3e6cd75ad0923ad535f2125f1c94cba6b9072f05))
* update module github.com/xanzy/go-gitlab to v0.93.2 ([76fd29b](https://github.com/lindell/multi-gitter/commit/76fd29b07a7561223adb474bf4fd5eedae6a4298))
* update module github.com/xanzy/go-gitlab to v0.94.0 ([38666d9](https://github.com/lindell/multi-gitter/commit/38666d9b5ae79783ea161a20ad566480e5c1b23a))
* update module github.com/xanzy/go-gitlab to v0.95.2 ([#446](https://github.com/lindell/multi-gitter/issues/446)) ([61f0110](https://github.com/lindell/multi-gitter/commit/61f0110512e4ff464315e3791e6393e62830f112))
* update module go to 1.19 ([2fdbcde](https://github.com/lindell/multi-gitter/commit/2fdbcde5f6160ece5b8fac38d65761278b9f0f6a))
* update module go-github to v47 ([#286](https://github.com/lindell/multi-gitter/issues/286)) ([68fb1e2](https://github.com/lindell/multi-gitter/commit/68fb1e28af23bbdd8c77323e572a3c43d06a4445))
* update module golang.org/x/oauth2 to v0.10.0 ([a82fc38](https://github.com/lindell/multi-gitter/commit/a82fc380c95ab1bec3a20a2da0465444b8d91cd2))
* update module golang.org/x/oauth2 to v0.11.0 ([7820c22](https://github.com/lindell/multi-gitter/commit/7820c228f69386a12df3ddc4c908bb1940b9aa2e))
* update module golang.org/x/oauth2 to v0.12.0 ([b9c258a](https://github.com/lindell/multi-gitter/commit/b9c258a2b2924ebafa37864159154dc7ed09449c))
* update module golang.org/x/oauth2 to v0.13.0 ([3b3e9c0](https://github.com/lindell/multi-gitter/commit/3b3e9c009d1d4a3634ce8decb75625f2ed8e4d0d))
* update module golang.org/x/oauth2 to v0.15.0 ([3874ccb](https://github.com/lindell/multi-gitter/commit/3874ccb93ab76c3a2724358f825880b750a5dbea))
* update module golang.org/x/oauth2 to v0.3.0 ([4ae8184](https://github.com/lindell/multi-gitter/commit/4ae81844ffe59a397753f046ae9cb87e2c1bac48))
* update module golang.org/x/oauth2 to v0.5.0 ([#327](https://github.com/lindell/multi-gitter/issues/327)) ([fefceee](https://github.com/lindell/multi-gitter/commit/fefceeea27305bb94702b53ad5cfa0bfb6cf973f))
* update module golang.org/x/oauth2 to v0.9.0 ([#363](https://github.com/lindell/multi-gitter/issues/363)) ([58030bd](https://github.com/lindell/multi-gitter/commit/58030bd778db562dedcb6cf48eef9dbb2a8039d9))
* update module spf13/cobra to v1.1.3 ([7a32bb6](https://github.com/lindell/multi-gitter/commit/7a32bb615e0969aa41618885023d797ea101cf5b))
* update module stretchr/testify to v1.7.0 ([0a06a24](https://github.com/lindell/multi-gitter/commit/0a06a247d93d34986504608d1ff437aa17869d53))
* update module xanzy/go-gitlab to v0.40.2 ([ce33ff5](https://github.com/lindell/multi-gitter/commit/ce33ff5e69a9d984199c0f00e0b8c57ef6bfbc93))
* update module xanzy/go-gitlab to v0.41.0 ([f713ee2](https://github.com/lindell/multi-gitter/commit/f713ee227b1013d3e2f293fa4d50dbdbbf980b17))
* update module xanzy/go-gitlab to v0.42.0 ([fd8e373](https://github.com/lindell/multi-gitter/commit/fd8e3737e9db7348cf271adab5d6a958e1a794f4))
* update module xanzy/go-gitlab to v0.43.0 ([1a44511](https://github.com/lindell/multi-gitter/commit/1a44511a7cb27aefb3e0e6d8e4309e3fc78f4756))
* update module xanzy/go-gitlab to v0.44.0 ([53f834b](https://github.com/lindell/multi-gitter/commit/53f834b29f9e801a9fa5d8416ad18b22635ab058))


### Miscellaneous

* added brew install ([#73](https://github.com/lindell/multi-gitter/issues/73)) ([3f56a4a](https://github.com/lindell/multi-gitter/commit/3f56a4aefe6a984b781e2e5792b929ea6b6962e3))
* added CODEOWNERS file ([7b85777](https://github.com/lindell/multi-gitter/commit/7b8577798fbcc3159dfc06d920ee5b33183f0ce9))
* bump github.com/sirupsen/logrus from 1.2.0 to 1.7.0 ([157ffa9](https://github.com/lindell/multi-gitter/commit/157ffa9833c711fec8585b28b544b9ba0231315d))
* bump github.com/xanzy/go-gitlab from 0.39.0 to 0.40.0 ([e65e1a8](https://github.com/lindell/multi-gitter/commit/e65e1a8a480202e81a3b2488d0c0d350fa3f265d))
* bump github.com/xanzy/go-gitlab from 0.40.0 to 0.40.1 ([550f302](https://github.com/lindell/multi-gitter/commit/550f302b23ac1301e84994846766c95d2011fb6e))
* fixed function documentations ([33e8995](https://github.com/lindell/multi-gitter/commit/33e89956b36cf59e605a4a7a8cf6239327224c54))
* fixed gitlab api pointer change ([5c39b22](https://github.com/lindell/multi-gitter/commit/5c39b2208cb3480eaab816a18d029bd9b20b7ddc))
* **gitlab:** fixed subgroup spelling ([46f178a](https://github.com/lindell/multi-gitter/commit/46f178ab097fefb3d07b65f3522e3e516638107a))
* improved logging and added stack trace if --log-level=trace is used ([#138](https://github.com/lindell/multi-gitter/issues/138)) ([abccc5f](https://github.com/lindell/multi-gitter/commit/abccc5f28ba22e3b2b99d6d3ee1c513a213caea7))
* improved the base-url description ([5e7ec24](https://github.com/lindell/multi-gitter/commit/5e7ec248b5dee732fc276d7c40a62a8c4fb76c1c))
* **master:** release 0.47.1 ([#372](https://github.com/lindell/multi-gitter/issues/372)) ([27cef05](https://github.com/lindell/multi-gitter/commit/27cef05d8559eeb3f953cd17e7212e9d7abea8cc))
* **master:** release 0.47.2 ([#374](https://github.com/lindell/multi-gitter/issues/374)) ([6a20d4f](https://github.com/lindell/multi-gitter/commit/6a20d4f5d0836feb31658697dd9f13864113c15f))
* **master:** release 0.47.3 ([#375](https://github.com/lindell/multi-gitter/issues/375)) ([313a02e](https://github.com/lindell/multi-gitter/commit/313a02ea361a20e7b4c8aacbc432ffa39acbe8b6))
* **master:** release 0.47.4 ([#379](https://github.com/lindell/multi-gitter/issues/379)) ([e410ded](https://github.com/lindell/multi-gitter/commit/e410ded374700e1d5caf04bd9b2899689d76fb07))
* **master:** release 0.48.0 ([#389](https://github.com/lindell/multi-gitter/issues/389)) ([f893e91](https://github.com/lindell/multi-gitter/commit/f893e912308f5d3fdbac8948562aa2c259703385))
* **master:** release 0.48.1 ([#405](https://github.com/lindell/multi-gitter/issues/405)) ([adf5f18](https://github.com/lindell/multi-gitter/commit/adf5f18bfadfa76d387972c3bfe2d161d45476dd))
* **master:** release 0.48.2 ([#424](https://github.com/lindell/multi-gitter/issues/424)) ([ff82ba7](https://github.com/lindell/multi-gitter/commit/ff82ba76d79407189e43e3b1fb37b51c0bdd68f0))
* **master:** release 0.49.0 ([#428](https://github.com/lindell/multi-gitter/issues/428)) ([89ef3b7](https://github.com/lindell/multi-gitter/commit/89ef3b77e51433735e667d0953c696b724f81b6f))
* **master:** release 0.49.1 ([#434](https://github.com/lindell/multi-gitter/issues/434)) ([1693f01](https://github.com/lindell/multi-gitter/commit/1693f018a8ea273e3787b9d4e4606dc6a2f2a860))
* moved brew release to Formula folder ([d1ae864](https://github.com/lindell/multi-gitter/commit/d1ae8644cd2a4e6138b3f23d23e792509ff7b3ef))
* release 0.15.1 ([06ce4af](https://github.com/lindell/multi-gitter/commit/06ce4af65668555e1f26b4e0c3aea611270717cb))
* release 0.15.2 ([c9e1596](https://github.com/lindell/multi-gitter/commit/c9e1596d2f4abbe246ac6f1981f0c227c597543f))
* release 0.15.3 ([be685d4](https://github.com/lindell/multi-gitter/commit/be685d4a4e096e6cb1edc6b6fe8cfbe6c648077a))
* release 0.16.0 ([4b3205c](https://github.com/lindell/multi-gitter/commit/4b3205cc40d39235facc3bf5118b889750b44046))
* release 0.16.1 ([1795895](https://github.com/lindell/multi-gitter/commit/1795895960a077d44d0d0f856234d051b6c19304))
* release 0.16.2 ([1be208d](https://github.com/lindell/multi-gitter/commit/1be208d0f3c9a4e27caaabf56d8da640b72f73d0))
* release 0.16.3 ([e30c408](https://github.com/lindell/multi-gitter/commit/e30c408dff1e3318b311e0296731b03e010b334a))
* release 0.16.4 ([0d2207c](https://github.com/lindell/multi-gitter/commit/0d2207c5560978a302596d4b39eed26b79907237))
* release 0.17.0 ([#60](https://github.com/lindell/multi-gitter/issues/60)) ([ece4247](https://github.com/lindell/multi-gitter/commit/ece4247cb587d1dbe690d81f78df6d2a0f21bb87))
* release 0.18.0 ([#63](https://github.com/lindell/multi-gitter/issues/63)) ([f87eb46](https://github.com/lindell/multi-gitter/commit/f87eb46df9ebca83a72d54ad78bd54721694aeae))
* release 0.19.0 ([#65](https://github.com/lindell/multi-gitter/issues/65)) ([f56465f](https://github.com/lindell/multi-gitter/commit/f56465face5a0cbcde705bbeadcfac9c5ca116b1))
* release 0.19.1 ([#67](https://github.com/lindell/multi-gitter/issues/67)) ([6953051](https://github.com/lindell/multi-gitter/commit/695305170ec9f489d6870847b25e6191b62b1fe5))
* release 0.20.0 ([#70](https://github.com/lindell/multi-gitter/issues/70)) ([02a2654](https://github.com/lindell/multi-gitter/commit/02a26542063d5ccaf029d087745082d40ce446af))
* release 0.20.1 ([#71](https://github.com/lindell/multi-gitter/issues/71)) ([f930b9e](https://github.com/lindell/multi-gitter/commit/f930b9e39cafd5dbfc2bf3329bbd3f229688f57c))
* release 0.20.2 ([#72](https://github.com/lindell/multi-gitter/issues/72)) ([e831ced](https://github.com/lindell/multi-gitter/commit/e831ced6095796edae16dadc229436cb9229e08b))
* release 0.20.3 ([#74](https://github.com/lindell/multi-gitter/issues/74)) ([8686185](https://github.com/lindell/multi-gitter/commit/8686185b3a8c16a23baf949e59a9b480f7bacd40))
* release 0.20.4 ([#75](https://github.com/lindell/multi-gitter/issues/75)) ([6f36672](https://github.com/lindell/multi-gitter/commit/6f36672fd5e797bcb3daf3ba4fe052cf6c583dd8))
* release 0.20.5 ([#76](https://github.com/lindell/multi-gitter/issues/76)) ([8f42eae](https://github.com/lindell/multi-gitter/commit/8f42eae65c373585310589c6296a7dcaaa899d03))
* release 0.21.0 ([#77](https://github.com/lindell/multi-gitter/issues/77)) ([736943f](https://github.com/lindell/multi-gitter/commit/736943f061eff704e460178ebbe49a6d1930bc92))
* release 0.21.1 ([#78](https://github.com/lindell/multi-gitter/issues/78)) ([ad04832](https://github.com/lindell/multi-gitter/commit/ad04832ddf772e0d0d84c61c3b1d89f59d81b709))
* release 0.22.0 ([#81](https://github.com/lindell/multi-gitter/issues/81)) ([f3b6dfb](https://github.com/lindell/multi-gitter/commit/f3b6dfbc3d9f39569858e6c0e5b97aa4904c0f3a))
* release 0.22.1 ([#83](https://github.com/lindell/multi-gitter/issues/83)) ([75a995c](https://github.com/lindell/multi-gitter/commit/75a995c05372d977a5fe71a9d2960e639cbc6f10))
* release 0.23.0 ([#85](https://github.com/lindell/multi-gitter/issues/85)) ([010f6ba](https://github.com/lindell/multi-gitter/commit/010f6ba07face772e6891c0161440725f975aef2))
* release 0.23.1 ([#86](https://github.com/lindell/multi-gitter/issues/86)) ([3c16c98](https://github.com/lindell/multi-gitter/commit/3c16c9876a8653dd0e0562b5eb9f47a983e81bc8))
* release 0.24.0 ([#88](https://github.com/lindell/multi-gitter/issues/88)) ([98fb52f](https://github.com/lindell/multi-gitter/commit/98fb52fdff55c3bead87c76f9c67cd05ff8b39b7))
* release 0.24.1 ([#90](https://github.com/lindell/multi-gitter/issues/90)) ([cd6c428](https://github.com/lindell/multi-gitter/commit/cd6c428587c5cc9c543714f0e107461ddc84ec1d))
* release 0.24.2 ([#94](https://github.com/lindell/multi-gitter/issues/94)) ([10b96fd](https://github.com/lindell/multi-gitter/commit/10b96fd10e7e6653e94b6a436b445f6ba3ec9b2d))
* release 0.25.0 ([#106](https://github.com/lindell/multi-gitter/issues/106)) ([9c64c99](https://github.com/lindell/multi-gitter/commit/9c64c994bb68bb5e56914fc58a90ddc1c7997e35))
* release 0.25.1 ([#111](https://github.com/lindell/multi-gitter/issues/111)) ([512d6c1](https://github.com/lindell/multi-gitter/commit/512d6c1d879c6841dd531edcddd9238ebbc567d6))
* release 0.25.2 ([#112](https://github.com/lindell/multi-gitter/issues/112)) ([9e1fbed](https://github.com/lindell/multi-gitter/commit/9e1fbedfaefccda4037bd8d658ebfbe251759f3c))
* release 0.25.3 ([#115](https://github.com/lindell/multi-gitter/issues/115)) ([e71ccb6](https://github.com/lindell/multi-gitter/commit/e71ccb6898059849307e67c5c261bacead9041da))
* release 0.25.4 ([#120](https://github.com/lindell/multi-gitter/issues/120)) ([c4906c8](https://github.com/lindell/multi-gitter/commit/c4906c864beaa43e3678d0322e2eae5c3c75cd06))
* release 0.25.5 ([#121](https://github.com/lindell/multi-gitter/issues/121)) ([874196d](https://github.com/lindell/multi-gitter/commit/874196d153d60d389eb2ee4c2e4481c86e1b96f1))
* release 0.25.6 ([#126](https://github.com/lindell/multi-gitter/issues/126)) ([c092d21](https://github.com/lindell/multi-gitter/commit/c092d219a5d45c0e0db02f8613cc53e465cc9834))
* release 0.26.0 ([#132](https://github.com/lindell/multi-gitter/issues/132)) ([d2bc434](https://github.com/lindell/multi-gitter/commit/d2bc434d1c1ebf1517c8a19aa992098f0411f2f2))
* release 0.26.1 ([#136](https://github.com/lindell/multi-gitter/issues/136)) ([776b929](https://github.com/lindell/multi-gitter/commit/776b9296c1c7d1054ff0e58c1d3650fd738791d4))
* release 0.27.0 ([#139](https://github.com/lindell/multi-gitter/issues/139)) ([37cd65d](https://github.com/lindell/multi-gitter/commit/37cd65d8905598154d623bb7019494ae6934df36))
* release 0.28.0 ([#141](https://github.com/lindell/multi-gitter/issues/141)) ([fbf3069](https://github.com/lindell/multi-gitter/commit/fbf30690105ebc46c3363c93c9e882206c3c015f))
* release 0.29.0 ([#152](https://github.com/lindell/multi-gitter/issues/152)) ([2b16fb5](https://github.com/lindell/multi-gitter/commit/2b16fb59d3d1af42c9897d8ffb88f35b1fd3f2f9))
* release 0.29.1 ([#158](https://github.com/lindell/multi-gitter/issues/158)) ([d298250](https://github.com/lindell/multi-gitter/commit/d298250191ac85970652117d497065b9bd14a169))
* release 0.29.2 ([#160](https://github.com/lindell/multi-gitter/issues/160)) ([df21838](https://github.com/lindell/multi-gitter/commit/df2183832d839db91bead86d157e52237d55ee51))
* release 0.30.0 ([81bbf3c](https://github.com/lindell/multi-gitter/commit/81bbf3c16fa288b72f24f7768ebe8d8825bd8a09))
* release 0.31.0 ([#168](https://github.com/lindell/multi-gitter/issues/168)) ([163a655](https://github.com/lindell/multi-gitter/commit/163a655499717168a29d36a0e368e3970cb1c9c4))
* release 0.31.1 ([#172](https://github.com/lindell/multi-gitter/issues/172)) ([89e8938](https://github.com/lindell/multi-gitter/commit/89e89382655328fcde0154b601f16cbeb4f7b958))
* release 0.32.0 ([#175](https://github.com/lindell/multi-gitter/issues/175)) ([3beecef](https://github.com/lindell/multi-gitter/commit/3beecef5e7d14264198860b1713ec730102cef7c))
* release 0.33.0 ([#181](https://github.com/lindell/multi-gitter/issues/181)) ([e3ae300](https://github.com/lindell/multi-gitter/commit/e3ae30091ed7676e7490fc4c15133ae83a85569e))
* release 0.33.1 ([#182](https://github.com/lindell/multi-gitter/issues/182)) ([9ad0269](https://github.com/lindell/multi-gitter/commit/9ad0269cfd94e31105861168a390e54db51bebdb))
* release 0.33.2 ([#186](https://github.com/lindell/multi-gitter/issues/186)) ([6212369](https://github.com/lindell/multi-gitter/commit/6212369ddd39a1cb0ae041ec06752fbb762c12e1))
* release 0.33.3 ([#190](https://github.com/lindell/multi-gitter/issues/190)) ([df64360](https://github.com/lindell/multi-gitter/commit/df6436026dc5bba88b71052f18e254d5e7e542f4))
* release 0.34.0 ([#199](https://github.com/lindell/multi-gitter/issues/199)) ([66c0060](https://github.com/lindell/multi-gitter/commit/66c0060c8c1891f89ebbfa8f0208eec56daf5138))
* release 0.35.0 ([#201](https://github.com/lindell/multi-gitter/issues/201)) ([fbbdea2](https://github.com/lindell/multi-gitter/commit/fbbdea281501ffbb9fc24f65a4c4298112a1ccf8))
* release 0.36.0 ([#203](https://github.com/lindell/multi-gitter/issues/203)) ([97aa91f](https://github.com/lindell/multi-gitter/commit/97aa91f8b531f477f542cffa9f12dbae2a030a72))
* release 0.36.1 ([#208](https://github.com/lindell/multi-gitter/issues/208)) ([760ee3c](https://github.com/lindell/multi-gitter/commit/760ee3ca681467cde153f2e11fd05a99d7120f7c))
* release 0.37.0 ([#213](https://github.com/lindell/multi-gitter/issues/213)) ([8e30ba9](https://github.com/lindell/multi-gitter/commit/8e30ba9cea00cf194f6db0412e11f6957fe8c5e8))
* release 0.38.0 ([#216](https://github.com/lindell/multi-gitter/issues/216)) ([cd32051](https://github.com/lindell/multi-gitter/commit/cd320515ad7af0e492d0e7095f98a70ee02ac990))
* release 0.38.1 ([#218](https://github.com/lindell/multi-gitter/issues/218)) ([e7b38e6](https://github.com/lindell/multi-gitter/commit/e7b38e6529e564383f382a41fd1972220108eecf))
* release 0.38.2 ([#219](https://github.com/lindell/multi-gitter/issues/219)) ([e54dfbc](https://github.com/lindell/multi-gitter/commit/e54dfbccf3da56ebb22c5d37b4166b73329680e8))
* release 0.38.3 ([#220](https://github.com/lindell/multi-gitter/issues/220)) ([12d9f50](https://github.com/lindell/multi-gitter/commit/12d9f50d7ef45ea504a63afedd4a13830ae4dffe))
* release 0.39.0 ([#231](https://github.com/lindell/multi-gitter/issues/231)) ([5f3acc0](https://github.com/lindell/multi-gitter/commit/5f3acc0633736be9c6bcee55a2931007682c3109))
* release 0.40.0 ([#233](https://github.com/lindell/multi-gitter/issues/233)) ([bf6ca70](https://github.com/lindell/multi-gitter/commit/bf6ca70784c0a0493712c8709651fdf7dd629fe0))
* release 0.40.1 ([#238](https://github.com/lindell/multi-gitter/issues/238)) ([f99b8ed](https://github.com/lindell/multi-gitter/commit/f99b8ede8731c6e88d9c6be99c10d3389475e7ec))
* release 0.41.0 ([#241](https://github.com/lindell/multi-gitter/issues/241)) ([2cdf5b3](https://github.com/lindell/multi-gitter/commit/2cdf5b3bf52c268e4ff450e8bfb1c5bf9a23835f))
* release 0.42.0 ([#248](https://github.com/lindell/multi-gitter/issues/248)) ([d8f9531](https://github.com/lindell/multi-gitter/commit/d8f9531e51445feb68bf5e8e064b7a1e83285a65))
* release 0.42.1 ([#251](https://github.com/lindell/multi-gitter/issues/251)) ([035bd3b](https://github.com/lindell/multi-gitter/commit/035bd3bec067f2cdd1e66cc560862fa4fbb008a2))
* release 0.42.2 ([#271](https://github.com/lindell/multi-gitter/issues/271)) ([ad5bb77](https://github.com/lindell/multi-gitter/commit/ad5bb777a6b2817b69364b4e1d72ef2b0ac0295c))
* release 0.42.3 ([#274](https://github.com/lindell/multi-gitter/issues/274)) ([e4bc83a](https://github.com/lindell/multi-gitter/commit/e4bc83aa64aa709ded9cfd04e8e40d8ebe65e988))
* release 0.42.4 ([#277](https://github.com/lindell/multi-gitter/issues/277)) ([28105ed](https://github.com/lindell/multi-gitter/commit/28105ed0354ce4c9687bb54b19bfb76b3dd86809))
* release 0.43.0 ([#293](https://github.com/lindell/multi-gitter/issues/293)) ([c12fb13](https://github.com/lindell/multi-gitter/commit/c12fb13b4889e174c684f0fc1ef46b52b5462a67))
* release 0.43.1 ([#297](https://github.com/lindell/multi-gitter/issues/297)) ([49e6477](https://github.com/lindell/multi-gitter/commit/49e6477619afc47b7b35e94ad53200349153f4d7))
* release 0.43.2 ([#298](https://github.com/lindell/multi-gitter/issues/298)) ([3436db3](https://github.com/lindell/multi-gitter/commit/3436db3c806343d2d83f58e65b35c30c906ef5d8))
* release 0.43.3 ([#304](https://github.com/lindell/multi-gitter/issues/304)) ([4925583](https://github.com/lindell/multi-gitter/commit/4925583322976f7f399b6276881bab50e2f5d544))
* release 0.44.0 ([#324](https://github.com/lindell/multi-gitter/issues/324)) ([6ec9e89](https://github.com/lindell/multi-gitter/commit/6ec9e89d6f2957476633edc7db19661e611860f6))
* release 0.44.1 ([#326](https://github.com/lindell/multi-gitter/issues/326)) ([fd839bb](https://github.com/lindell/multi-gitter/commit/fd839bb33d020298a1a75ffd69b42b08f9ed90ce))
* release 0.44.2 ([#335](https://github.com/lindell/multi-gitter/issues/335)) ([d08e726](https://github.com/lindell/multi-gitter/commit/d08e7264018d17624b1d00822470a73e3d18b3ca))
* release 0.45.0 ([#338](https://github.com/lindell/multi-gitter/issues/338)) ([1f2b956](https://github.com/lindell/multi-gitter/commit/1f2b9568b29055edb3ec05631b509504842688d9))
* release 0.46.0 ([#344](https://github.com/lindell/multi-gitter/issues/344)) ([451167c](https://github.com/lindell/multi-gitter/commit/451167c22f458bfc242f5ed80e064ce5bbe1344a))
* release 0.47.0 ([#349](https://github.com/lindell/multi-gitter/issues/349)) ([c5a6730](https://github.com/lindell/multi-gitter/commit/c5a67305f810c87ec8a2973aedd9e2c002d29314))
* removed brew download strategy ([6c35be5](https://github.com/lindell/multi-gitter/commit/6c35be50c14dff82363931bd12ade4de204103c1))
* update to Go 1.21 ([#402](https://github.com/lindell/multi-gitter/issues/402)) ([d454e5b](https://github.com/lindell/multi-gitter/commit/d454e5b17e90fa2b2d36fa60d15b51df3e7701fa))
* updated to go 1.16 ([c8fa961](https://github.com/lindell/multi-gitter/commit/c8fa96154f2925ea724bf9ff0a74027dfc0a9286))
* updated to go modules to 1.18 ([#272](https://github.com/lindell/multi-gitter/issues/272)) ([d70a038](https://github.com/lindell/multi-gitter/commit/d70a03855a0a3da5b70e0ee85f53326cc4c137e6))

## [0.49.1](https://github.com/lindell/multi-gitter/compare/v0.49.0...v0.49.1) (2024-01-05)


### Bug Fixes

* cmdgit BranchExist should match full branch name ([#433](https://github.com/lindell/multi-gitter/issues/433)) ([c7a6e69](https://github.com/lindell/multi-gitter/commit/c7a6e69aafb87efafa25248c73d7afa1a8ee76ac))


### Dependencies

* update golang.org/x/exp digest to 02704c9 ([f7e273e](https://github.com/lindell/multi-gitter/commit/f7e273e6ebf4ee52d16e76e8472de87eb9c3a52e))
* update golang.org/x/exp digest to 6522937 ([cab859f](https://github.com/lindell/multi-gitter/commit/cab859fdb40c181e9b1dc58448290eec800f6aae))
* update module code.gitea.io/sdk/gitea to v0.17.0 ([3949d2f](https://github.com/lindell/multi-gitter/commit/3949d2f4a3d7f315ca48a6c8c18910a1bcce3443))
* update module github.com/go-git/go-git/v5 to v5.10.1 ([40f62a1](https://github.com/lindell/multi-gitter/commit/40f62a11cacfd429e232451aeab9e11825de9ce3))
* update module github.com/go-git/go-git/v5 to v5.11.0 [security] ([5fef1cc](https://github.com/lindell/multi-gitter/commit/5fef1ccc936c2d6f31a8d4ae7ac0ebed6409f2c1))
* update module github.com/google/go-github/v56 to v57 ([#445](https://github.com/lindell/multi-gitter/issues/445)) ([28a3279](https://github.com/lindell/multi-gitter/commit/28a32793edd1c63bb66b2008ef333f2f7b31c07b))
* update module github.com/spf13/cobra to v1.8.0 ([f976ac2](https://github.com/lindell/multi-gitter/commit/f976ac2d91254e92f5e937f9e0fad99802b8ce79))
* update module github.com/spf13/viper to v1.18.2 ([29a205e](https://github.com/lindell/multi-gitter/commit/29a205e9db876a0baf6091415f47085bb2876b5e))
* update module github.com/xanzy/go-gitlab to v0.94.0 ([38666d9](https://github.com/lindell/multi-gitter/commit/38666d9b5ae79783ea161a20ad566480e5c1b23a))
* update module github.com/xanzy/go-gitlab to v0.95.2 ([#446](https://github.com/lindell/multi-gitter/issues/446)) ([61f0110](https://github.com/lindell/multi-gitter/commit/61f0110512e4ff464315e3791e6393e62830f112))
* update module golang.org/x/oauth2 to v0.15.0 ([3874ccb](https://github.com/lindell/multi-gitter/commit/3874ccb93ab76c3a2724358f825880b750a5dbea))

## [0.49.0](https://github.com/lindell/multi-gitter/compare/v0.48.2...v0.49.0) (2023-11-09)


### Features

* add repo-include and repo-exclude options ([#426](https://github.com/lindell/multi-gitter/issues/426)) ([faf0092](https://github.com/lindell/multi-gitter/commit/faf0092a4c5411420f721973595588f6303253da))
* **github:** use code search to select repositories ([#399](https://github.com/lindell/multi-gitter/issues/399)) ([d8b75f4](https://github.com/lindell/multi-gitter/commit/d8b75f4ed173c2389dcb3dd89a4213fd8ce78f42))

## [0.48.2](https://github.com/lindell/multi-gitter/compare/v0.48.1...v0.48.2) (2023-11-04)


### Bug Fixes

* **typo:** fixed panicked typo ([#423](https://github.com/lindell/multi-gitter/issues/423)) ([9b5774a](https://github.com/lindell/multi-gitter/commit/9b5774ad8e95f849d9af0b76a58183c7bbdf0462))

## [0.48.1](https://github.com/lindell/multi-gitter/compare/v0.48.0...v0.48.1) (2023-11-02)


### Bug Fixes

* mitigated script injection in test install action ([#404](https://github.com/lindell/multi-gitter/issues/404)) ([f77ef91](https://github.com/lindell/multi-gitter/commit/f77ef91416f3b13080372bb00831824eea37c130))


### Dependencies

* update github.com/gfleury/go-bitbucket-v1 digest to 6e30c57 ([#382](https://github.com/lindell/multi-gitter/issues/382)) ([c201f42](https://github.com/lindell/multi-gitter/commit/c201f425b4fbde46c3d27fa380ab099e0a6f8608))
* update module github.com/go-git/go-git/v5 to v5.10.0 ([077d637](https://github.com/lindell/multi-gitter/commit/077d637789ce7943212ca01827dd06fe314ad75d))
* update module github.com/google/go-github/v55 to v56 ([#412](https://github.com/lindell/multi-gitter/issues/412)) ([946afac](https://github.com/lindell/multi-gitter/commit/946afac76f4a1884823da591d04236daaf3fb795))
* update module github.com/spf13/viper to v1.17.0 ([cfe3409](https://github.com/lindell/multi-gitter/commit/cfe3409acee757b78c542e106436e7ca177cf82b))
* update module github.com/xanzy/go-gitlab to v0.93.2 ([76fd29b](https://github.com/lindell/multi-gitter/commit/76fd29b07a7561223adb474bf4fd5eedae6a4298))
* update module golang.org/x/oauth2 to v0.13.0 ([3b3e9c0](https://github.com/lindell/multi-gitter/commit/3b3e9c009d1d4a3634ce8decb75625f2ed8e4d0d))

## [0.48.0](https://github.com/lindell/multi-gitter/compare/v0.47.4...v0.48.0) (2023-10-27)


### Features

* added --plain-output flag ([#378](https://github.com/lindell/multi-gitter/issues/378)) ([356f471](https://github.com/lindell/multi-gitter/commit/356f471f720e6efe0146bfabcb99b2c214d9c4db))
* use repo search to select repositories ([#397](https://github.com/lindell/multi-gitter/issues/397)) ([0f8c2dc](https://github.com/lindell/multi-gitter/commit/0f8c2dcceeb9b635c7405e878df608318b518de4))


### Bug Fixes

* **gitlab:** make sure entire project path is shown ([#392](https://github.com/lindell/multi-gitter/issues/392)) ([280dfc7](https://github.com/lindell/multi-gitter/commit/280dfc7b6ed5c8a5345f2207e58e81202a451e33))


### Dependencies

* update module code.gitea.io/sdk/gitea to v0.16.0 ([ff04448](https://github.com/lindell/multi-gitter/commit/ff0444885bfbf657b4813efdb53986af67b0ac58))
* update module github.com/go-git/go-git/v5 to v5.9.0 ([a954c7c](https://github.com/lindell/multi-gitter/commit/a954c7c0c20cb0a930203c92143456c1571d2e2f))
* update module github.com/google/go-github/v53 to v55 ([#390](https://github.com/lindell/multi-gitter/issues/390)) ([d92cb29](https://github.com/lindell/multi-gitter/commit/d92cb29cc96e8a78d19ccdd1816643510169f6c9))
* update module github.com/xanzy/go-gitlab to v0.92.1 ([19c2925](https://github.com/lindell/multi-gitter/commit/19c292523f3442766f09d5811ec9b38ecb26c9bf))
* update module github.com/xanzy/go-gitlab to v0.92.3 ([3e6cd75](https://github.com/lindell/multi-gitter/commit/3e6cd75ad0923ad535f2125f1c94cba6b9072f05))
* update module golang.org/x/oauth2 to v0.12.0 ([b9c258a](https://github.com/lindell/multi-gitter/commit/b9c258a2b2924ebafa37864159154dc7ed09449c))


### Miscellaneous

* update to Go 1.21 ([#402](https://github.com/lindell/multi-gitter/issues/402)) ([d454e5b](https://github.com/lindell/multi-gitter/commit/d454e5b17e90fa2b2d36fa60d15b51df3e7701fa))

## [0.47.4](https://github.com/lindell/multi-gitter/compare/v0.47.3...v0.47.4) (2023-09-12)


### Bug Fixes

* don't panic when a PR branch is deleted ([#385](https://github.com/lindell/multi-gitter/issues/385)) ([fca98e1](https://github.com/lindell/multi-gitter/commit/fca98e134a6640d79aae6bffbe154d42b3393d40))


### Dependencies

* update module github.com/go-git/go-git/v5 to v5.8.1 ([8de2001](https://github.com/lindell/multi-gitter/commit/8de20010a3c8c7ac5eecb6c366975d164e9a8409))
* update module github.com/xanzy/go-gitlab to v0.88.0 ([c0f0bb4](https://github.com/lindell/multi-gitter/commit/c0f0bb4ba45c9834fe67a44774a96a984a16913c))
* update module github.com/xanzy/go-gitlab to v0.90.0 ([fd66cd2](https://github.com/lindell/multi-gitter/commit/fd66cd200c08b911fb19dae40dafa6c91170a269))
* update module golang.org/x/oauth2 to v0.10.0 ([a82fc38](https://github.com/lindell/multi-gitter/commit/a82fc380c95ab1bec3a20a2da0465444b8d91cd2))
* update module golang.org/x/oauth2 to v0.11.0 ([7820c22](https://github.com/lindell/multi-gitter/commit/7820c228f69386a12df3ddc4c908bb1940b9aa2e))

## [0.47.3](https://github.com/lindell/multi-gitter/compare/v0.47.2...v0.47.3) (2023-07-16)


### Bug Fixes

* changed the order of version and os in the released archives ([656e9de](https://github.com/lindell/multi-gitter/commit/656e9de09c0bdb9064d47bd2b6adeb23f971eea3))

## [0.47.2](https://github.com/lindell/multi-gitter/compare/v0.47.1...v0.47.2) (2023-07-16)


### Bug Fixes

* fixed archive release names ([1462b90](https://github.com/lindell/multi-gitter/commit/1462b901c373c0a1f04c73f38784fac01cb65189))

## [0.47.1](https://github.com/lindell/multi-gitter/compare/v0.47.0...v0.47.1) (2023-07-15)


### Bug Fixes

* ensure pull request exists during conflict resolution ([#369](https://github.com/lindell/multi-gitter/issues/369)) ([2b7166a](https://github.com/lindell/multi-gitter/commit/2b7166a203ddd7d939bc28a90ee26e08adb5a771))


### Dependencies

* update github.com/gfleury/go-bitbucket-v1 digest to 8d7be58 ([1d8a19f](https://github.com/lindell/multi-gitter/commit/1d8a19f8316dd8874a0755189902dce24acd09c5))
* update module github.com/go-git/go-git/v5 to v5.7.0 ([d407eb9](https://github.com/lindell/multi-gitter/commit/d407eb9ef7bb06ca54b0d28ad992c8832ecb6c41))
* update module github.com/google/go-github/v50 to v50.2.0 ([#359](https://github.com/lindell/multi-gitter/issues/359)) ([481c5da](https://github.com/lindell/multi-gitter/commit/481c5da7223eb798bd04b21eb839aa6aa563a763))
* update module github.com/google/go-github/v50 to v53 ([#362](https://github.com/lindell/multi-gitter/issues/362)) ([ff85919](https://github.com/lindell/multi-gitter/commit/ff85919b1a264ce042e9e2bbca24e356be15e623))
* update module github.com/sirupsen/logrus to v1.9.3 ([3f0d4c7](https://github.com/lindell/multi-gitter/commit/3f0d4c707636a595404839deecc6be9204cf0407))
* update module github.com/spf13/cobra to v1.7.0 ([#360](https://github.com/lindell/multi-gitter/issues/360)) ([212a9b6](https://github.com/lindell/multi-gitter/commit/212a9b686db37c58205d7f2e20778ec592e15388))
* update module github.com/spf13/viper to v1.16.0 ([#364](https://github.com/lindell/multi-gitter/issues/364)) ([af41308](https://github.com/lindell/multi-gitter/commit/af4130863a62217bcf2c5f4fda163b521112a7a0))
* update module github.com/stretchr/testify to v1.8.4 ([#361](https://github.com/lindell/multi-gitter/issues/361)) ([2686055](https://github.com/lindell/multi-gitter/commit/2686055d74f2f2af416ff6ed05eb4f47647648be))
* update module github.com/xanzy/go-gitlab to v0.86.0 ([#365](https://github.com/lindell/multi-gitter/issues/365)) ([79b9a0e](https://github.com/lindell/multi-gitter/commit/79b9a0e57a45ca5c55af310d1a619018166d2362))
* update module golang.org/x/oauth2 to v0.9.0 ([#363](https://github.com/lindell/multi-gitter/issues/363)) ([58030bd](https://github.com/lindell/multi-gitter/commit/58030bd778db562dedcb6cf48eef9dbb2a8039d9))

## [0.47.0](https://www.github.com/lindell/multi-gitter/compare/v0.46.0...v0.47.0) (2023-05-02)


### Features

* **github:** added option to add team reviewers ([#351](https://www.github.com/lindell/multi-gitter/issues/351)) ([bfe05b9](https://www.github.com/lindell/multi-gitter/commit/bfe05b9b5c53307a8c429278b6491bafe57a2f26))
* OS and Arch info added to the version command ([#348](https://www.github.com/lindell/multi-gitter/issues/348)) ([12c0422](https://www.github.com/lindell/multi-gitter/commit/12c04221fa678bf230a50f4a4386aebe4cfa593f))


### Bug Fixes

* update error message to fix grammar ([#345](https://www.github.com/lindell/multi-gitter/issues/345)) ([5a4c990](https://www.github.com/lindell/multi-gitter/commit/5a4c990b7f9325d1a5c5b4ee619e792478517754))


### Dependencies

* update module github.com/go-git/go-git/v5 to v5.6.1 ([1fddf2e](https://www.github.com/lindell/multi-gitter/commit/1fddf2e26cd3cb9f201325decc356d26a1f1444e))

## [0.46.0](https://www.github.com/lindell/multi-gitter/compare/v0.45.0...v0.46.0) (2023-04-09)


### Features

* option to skip repositories that are forks ([#341](https://www.github.com/lindell/multi-gitter/issues/341)) ([941731b](https://www.github.com/lindell/multi-gitter/commit/941731bfc0a9a89b2abef18286a7a3b06ab5d1db))

## [0.45.0](https://www.github.com/lindell/multi-gitter/compare/v0.44.2...v0.45.0) (2023-04-01)


### Features

* set DRY_RUN when --dry-run is used ([#337](https://www.github.com/lindell/multi-gitter/issues/337)) ([e4390ee](https://www.github.com/lindell/multi-gitter/commit/e4390ee8eddc8a37ea73fd5f29d3e9151221900e))

### [0.44.2](https://www.github.com/lindell/multi-gitter/compare/v0.44.1...v0.44.2) (2023-03-25)


### Bug Fixes

* **github:** allow GitHub app tokens to be used ([#334](https://www.github.com/lindell/multi-gitter/issues/334)) ([8d86544](https://www.github.com/lindell/multi-gitter/commit/8d865447765d70c6bdde393ef3ee450355b61ae0))

### [0.44.1](https://www.github.com/lindell/multi-gitter/compare/v0.44.0...v0.44.1) (2023-03-16)


### Bug Fixes

* **github:** chunk get pull requests ([#330](https://www.github.com/lindell/multi-gitter/issues/330)) ([488cd63](https://www.github.com/lindell/multi-gitter/commit/488cd6339bd659e88cd89d207ae1d5eb5e40b988))


### Dependencies

* bump golang.org/x/net from 0.3.0 to 0.7.0 ([#325](https://www.github.com/lindell/multi-gitter/issues/325)) ([9d0ea43](https://www.github.com/lindell/multi-gitter/commit/9d0ea43c5ad46a76d8a596e0d33082ce51579aa6))
* update module github.com/google/go-github/v50 to v50.1.0 ([366f201](https://www.github.com/lindell/multi-gitter/commit/366f201f77bc46aaeedfa34fd6ed40dc99a007f6))
* update module github.com/stretchr/testify to v1.8.2 ([4f3ce98](https://www.github.com/lindell/multi-gitter/commit/4f3ce98307815934149fbef6da211280bbd581dd))
* update module golang.org/x/oauth2 to v0.5.0 ([#327](https://www.github.com/lindell/multi-gitter/issues/327)) ([fefceee](https://www.github.com/lindell/multi-gitter/commit/fefceeea27305bb94702b53ad5cfa0bfb6cf973f))

## [0.44.0](https://www.github.com/lindell/multi-gitter/compare/v0.43.3...v0.44.0) (2023-02-13)


### Features

* filter repositories using topics ([#320](https://www.github.com/lindell/multi-gitter/issues/320)) ([d3c5403](https://www.github.com/lindell/multi-gitter/commit/d3c54034d56f7826dfa7dbc13851b152334fbf58))


### Bug Fixes

* input description for platform args ([#321](https://www.github.com/lindell/multi-gitter/issues/321)) ([6816c86](https://www.github.com/lindell/multi-gitter/commit/6816c86ea534112ef55c24036940f2c44b50e1ab))


### Dependencies

* update module github.com/go-git/go-git/v5 to v5.5.1 ([#312](https://www.github.com/lindell/multi-gitter/issues/312)) ([7c9136c](https://www.github.com/lindell/multi-gitter/commit/7c9136c7427fcb1d103c0f1e004059f7a793f25d))
* update module github.com/go-git/go-git/v5 to v5.5.2 ([83032fb](https://www.github.com/lindell/multi-gitter/commit/83032fb7e4cd197c703c7eb8e272938a315c5ef5))
* update module github.com/google/go-github/v47 to v48 ([#311](https://www.github.com/lindell/multi-gitter/issues/311)) ([c3ffe09](https://www.github.com/lindell/multi-gitter/commit/c3ffe095a56d450ea2d5aa81384a4578cc069a12))
* update module github.com/google/go-github/v48 to v50 ([#319](https://www.github.com/lindell/multi-gitter/issues/319)) ([6f1fe90](https://www.github.com/lindell/multi-gitter/commit/6f1fe9009131d1d0326c7f63cf6d07b2cd1983f3))
* update module github.com/xanzy/go-gitlab to v0.77.0 ([5e39817](https://www.github.com/lindell/multi-gitter/commit/5e3981748c08e5ee4d63b0a1de1b8b4df50b84fc))
* update module golang.org/x/oauth2 to v0.3.0 ([4ae8184](https://www.github.com/lindell/multi-gitter/commit/4ae81844ffe59a397753f046ae9cb87e2c1bac48))

### [0.43.3](https://www.github.com/lindell/multi-gitter/compare/v0.43.2...v0.43.3) (2022-12-30)


### Bug Fixes

* allow individual merge failures without aborting all merges ([#310](https://www.github.com/lindell/multi-gitter/issues/310)) ([e067502](https://www.github.com/lindell/multi-gitter/commit/e067502f2366e86ce6b979130233d41017f40caf))
* typos in docs and code ([#303](https://www.github.com/lindell/multi-gitter/issues/303)) ([45ddb60](https://www.github.com/lindell/multi-gitter/commit/45ddb60cd438e9251b670c0884cc851c29f09d8f))


### Dependencies

* update module github.com/spf13/viper to v1.14.0 ([294eced](https://www.github.com/lindell/multi-gitter/commit/294eced3db7de848f4c721a651e4833007bb3281))
* update module github.com/xanzy/go-gitlab to v0.76.0 ([0ce73dc](https://www.github.com/lindell/multi-gitter/commit/0ce73dc36d03ea6fa94a7203c0825028e1a97c0b))

### [0.43.2](https://www.github.com/lindell/multi-gitter/compare/v0.43.1...v0.43.2) (2022-11-10)


### Bug Fixes

* **github:** made sure GraphQL requests, with non-GraphQL errors is reported properly ([#301](https://www.github.com/lindell/multi-gitter/issues/301)) ([d7e1fda](https://www.github.com/lindell/multi-gitter/commit/d7e1fda392d1046fa93c2f0304295c0fd7872292))
* **github:** use other format when cloning with token to support more token formats ([#302](https://www.github.com/lindell/multi-gitter/issues/302)) ([a74cc60](https://www.github.com/lindell/multi-gitter/commit/a74cc60d8257e69a99f224150f45400223b5911b))


### Dependencies

* update module github.com/spf13/cobra to v1.6.1 ([cfc2861](https://www.github.com/lindell/multi-gitter/commit/cfc2861089ac4e8697202c6bcf5daf01021fd8bf))
* update module github.com/stretchr/testify to v1.8.1 ([e5158dc](https://www.github.com/lindell/multi-gitter/commit/e5158dc7b3536ec681d090ae18d5db88feb19fe0))

### [0.43.1](https://www.github.com/lindell/multi-gitter/compare/v0.43.0...v0.43.1) (2022-10-25)


### Bug Fixes

* **github:** add retry mechanism to all requests ([#289](https://www.github.com/lindell/multi-gitter/issues/289)) ([89a642c](https://www.github.com/lindell/multi-gitter/commit/89a642c551fa032d4389883acd26146109ebfa99))
* added support for cancellation of git remote commands ([#296](https://www.github.com/lindell/multi-gitter/issues/296)) ([83bfbc7](https://www.github.com/lindell/multi-gitter/commit/83bfbc764a920ffbbedc5a80028b889391360dfd))

## [0.43.0](https://www.github.com/lindell/multi-gitter/compare/v0.42.4...v0.43.0) (2022-10-18)


### Features

* added option to add labels to creates pull requests ([#292](https://www.github.com/lindell/multi-gitter/issues/292)) ([fccf678](https://www.github.com/lindell/multi-gitter/commit/fccf678c384ff01cd0247c35860acf0d257e41a7))

### [0.42.4](https://www.github.com/lindell/multi-gitter/compare/v0.42.3...v0.42.4) (2022-10-17)


### Bug Fixes

* **github:** added missing closed pr status ([#290](https://www.github.com/lindell/multi-gitter/issues/290)) ([9e3644a](https://www.github.com/lindell/multi-gitter/commit/9e3644a899197470985707112238f9990df99564))
* better error message when the same feature and base branch is used ([#281](https://www.github.com/lindell/multi-gitter/issues/281)) ([ff98e8f](https://www.github.com/lindell/multi-gitter/commit/ff98e8fff43a3ccc14feb2f24f6c54c134dbc21c))
* retry when encounting the GitHub rate limit ([#280](https://www.github.com/lindell/multi-gitter/issues/280)) ([008a26a](https://www.github.com/lindell/multi-gitter/commit/008a26ae0182c8e69cc0679a0d7ee776415c1d62))


### Dependencies

* update golang.org/x/oauth2 digest to f213421 ([dfd1837](https://www.github.com/lindell/multi-gitter/commit/dfd18376c44f4b9a36628142fc52277b0aa48ec0))
* update module github.com/spf13/viper to v1.13.0 ([099a9c9](https://www.github.com/lindell/multi-gitter/commit/099a9c9f3e5341b207ec298356319303df67a31f))
* update module github.com/xanzy/go-gitlab to v0.73.1 ([badc233](https://www.github.com/lindell/multi-gitter/commit/badc23353c3f611031fa70483b6892cf8f2e0b99))
* update module go to 1.19 ([2fdbcde](https://www.github.com/lindell/multi-gitter/commit/2fdbcde5f6160ece5b8fac38d65761278b9f0f6a))
* update module go-github to v47 ([#286](https://www.github.com/lindell/multi-gitter/issues/286)) ([68fb1e2](https://www.github.com/lindell/multi-gitter/commit/68fb1e28af23bbdd8c77323e572a3c43d06a4445))

### [0.42.3](https://www.github.com/lindell/multi-gitter/compare/v0.42.2...v0.42.3) (2022-08-12)


### Bug Fixes

* deleted files are now detected with `git-type: go` ([#273](https://www.github.com/lindell/multi-gitter/issues/273)) ([b21509c](https://www.github.com/lindell/multi-gitter/commit/b21509c3b392fbe681c1cc5315ccdc63942abbc2))

### [0.42.2](https://www.github.com/lindell/multi-gitter/compare/v0.42.1...v0.42.2) (2022-08-07)


### Bug Fixes

* **github:** correctly map merged prs ([#268](https://www.github.com/lindell/multi-gitter/issues/268)) ([0474040](https://www.github.com/lindell/multi-gitter/commit/0474040c86d7c213987fb5dfc503d6c46288cb3a))


### Miscellaneous

* updated to go modules to 1.18 ([#272](https://www.github.com/lindell/multi-gitter/issues/272)) ([d70a038](https://www.github.com/lindell/multi-gitter/commit/d70a03855a0a3da5b70e0ee85f53326cc4c137e6))

### [0.42.1](https://www.github.com/lindell/multi-gitter/compare/v0.42.0...v0.42.1) (2022-08-06)


### Bug Fixes

* **github:** allow only pull permission for print command ([#262](https://www.github.com/lindell/multi-gitter/issues/262)) ([582c706](https://www.github.com/lindell/multi-gitter/commit/582c70633e0dcef04699c62b160a57fedb50c00d))


### Miscellaneous

* **gitlab:** fixed subgroup spelling ([46f178a](https://www.github.com/lindell/multi-gitter/commit/46f178ab097fefb3d07b65f3522e3e516638107a))


### Dependencies

* update github.com/eiannone/keyboard digest to 0d22619 ([f38b2f5](https://www.github.com/lindell/multi-gitter/commit/f38b2f5b2bf17afab06ed709c6b5ba4c9c9a6075))
* update golang.org/x/oauth2 digest to 128564f ([63eab95](https://www.github.com/lindell/multi-gitter/commit/63eab95d69715ca1794f05447e3a588b70d99e62))
* update golang.org/x/oauth2 digest to 2104d58 ([c4605e2](https://www.github.com/lindell/multi-gitter/commit/c4605e2c8db5440fc68b5b779ee46b360c4459db))
* update module github.com/sirupsen/logrus to v1.9.0 ([1c8201a](https://www.github.com/lindell/multi-gitter/commit/1c8201a3b748b511effeb02b98ae1485f0e97b31))
* update module github.com/spf13/cobra to v1.5.0 ([eb15db4](https://www.github.com/lindell/multi-gitter/commit/eb15db4fc137d5a23f3d137a778e4235b26b4e74))
* update module github.com/spf13/viper to v1.12.0 ([ffdf7ae](https://www.github.com/lindell/multi-gitter/commit/ffdf7ae901b1f98fdd279fd5d1aeb72095bf2ad2))
* update module github.com/stretchr/testify to v1.8.0 ([81a757c](https://www.github.com/lindell/multi-gitter/commit/81a757c88d3942ff5295a8edfcec6ad99ff9e669))
* update module github.com/xanzy/go-gitlab to v0.68.0 ([78ca7f3](https://www.github.com/lindell/multi-gitter/commit/78ca7f343b846c56a7fdc76f59ec6a7591d7c3da))
* update module github.com/xanzy/go-gitlab to v0.68.2 ([dab0985](https://www.github.com/lindell/multi-gitter/commit/dab09851cf9b694096ec11084568685efcd860ff))

## [0.42.0](https://www.github.com/lindell/multi-gitter/compare/v0.41.0...v0.42.0) (2022-05-06)


### Features

* **github:** use graphql endpoint to get pull request status ([#242](https://www.github.com/lindell/multi-gitter/issues/242)) ([60bbbdf](https://www.github.com/lindell/multi-gitter/commit/60bbbdf526ce6e87d0f952fdae4858e4c6954952))


### Dependencies

* update github.com/gfleury/go-bitbucket-v1 digest to 711d7d5 ([0ecf9ee](https://www.github.com/lindell/multi-gitter/commit/0ecf9ee2e17b162218b251db63b5e9946369df39))
* update module github.com/mitchellh/mapstructure to v1.5.0 ([152c3b1](https://www.github.com/lindell/multi-gitter/commit/152c3b1aeebeb996221d28ca6072342e9188493c))
* update module github.com/spf13/cobra to v1.4.0 ([d378aa0](https://www.github.com/lindell/multi-gitter/commit/d378aa00a1b54f8482ed4d18d94ef47651b44273))
* update module github.com/spf13/viper to v1.11.0 ([a8e01ef](https://www.github.com/lindell/multi-gitter/commit/a8e01ef104ee09a38fe3d3551caa97fb5d00dc0e))
* update module github.com/stretchr/testify to v1.7.1 ([7506930](https://www.github.com/lindell/multi-gitter/commit/750693020fffae4aef811558fcdd65abf3662935))

## [0.41.0](https://www.github.com/lindell/multi-gitter/compare/v0.40.1...v0.41.0) (2022-03-28)


### Features

* **gitlab:** skip archived repos ([#240](https://www.github.com/lindell/multi-gitter/issues/240)) ([10df83a](https://www.github.com/lindell/multi-gitter/commit/10df83a80e486dc7f2c8ec085c1111d8fab42cfb))

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
