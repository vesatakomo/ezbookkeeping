const licenses = [
    {
        name: 'Gin',
        copyright: 'Copyright (c) 2014 Manuel Martínez-Almeida',
        url: 'https://github.com/gin-gonic/gin',
        licenseUrl: 'https://github.com/gin-gonic/gin/blob/master/LICENSE'
    },
    {
        name: 'GZIP gin\'s middleware',
        copyright: 'Copyright (c) 2017 Gin-Gonic',
        url: 'https://github.com/gin-contrib/gzip',
        licenseUrl: 'https://github.com/gin-contrib/gzip/blob/master/LICENSE'
    },
    {
        name: 'xorm',
        copyright: 'Copyright (c) 2013 - 2015 The Xorm Authors',
        url: 'https://xorm.io/xorm',
        licenseUrl: 'https://gitea.com/xorm/xorm/src/branch/master/LICENSE'
    },
    {
        name: 'Logrus',
        copyright: 'Copyright (c) 2014 Simon Eskildsen',
        url: 'https://github.com/sirupsen/logrus',
        licenseUrl: 'https://github.com/sirupsen/logrus/blob/master/LICENSE'
    },
    {
        name: 'cli',
        copyright: 'Copyright (c) 2016 Jeremy Saenz & Contributors',
        url: 'https://github.com/urfave/cli',
        licenseUrl: 'https://github.com/urfave/cli/blob/master/LICENSE'
    },
    {
        name: 'INI',
        copyright: 'Copyright 2014 Unknwon',
        url: 'https://gopkg.in/ini.v1',
        licenseUrl: 'https://github.com/go-ini/ini/blob/master/LICENSE'
    },
    {
        name: 'Package validator',
        copyright: 'Copyright (c) 2015 Dean Karn',
        url: 'https://github.com/go-playground/validator',
        licenseUrl: 'https://github.com/go-playground/validator/blob/master/LICENSE'
    },
    {
        name: 'jwt-go',
        copyright: 'Copyright (c) 2012 Dave Grijalva',
        url: 'https://github.com/dgrijalva/jwt-go',
        licenseUrl: 'https://github.com/dgrijalva/jwt-go/blob/master/LICENSE'
    },
    {
        name: 'otp',
        url: 'https://github.com/pquerna/otp',
        licenseUrl: 'https://github.com/pquerna/otp/blob/master/LICENSE'
    },
    {
        name: 'Go-MySQL-Driver',
        url: 'https://github.com/go-sql-driver/mysql',
        licenseUrl: 'https://github.com/go-sql-driver/mysql/blob/master/LICENSE'
    },
    {
        name: 'pq',
        copyright: 'Copyright (c) 2011-2013, \'pq\' Contributors Portions Copyright (C) 2011 Blake Mizerany',
        url: 'https://github.com/lib/pq',
        licenseUrl: 'https://github.com/lib/pq/blob/master/LICENSE.md'
    },
    {
        name: 'go-sqlite3',
        copyright: 'Copyright (c) 2014 Yasuhiro Matsumoto',
        url: 'https://github.com/mattn/go-sqlite3',
        licenseUrl: 'https://github.com/mattn/go-sqlite3/blob/master/LICENSE'
    },
    {
        name: 'Go Cryptography',
        copyright: 'Copyright (c) 2009 The Go Authors. All rights reserved.',
        url: 'https://golang.org/x/crypto',
        licenseUrl: 'https://github.com/golang/crypto/blob/master/LICENSE'
    },
    {
        name: 'Testify',
        copyright: 'Copyright (c) 2012-2020 Mat Ryer, Tyler Bunnell and contributors.',
        url: 'https://github.com/stretchr/testify',
        licenseUrl: 'https://github.com/stretchr/testify/blob/master/LICENSE'
    },
    {
        name: 'vue',
        copyright: 'Copyright (c) 2013-present, Yuxi (Evan) You',
        url: 'https://github.com/vuejs/vue',
        licenseUrl: 'https://github.com/vuejs/vue/blob/dev/LICENSE'
    },
    {
        name: 'vue-i18n',
        copyright: 'Copyright (c) 2016 kazuya kawaguchi',
        url: 'https://github.com/kazupon/vue-i18n',
        licenseUrl: 'https://github.com/kazupon/vue-i18n/blob/v8.x/LICENSE'
    },
    {
        name: 'vue-i18n-filter',
        copyright: 'Copyright (c) 2018 +v',
        url: 'https://github.com/chiaweilee/vue-i18n-filter',
        licenseUrl: 'https://github.com/chiaweilee/vue-i18n-filter/blob/master/LICENSE'
    },
    {
        name: 'vue-moment',
        copyright: 'Copyright (c) 2017 Brock Petrie',
        url: 'https://github.com/brockpetrie/vue-moment',
        licenseUrl: 'https://github.com/brockpetrie/vue-moment/blob/master/LICENSE'
    },
    {
        name: 'vue-clipboard2',
        copyright: 'Copyright (c) 2017 Inndy <inndy \\dot tw \\at gmail \\dot com>',
        url: 'https://github.com/Inndy/vue-clipboard2',
        licenseUrl: 'https://github.com/Inndy/vue-clipboard2/blob/master/LICENSE'
    },
    {
        name: 'core-js',
        copyright: 'Copyright (c) 2014-2020 Denis Pushkarev',
        url: 'https://github.com/zloirock/core-js',
        licenseUrl: 'https://github.com/zloirock/core-js/blob/master/LICENSE'
    },
    {
        name: 'Framework7',
        copyright: 'Copyright (c) 2014 Vladimir Kharlampidi',
        url: 'https://framework7.io/',
        licenseUrl: 'https://github.com/framework7io/framework7/blob/master/LICENSE'
    },
    {
        name: 'Framework7-vue',
        copyright: 'Copyright (c) 2014 Vladimir Kharlampidi',
        url: 'https://framework7.io/vue/',
        licenseUrl: 'https://github.com/framework7io/framework7/blob/master/LICENSE'
    },
    {
        name: 'Framework7-icons',
        copyright: 'Copyright (c) 2016 Vladimir Kharlampidi',
        url: 'https://framework7.io/icons/',
        licenseUrl: 'https://github.com/framework7io/framework7-icons'
    },
    {
        name: 'axios',
        copyright: 'Copyright (c) 2014-present Matt Zabriskie',
        url: 'https://github.com/axios/axios',
        licenseUrl: 'https://github.com/axios/axios/blob/master/LICENSE'
    },
    {
        name: 'Moment.js',
        copyright: 'Copyright (c) JS Foundation and other contributors',
        url: 'https://momentjs.com',
        licenseUrl: 'https://github.com/moment/moment/blob/develop/LICENSE'
    },
    {
        name: 'js-cookie',
        copyright: 'Copyright (c) 2018 Copyright 2018 Klaus Hartl, Fagner Brack, GitHub Contributors',
        url: 'https://github.com/js-cookie/js-cookie',
        licenseUrl: 'https://github.com/js-cookie/js-cookie/blob/master/LICENSE'
    }
];

export default licenses;
