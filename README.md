Local PHP Security Checker
==========================

The Local PHP Security Checker is a command line tool that checks if your PHP
application depends on PHP packages with known security vulnerabilities. It
uses the [Security Advisories Database][1] behind the scenes.

Download a binary from the [Releases page on Github][2], rename it to
`local-php-security-checker` and make it executable.

From a directory containing a PHP project that uses Composer, check for known
vulnerabilities by running the binary without arguments or flags:

    $ local-php-security-checker

You can also pass a `--path` to check a specific directory:

    $ local-php-security-checker --path=/path/to/php/project
    $ local-php-security-checker --path=/path/to/php/project/composer.lock

By default, the output is optimized for terminals, change it via the `--format`
flag (supported formats: `ansi`, `markdown`, `json`, and `yaml`):

    $ local-php-security-checker --format=json

When running the command, it checks for an updated vulnerability database and
downloads it from Github if it changed since the last run. If you want to avoid
the HTTP round-trip, use `--local`. To force a database update without checking
for a project, use `--update-cache`.

If you want to continuously check for security issues on your applications in
production, you can use this tool in combination with [croncape][3] to get an
email whenever a new security issue is detected:

    MAILTO=sysadmins@example.com
    50 23 * * * croncape php-security-checker --path=/path/to/php/project

[1]: https://github.com/FriendsOfPHP/security-advisories
[2]: https://github.com/fabpot/local-php-security-checker/releases
[3]: https://github.com/symfonycorp/croncape
