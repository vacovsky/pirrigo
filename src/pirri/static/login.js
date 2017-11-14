(function() {

    var app = angular.module('pirrilogin', [
        'chart.js',
        'ngCookies',
        'base64'
    ]);
    app.Root = '/';

    app.controller('PirriLoginControl', function($scope, $rootScope, $http, $base64, $cookies, $window) {

        $scope.username = undefined;
        $scope.password = undefined;
        $scope.banner = undefined;

        this.encodeForAuth = function(username, password) {
            return $base64.encode(username + ":" + password);
        };

        this.Login = function(username, password) {
            var encoded = this.encodeForAuth(username, password);
            $http.defaults.headers.common['Authorization'] = 'Basic ' + encoded;

            $http.post('/login/verify')
                .then(function(response) {
                    if (response.status === 200) {
                        $cookies.put('Authorization', $http.defaults.headers.common['Authorization']);
                        $window.location.href = '/home';
                    } else {
                        alert('Incorrect login credentials...');
                    }
                });
        };

        this.getBanner = function() {
            $http.get('/metadata').then(function(response){
                console.log(response.data.banner);
                $scope.banner = console.log(response.data.banner);
            })
        }
    });
})();