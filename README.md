Local PHP Security Checker
==========================

The Local PHP Security Checker is a command line tool written in GO that checks if your PHP
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
     _                    _   ____  _   _ ____
    | |    ___   ___ __ _| | |  _ \| | | |  _ \
    | |   / _ \ / __/ _` | | | |_) | |_| | |_) |
    | |__| (_) | (_| (_| | | |  __/|  _  |  __/
    |_____\___/ \___\__,_|_| |_|   |_| |_|_|
     ____                       _ _            ____ _               _
    / ___|  ___  ___ _   _ _ __(_) |_ _   _   / ___| |__   ___  ___| | _____ _ __
    \___ \ / _ \/ __| | | | '__| | __| | | | | |   | '_ \ / _ \/ __| |/ / _ \ '__|
     ___) |  __/ (__| |_| | |  | | |_| |_| | | |___| | | |  __/ (__|   <  __/ |
    |____/ \___|\___|\__,_|_|  |_|\__|\__, |  \____|_| |_|\___|\___|_|\_\___|_|
                                      |___/
    1.0, built at Fri, 15 Jan 2021

    database known vulnerabilities 1065
    database version Fri, 15 Jan 2021 17:39:07 GMT

    composer://composer/ca-bundle@1.2.5
    composer://doctrine/annotations@v1.2.0 vulnerability found !
    composer://doctrine/cache@v1.6.2
    composer://doctrine/collections@v1.4.0
    composer://doctrine/common@v2.7.3
    composer://doctrine/dbal@v2.5.13
    composer://doctrine/doctrine-bundle@1.10.3
    composer://doctrine/doctrine-cache-bundle@1.3.5
    composer://doctrine/inflector@v1.2.0
    composer://doctrine/instantiator@1.0.5
    composer://doctrine/lexer@1.0.2
    composer://fig/link-util@1.1.0
    composer://friendsofsymfony/rest-bundle@2.3.1
    composer://incenteev/composer-parameter-handler@v2.1.3
    composer://jdorn/sql-formatter@v1.2.17
    composer://jms/metadata@1.7.0
    composer://jms/parser-lib@1.0.0
    composer://jms/serializer@1.14.0
    composer://jms/serializer-bundle@2.3.1
    composer://monolog/monolog@1.25.3
    composer://nelmio/cors-bundle@1.5.4
    composer://onup/utilsbundle@dev-php7-sf3
    composer://paragonie/random_compat@v9.99.99
    composer://phpcollection/phpcollection@0.5.0
    composer://phpoption/phpoption@1.7.2
    composer://psr/cache@1.0.1
    composer://psr/container@1.0.0
    composer://psr/link@1.0.0
    composer://psr/log@1.1.2
    composer://psr/simple-cache@1.0.1
    composer://sensio/distribution-bundle@v5.0.21
    composer://sensio/framework-extra-bundle@v3.0.29
    composer://sensiolabs/security-checker@v4.1.8
    composer://symfony/monolog-bundle@v3.2.0
    composer://symfony/polyfill-apcu@v1.13.1
    composer://symfony/polyfill-ctype@v1.13.1
    composer://symfony/polyfill-intl-icu@v1.13.1
    composer://symfony/polyfill-mbstring@v1.13.1
    composer://symfony/polyfill-php56@v1.13.1
    composer://symfony/polyfill-php70@v1.13.1
    composer://symfony/polyfill-util@v1.13.1
    composer://symfony/symfony@v3.4.35
    composer://twig/twig@v2.12.3
    composer://willdurand/jsonp-callback-validator@v1.1.0
    composer://willdurand/negotiation@v2.3.1
    composer://zircote/swagger-php@2.0.14
    composer://hamcrest/hamcrest-php@v2.0.0
    composer://mockery/mockery@1.3.1
    composer://sensio/generator-bundle@v3.1.7
    composer://symfony/phpunit-bridge@v3.4.36

    Symfony Security Check Report
    =============================

    1 package has known vulnerabilities.

    doctrine/annotations (v1.2.0)
    -----------------------------

    * [CVE-2015-5723][]: Security Misconfiguration Vulnerability in various Doctrine projects

    [CVE-2015-5723]: https://www.doctrine-project.org/2015/08/31/security_misconfiguration_vulnerability_in_various_doctrine_projects.html

    Note that this checker can only detect vulnerabilities that are referenced in the security advisories database.
    Execute this command regularly to check the newly discovered vulnerabilities.


The output can been optimized for terminals, change it via the `--quiet` and `--format`
flags (supported formats: `ansi`, `markdown`, `json`, and `yaml`):

    $ local-php-security-checker --quiet --format=json
    {
        "doctrine/annotations": {
            "version": "v1.2.0",
            "advisories": [
                {
                    "title": "Security Misconfiguration Vulnerability in various Doctrine projects",
                    "link": "https://www.doctrine-project.org/2015/08/31/security_misconfiguration_vulnerability_in_various_doctrine_projects.html",
                    "cve": "CVE-2015-5723"
                }
            ]
        }
    }

When running the command, it checks for an updated vulnerability database and
downloads it from [Security Advisories Database][1] if it changed since the last run. If you want to avoid
the HTTP round-trip, use `--local`. To force a database update without checking
for a project, use `--update-cache`.

[1]: https://github.com/FriendsOfPHP/security-advisories
[2]: https://github.com/fabpot/local-php-security-checker/releases
