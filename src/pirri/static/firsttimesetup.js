(function() {
    var app = angular.module('pirrisetup', []);
    app.Root = '/';
    app.config(['$interpolateProvider',
        function($interpolateProvider) {
            $interpolateProvider.startSymbol('{[');
            $interpolateProvider.endSymbol(']}');
        }
    ]);

    app.controller('PirriSetupControl', function($rootScope, $scope, $http, $timeout, $filter) {
        $scope.settingsModel = undefined;


        this.populateCurrentOrDefaultSettings = function() {
            $http.get('/settings/get')
                .success(function(data, status, headers, config) {
                    $scope.settingsModel = data;
                })
                .error(function(data, status, headers, config) {})
        }
        this.updateSettings = function() {
            $http.get('/settings/post')
                .success(function(data, status, headers, config) {

                })
                .error(function(data, status, headers, config) {})
        }
    });
})();