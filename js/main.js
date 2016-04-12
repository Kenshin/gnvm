"use strict";
$( document ).ready( function() {

    var lang = navigator.language;

    // set 'en-XX' to 'en'
    if ( lang.split("-")[0] == "en" ) {
        lang = "en";
    }

    // set cookie
    var newLang = $.cookie( "lang" );
    if ( newLang !== undefined ) {
        lang = newLang;
    }

    // get http://xxxx.index.html?lang=xxx
    newLang = getQueryString( "lang" );
    if ( newLang !== null ) {
        lang = newLang;
    }

    // set lower case
    lang = lang.toLowerCase();

    // set default language
    if ( lang != "en" && lang !== "zh-cn" && lang !== "zh-tw" ) {
        lang = "en";
    }

    // set options
    var options = {
        load         : "current",
        lng          : lang,
        lowerCaseLng : true,
        fallbackLng  : false,
        cookieName   : "lang",
        useCookie    : true,
        detectLngQS  : "lang"
    };

    console.log(options)

    // i18n init
    i18n.init( options, function(t) {
        $(".side").i18n();
    });

});

function getQueryString( name ) {
    var reg = new RegExp( "(^|&)" + name + "=([^&]*)(&|$)", "i" );
    var r = window.location.search.substr(1).match( reg );
    if (r !== null) return unescape(r[2]); return null;
}
