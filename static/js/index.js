document.addEventListener('DOMContentLoaded', function() {
    const submitBtn = document.getElementById('submitBtn')
    const example = document.getElementById('example')
    example.innerHTML = `
    Original URL: http://www.example.com
    code: abcdef
    Shortened URL: ${window.location.origin}/abcdef
    `;

    submitBtn.addEventListener('click', function(event) {
        event.preventDefault();
        const url = document.getElementById('url').value;
        const resultDiv = document.getElementById('result');

        const formData = new FormData();
        formData.append('url', url);
        fetch('/new', {
            method: 'POST',
            body: formData,
        })
            .then(res => res.json())
            .then(function(json) {
                const shortenedUrl = window.location.origin + '/' + json.code;

                const shortenedUrlLabel = document.createElement('span')
                shortenedUrlLabel.innerHTML = 'Shortened URL: ';

                const anchor = document.createElement('a')
                anchor.setAttribute('href', shortenedUrl);
                anchor.setAttribute('target', '_blank');
                anchor.innerHTML = shortenedUrl;

                resultDiv.innerHTML = "";
                resultDiv.appendChild(shortenedUrlLabel);
                resultDiv.appendChild(anchor);
            })
            .catch(function(err) { console.error(err) });
    })
});
