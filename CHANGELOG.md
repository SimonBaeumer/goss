# v0.6.0

 - Fixed a bug where DNS record lookups never were fired wihout a nameserver
 - Removed global `--package` option, added `package-manager` property to `package` resource
 - Removed code dependency on `*cli.Context`

# v0.5.0

 - Add certificate authentication to `http` resource
 
# v0.4.0

 - Added http header validation
 - Added headers in http resource
 - Moved project from aelsabbahy/goss to SimonBaeumer/goss
 