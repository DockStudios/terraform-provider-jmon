# Changelog

## [2.2.3](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/compare/v2.2.2...v2.2.3) (2023-08-21)


### Bug Fixes

* Perform deep compare of steps structure before updating state to avoid non-functional changes showing differences in Terraform plan ([aa64397](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/commit/aa64397b501e178e6f121d6d4b7c2d3dcc04fc69)), closes [#7](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/issues/7)

## [2.2.2](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/compare/v2.2.1...v2.2.2) (2023-08-11)


### Bug Fixes

* **docs:** Fix example for JSON matching ([85f98dc](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/commit/85f98dc3f5bfdbb3948c972197fdb1ac91c78dd9))
* Fix call to provider ([d69e539](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/commit/d69e539db70c918c28f1c72b02305d64c8ff7876))

## [2.2.1](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/compare/v2.2.0...v2.2.1) (2023-08-11)


### Bug Fixes

* Fix references to matthewjohn, after moving repository to dockstudios ([3be1124](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/commit/3be11242c51dbc10f3106b6fbc21049a780e0a82)), closes [#5](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/issues/5)

# [2.2.0](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/compare/v2.1.1...v2.2.0) (2023-08-11)


### Features

* Add support for setting check attributes ([e6d5328](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/commit/e6d532806716e087628d84c264c34bef3d23bb9d)), closes [#6](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/issues/6)

# [2.1.0](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/compare/v2.0.0...v2.1.0) (2023-08-10)


### Bug Fixes

* Add headers to all requests made in check resource ([a30f4e4](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/commit/a30f4e479e32a2dbb1b936d95190a3e718399a78)), closes [#3](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/issues/3)
* Add semantic-release config ([c71fedb](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/commit/c71fedbc1991a4a7eb49d10d6f25c3a2c76c3a94)), closes [#5](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/issues/5)


### Features

* Add api_key attribute to provider and pass API key header in request headers ([dc882a5](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/commit/dc882a52e8e593ce5cc371ae431756032ae6dab6)), closes [#3](https://gitlab.dockstudios.co.uk/pub/jmon/jmon-terraform-provider/issues/3)
