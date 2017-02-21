(function() {
	
    var app = angular.module('pirriweb', ['chart.js', 'ngCookies','ngMessages', 'ngMaterial']);
    app.Root = '/';
    app.config(['$interpolateProvider',
        function($interpolateProvider) {
            $interpolateProvider.startSymbol('{[');
            $interpolateProvider.endSymbol(']}');
        }
    ]);
	
	/*
	angular.module 'app.components'
 .directive 'autoChangeStringDates', ->
   directive =
     restrict: 'A'
     require: 'ngModel'
     priority: 2000
     link: (scope, el, attrs, ngModelController) ->
       ngModelController.$formatters.push((input) ->
         if typeof input == Date
          return input
         else
           return Date.create(input, {fromUTC: true})
       )
    return
	
	<md-datepicker ng-model='myModel' auto-change-string-dates></md-datepicker>

	*/

    app.controller('PirriControl', function($rootScope, $scope, $http, $timeout, $filter, $cookies) {
        $rootScope.updateInterval = 6000;
        $scope.chartData1 = {
			labels: [],
			series: [],
			data:[]
		};
        $scope.chartData2 = {			
			labels: [],
			series: [],
			data:[]
			}; 
        $scope.chartData3 = {			
			labels: [],
			series: [],
			data:[]
			};
        $scope.chartData4 = {			
			labels: [],
			series: [],
			data:[]
			};
        $scope.beatheart = false;

        $scope.randomColor = function() {
            var letters = '0123456789ABCDEF';
            var color = '#';
            for (var i = 0; i < 6; i++) {
                color += letters[Math.floor(Math.random() * 16)];
            }
            return color;
        };

        $scope.calEvents = []
        this.getCalEvents = function() {
            $scope.beatheart = true;
            $http.get('/schedule/all')
                .then(function(response) {
					$scope.calEvents = response.data.StationSchedule;
					$scope.beatheart = false;
                })
        };

		$scope.StrToDate = function (str) {
            return new Date(str);
        }
		$scope.setTabCookie = function() {
  			$cookies.put('lastTab', $scope.currentPage);
		}
		
		$scope.schedule = []
        $scope.weatherData = {};
        $scope.settingsData = {};
        $scope.currentPage = 'home'; // history / home / settings / add
        $scope.stations = undefined;
        $scope.navTitle = "All Stations";
        $scope.gpio_pins = undefined;
        $scope.searchResults = {};
        $scope.searchText = "";
        $scope.showSearchResults = false;
        $scope.history = [];
        $scope.historyScope = "All Stations";
        $scope.schedule = undefined;
        $scope.gpio_add_model = {
            default_message: "Select GPIO",
            GPIO: undefined
        };
        $scope.edit_station_model = {
            tempID: undefined,
            GPIO: undefined,
            Notes: undefined
        };
        $scope.durationIntervals = [1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60];
        $scope.show_gpio_diagram = false;

        this.filterForKeys = function(searchText) {
            $scope.searchResults = [];
            $scope.stations.forEach(function(k) {
                var n = k.search(searchText);
                if (n >= 0) {
                    $scope.searchResults[k] = true;
                }
            });
            if (Object.keys($scope.searchResults).length > 0) {
                $scope.showSearchResults = true;
            } else {
                $scope.showSearchResults = false;
            }
            if (searchText === "") {
                $scope.searchResults = [];
                $scope.showSearchResults = false;
            }
        };

        this.resetAddForm = function() {
            $scope.gpio_add_model = {
                default_message: "Select GPIO",
                GPIO: undefined
            };
        };

        this.setGPIO = function(gpio) {
            $scope.gpio_add_model.GPIO = gpio;
        };

        this.setEditingStationInfo = function(station) {
            $scope.stationModel = station;
        };


		
        this.setPage = function(pageName) {
            $scope.currentPage = pageName;
        };


        $scope.dripnodes = {};
        $scope.watercost = 0.0021;
        this.getWaterUsageStats = function() {
            $http.get('/stats/gallons')
                .then(function(response) {
                    $scope.dripnodes = response.data.water_usage;
                })
            this.getWaterNodeEntries();
        };

        this.getUsageDataForChart1 = function() {
            $http.get('/stats?id=1')
                .then(function(response) {
                    $scope.chartData1.labels = response.data.chartData.labels;
                    $scope.chartData1.series = response.data.chartData.series;
                    $scope.chartData1.data = response.data.chartData.data;
                })

            $scope.chartData1.options = {
                title: {
                    display: true,
                    text: 'Total Usage in Minutes (last 30 days)'
                },
                scaleStartValue: 0,
                legend: {
                    display: true,
                    labels: {
                        //fontColor: 'rgb(255, 99, 132)'
                    }
                },
            };
        };

        $scope.monthly_cost = 0;
        this.calcMonthlyCost = function() {
            $scope.monthly_cost = 0;
            angular.forEach($scope.dripnodes, function(value, key) {
                $scope.monthly_cost += value['usage_last_30'] * $scope.watercost
            })
        };

        this.getUsageDataForChart2 = function() {
            $http.get('/stats?id=2')
                .then(function(response) {
                    $scope.chartData2.labels = response.data.chartData.labels;
                    $scope.chartData2.series = response.data.chartData.series;
                    $scope.chartData2.data = response.data.chartData.data;
                })
            $scope.chartData2.options = {
                scaleStartValue: 0,
                title: {
                    display: true,
                    text: 'Usage in Minutes by Day of Week (last 7 days)'
                },
                scaleStartValue: 0,
                legend: {
                    display: true,
                    labels: {}
                },
            };
        };


        this.getUsageDataForChart3 = function() {
            $http.get('/stats?id=3')
                .then(function(response) {
                    $scope.chartData3.labels = response.data.chartData.labels;
                    $scope.chartData3.series = response.data.chartData.series;
                    $scope.chartData3.data = response.data.chartData.data;
                })
            $scope.chartData3.options = {
                scaleStartValue: 0,
                title: {
                    display: true,
                    text: 'Usage in Minutes Per Station by Day of Week (last 30 days)'
                },
                scaleStartValue: 0,
                legend: {
                    display: true,
                    labels: {}
                },
            };
        };

        this.getUsageDataForChart4 = function() {
            $http.get('/stats?id=4')
                .then(function(response) {
                    $scope.chartData4.labels = response.data.chartData.labels;
                    $scope.chartData4.series = response.data.chartData.series;
                    $scope.chartData4.data = response.data.chartData.data;
                })
            $scope.chartData4.options = {
                scaleStartValue: 0,
                title: {
                    display: true,
                    text: 'Station Activity by Hour of the Day (last 30 days)'
                },
                scaleStartValue: 0,
                legend: {
                    display: true,
                    labels: {}
                },
            };
        };

        this.loadStatsData = function() {
            $scope.beatheart = true;
            Chart.defaults.global.defaultFontColor = "#fff";
//            this.getUsageDataForChart1();
//            this.getUsageDataForChart2();
//            this.getUsageDataForChart3();
//            this.getUsageDataForChart4();
            $scope.beatheart = false;
        };


        $scope.stationModel = {}
        this.submitEditStation = function() {
            $http.post('/station/edit', $scope.stationModel)
                .success(function(response) {})
            $scope.stationModel = {};
            $scope.stationModel = undefined;
        };
        this.submitDeleteStation = function() {};
        this.submitAddStation = function() {};

        $scope.scheduleModel = {};

        this.addScheduleButton = function() {
            $scope.scheduleModel = undefined;
            $scope.scheduleModel = {
                ID: 0,
				tempID: 0,
                Sunday: false,
                Monday: false,
                Tuesday: false,
                Wednesday: false,
                Thursday: false,
                Friday: false,
                Saturday: false,
                Repeating: false,
                StationID: undefined,
                StartTime: undefined,
                StartDate: undefined,
                EndDate: 30000000,
                Duration: 0,
            };
            $scope.schedule.unshift($scope.scheduleModel)
        };

        this.addScheduleButtonFromCalendar = function(startTime) {
            $scope.scheduleModel = undefined;
            $scope.scheduleModel = {
                ID: 0,
				tempID: 0,
                Sunday: false,
                Monday: false,
                Tuesday: false,
                Wednesday: false,
                Thursday: false,
                Friday: false,
                Saturday: false,
                Repeating: false,
                StationID: undefined,
                StartTime: undefined,
                StartDate: undefined,
                EndDate: 30000000,
                Duration: 0,
            };
            $scope.currentPage = 'calendar';

            $scope.schedule.unshift($scope.scheduleModel)
        };

        this.submitEditSchedule = function() {
            $http.post('/schedule/edit', $scope.scheduleModel)
                .then(function(response) {
					$scope.schedule = response.data.stationSchedules
				})
            $scope.scheduleModel = {};
            $scope.scheduleModel = undefined;
            this.refresh();
        };

        this.submitAddSchedule = function() {
            this.convertScheduleBoolToInt();
            $http.post('/schedule/add', $scope.scheduleModel)
                .success(function(response) {})
                // cleanup
            $scope.scheduleModel = {};
            $scope.scheduleModel = undefined;
            this.refresh();
        };

        this.mapModelForSchedEdit = function(currentModel) {
            $scope.scheduleModel = currentModel;
        };

        this.mapModelForSchedEditFromCalClick = function(id) {
            $scope.currentPage = 'calendar';
            console.log($scope.currentPage);
            var sch = $filter('filter')($scope.schedule, {
                id: id
            })[0];
            $scope.scheduleModel = sch;
        };

        this.submitDeleteSchedule = function(schedule_id) {
            $http.post('/schedule/delete', {
                    ID: schedule_id
                })
                .then(function(response) {})
            $scope.scheduleModel = {};
            $scope.scheduleModel = undefined;
            this.refresh();
        };

        $scope.singleRunModel = {};
        this.submitSingleRun = function() {
            $http.post('/station/run', $scope.singleRunModel)
                .then(function(response) {})
            $scope.singleRunModel = {};
            $scope.singleRunMinField = undefined;
        };

        this.refresh = function() {
            this.getSchedule();
            this.loadStations();
            this.getLastStationRun();
            this.getNextStationRun();
            this.loadGPIO();
            this.loadStatsData();
            this.getWaterUsageStats();
            this.loadSettings();
            this.loadWeather();
        };
        this.loadStations = function() {
            $http.get('/station/all')
                .then(function(response) {
//					console.log(response.data)
                    $scope.stations = response.data.stations;
//                    angular.forEach($scope.stations, function(value, key) {
//                        value['cal_color'] = $scope.randomColor();
//                    })
                })
        };

        this.loadSettings = function() {
            $http.get('/settings/load')
                .then(function(response) {
                    $scope.settingsData = response.data.data;
                })
        };

        this.loadWeather = function() {
            $scope.beatheart = true;
            $http.get('/weather')
                .then(function(response) {
//                    $scope.weatherData = response.data;
//                    $scope.weatherData.current.sys.sunrise_t = moment(data.current.sys.sunrise * 1000).fromNow();
//                    $scope.weatherData.current.sys.sunset_t = moment(data.current.sys.sunset * 1000).fromNow();

                })
        };

        this.loadGPIO = function() {
            $http.get('/gpio/list')
				.then(function(response) {                    
					$scope.gpio_pins = response.data.gpio_pins;
                })
        };

        this.loadHistory = function(station) {
            var query = '?station=' + station + '&earliest=-168';
            $http.get('/history' + query)
                .then(function(response) {
                    $scope.history = response.data.history;
                })
        };

        this.prettyTime = function(uglyTime) {
            if (uglyTime !== undefined && uglyTime !== null) {
                // console.log(uglyTime)
                var pt = moment(uglyTime).calendar();
                return pt
            } else {
                return "Never"
            }
        }
        this.getSchedule = function() {
            $http.get('/schedule/all')
                .then(function(response) {
                    $scope.schedule = response.data.stationSchedules;
                })
        };

        $scope.lastStationRunHash = {}
        this.getLastStationRun = function() {
            $http.get('/station/lastruns')
                .then(function(response) {
                    $scope.lastStationRunHash = response.data.lastrunlist;
                })
                // console.log($scope.lastStationRunHash);
        };

        $scope.nextStationRunHash = {}
        this.getNextStationRun = function() {
            $http.get('/station/nextruns')
                .then(function(response) {
                    $scope.nextStationRunHash = response.data.nextrunlist;
                })
                
        };

        $scope.waterNodeEntries = [];
        $scope.waterNodeModel = {};
        this.getWaterNodeEntries = function() {
            $http.get('/station/nodes')
                .then(function(response) {
                    $scope.waterNodeEntries = response.data.dripnodes;
                })
        }
        this.submitEditNodeEntry = function() {
            $http.post('/station/nodes', $scope.waterNodeModel)
                .then(function(response) {
                    // console.log($scope.singleRunModel, data)
                })
                // cleanup
            $scope.waterNodeModel = undefined;
            $scope.waterNodeModel = {};
        };

        this.submitAddNodeEntry = function() {
            $scope.waterNodeModel.new = true;
            $http.post('/station/nodes', $scope.waterNodeModel)
                .then(function(response) {
                    // console.log($scope.singleRunModel, data)
                })
                // cleanup
            $scope.waterNodeModel = undefined;
            $scope.waterNodeModel = {};
        };

        this.submitDeleteNodeEntry = function(nodeid) {
            $scope.waterNodeModel.id = nodeid;
            $http.post('/station/nodes/delete', $scope.waterNodeModel)
                .then(function(response) {
                    // console.log($scope.singleRunModel, data)
                })
                // cleanup
            $scope.waterNodeModel = undefined;
            $scope.waterNodeModel = {};
        };

        this.mapModelForWaterNodeEdit = function(currentModel) {
            $scope.waterNodeModel = currentModel;
            // console.log($scope.scheduleModel)
        };

        this.addWaterNodeButton = function() {
            $scope.waterNodeModel = undefined;
            $scope.waterNodeModel = {
                id: '-',
                sid: 'Select Station ID',
                gph: '',
                count: 0,
                new: true
            };
            // console.log($scope.scheduleModel)
            $scope.waterNodeEntries.unshift($scope.waterNodeModel)
        };

        this.autoLoader = function() {
            this.getCalEvents();
            this.getSchedule();
            this.loadStations();
            this.getLastStationRun();
            this.getNextStationRun();
            this.loadGPIO();
            this.loadStatsData();
            this.loadHistory();
            this.calcMonthlyCost();
            //this.loadSettings();
            //this.loadWeather();
			if ($cookies.get('lastTab') != undefined) {
				$scope.currentPage = $cookies.get('lastTab');
			}

        };
        $scope.loader = this.autoLoader;
        // $scope.intervalFunction = function() {
        //     $timeout(function() {
        //         $scope.loader();
        //         $scope.intervalFunction();
        //     }, $rootScope.updateInterval)
        // };
        //$scope.intervalFunction();

        this.autoLoader();

    });
})();