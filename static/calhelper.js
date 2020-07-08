function getRandomColor() {
    var letters = '0123456789ABCDEF';
    var color = '#';
    for (var i = 0; i < 6; i++) {
        color += letters[Math.floor(Math.random() * 16)];
    }
    return color;
};

function poplulateCalEvents() {
    calEvents = angular.element(document.getElementById('PirriControl')).scope().calEvents;
    return angular.element(document.getElementById('PirriControl')).scope().calEvents
}


