flatpickr("#trip-duration", {
    enableTime: true,
    noCalendar: true,
    dateFormat: "H:i",
    time_24hr: true,
    inline: true,
    defaultDate: "{{ if .TripDuration }}{{ .TripDuration }}{{ else }}01:00{{ end }}"
});

flatpickr("#alternate-duration", {
    enableTime: true,
    noCalendar: true,
    dateFormat: "H:i",
    time_24hr: true,
    inline: true,
    defaultDate: "{{ if .AlternateDuration }}{{ .AlternateDuration }}{{ else }}00:30{{ end }}"
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
