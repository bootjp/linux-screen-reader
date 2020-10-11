// ==UserScript==
// @name            Selection Text Reader
// @namespace       https://bootjp.me/
// @description     Context menu to execute UserScript
// @version         0.1
// @author          bootjp
// @include         *
// @grant           GM_xmlhttpRequest
// @run-at          context-menu
// ==/UserScript==]


(function() {
    'use strict';
    (function () {
        const selection = document.getSelection().toString();
        function send(url = '', data = {}) {
            fetch(url, {
                method: 'POST',
                mode: 'cors',
                cache: 'no-cache',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: new URLSearchParams(data).toString()
            }).catch(e => {
                // csp blocking, fallback clipboard
                let last = ""
                navigator.clipboard.readText()
                .then(function(text){
                   last = text;
                });

                navigator.clipboard.writeText( `--screen_reader ${data.text}`);

                // 副作用を減らすためもとのクリップボードに戻す
                setTimeout(function () {
                    navigator.clipboard.writeText(last.toString())
                }, 1001)
            })
        }
        send('http://localhost:50500', {text: selection})
    }());
})();