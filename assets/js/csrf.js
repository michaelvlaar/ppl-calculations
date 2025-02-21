document.body.addEventListener('htmx:configRequest', function (evt) {
    evt.detail.headers['X-CSRF-Token'] = document.querySelector('meta[name="csrf"]').getAttribute("content");
});
