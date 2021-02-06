// ==UserScript==
// @name            Selection Text Reader
// @namespace       https://bootjp.me/
// @description     Read the selected text from the context menu with the goolge text to speech API
// @version         0.1
// @author          bootjp
// @include         *
// @grant           GM_xmlhttpRequest
// @grant           GM_setClipboard
// @run-at          context-menu
// ==/UserScript==

(function () {
    'use strict';
    const selection = document.getSelection().toString();
    const host = 'localhost:50500';
    const schema = 'http';
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
            if (e.message !== 'Failed to fetch') {
                return
            }

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
    send(`${schema}://${host}/tts/speech`, {text: selection})
}());