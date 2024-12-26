

flatpickr("#trip-duration", {
    enableTime: true,
    noCalendar: true,
    dateFormat: "H:i",
    time_24hr: true,
    inline: true,
});

flatpickr("#alternate-duration", {
    enableTime: true,
    noCalendar: true,
    dateFormat: "H:i",
    time_24hr: true,
    inline: true,
});

htmx.on('htmx:beforeHistorySave', function() {
    let tripDurationElem = document.getElementById("trip-duration");
    let alternateDurationElem = document.getElementById("alternate-duration");

    if (tripDurationElem && tripDurationElem._flatpickr) {
        tripDurationElem._flatpickr.destroy();
    }
    if (alternateDurationElem && alternateDurationElem._flatpickr) {
        alternateDurationElem._flatpickr.destroy();
    }
})
