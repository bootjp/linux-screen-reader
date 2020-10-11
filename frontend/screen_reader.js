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
        }).catch(r => {
            console.log(r)

        })

    }
    send('http://localhost:50500', {text: selection})
}())