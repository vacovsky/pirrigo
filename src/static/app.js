(function () {

    var app = angular.module('pirriweb', [
        'chart.js',
        'ngCookies',
        'ngMessages',
        // 'ngMaterial',
        'ngAnimate',
        // 'ngSanitize',
        // 'ui.calendar',
        // 'ui.bootstrap',
        'ui.rCalendar',
        'angularMoment'
    ]).factory('moment', function ($window) {
        return $window.moment;
    });
    app.Root = '/';

    app.filter('reverse', function () {
        return function (items) {
            return items.slice().reverse();
        };
    });

    app.controller('PirriControl', function ($scope, $rootScope, $http, $timeout, $filter, $cookies, $scope, $compile) {


        $rootScope.updateInterval = 6000;
        $scope.events = [{
            title: 'All Day Event',
            start: new Date(y, m, 1),
            color: 'orange'
        }];
        $scope.eventSource = [];

        $scope.runStatus = {
            IsIdle: true
        };

        $scope.intervalFunction = function () {
            $timeout(function () {
                $scope.getRunStatus();
                $scope.intervalFunction();
            }, $rootScope.updateInterval);
        };

        $scope.intervalFunction();

        $scope.getRunStatus = function () {
            $http.get('/status/run')
                .then(function (response) {
                    $scope.runStatus = response.data;
                    var st = new Date(response.data.StartTime);
                    var fin = new Date(
                        st.getFullYear(),
                        st.getMonth(),
                        st.getDate(),
                        st.getHours(),
                        st.getMinutes(),
                        st.getSeconds() + $scope.runStatus.Duration
                    );
                    $scope.runStatus.fin = fin;
                    $scope.runStatus.td = (new Date() - new Date($scope.runStatus.StartTime).getTime()) / 1000;
                    $scope.runStatus.pComplete = $scope.runStatus.td / $scope.runStatus.Duration;
                });
        };

        this.cancelStationRun = function () {
            $http.get('/status/cancel')
                .then(function (response) {
                    $scope.runStatus = response.data;
                });
        };

        this.login = function () {
            $http.post('/home', {
                username: username,
                password: password
            })
                .then(function (response) {
                    // callback($window.location.href); or var landingUrl = "http://" + $window.location.host + "/login";
                    // $window.location.href = landingUrl;
                });
        };

        $scope.chartData1 = {
            labels: [],
            series: [],
            data: [],
            options: {
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
            }
        };
        $scope.chartData2 = {
            labels: [],
            series: [],
            data: [],
            options: {
                scaleStartValue: 0,
                title: {
                    display: true,
                    text: 'Usage in Minutes by Day of Week (last 7 days)'
                },
                legend: {
                    display: true,
                    labels: {}
                },
            }
        };
        $scope.chartData3 = {
            labels: [],
            series: [],
            data: [],
            options: {
                scaleStartValue: 0,
                title: {
                    display: true,
                    text: 'Usage in Minutes Per Station by Day of Week (last 30 days)'
                },
                legend: {
                    display: true,
                    labels: {}
                },
            }
        };
        $scope.chartData4 = {
            labels: [],
            series: [],
            data: [],
            options: {
                scaleStartValue: 0,
                title: {
                    display: true,
                    text: 'Station Activity in Minutes by Hour of the Day (last 30 days)'
                },
                legend: {
                    display: true,
                    labels: {}
                },
            }
        };
        $scope.beatheart = false;

        /*
        Datepicker Crap
        */
        $scope.formats = ['dd-MMMM-yyyy', 'yyyy/MM/dd', 'dd.MM.yyyy', 'shortDate'];
        $scope.format = $scope.formats[0];
        $scope.altInputFormats = ['M!/d!/yyyy'];

        $scope.dateOptions = {
            // dateDisabled: disabled,
            formatYear: 'yy',
            maxDate: new Date(2050, 5, 22),
            minDate: new Date(),
            startingDay: 1
        };

        $scope.inlineOptions = {
            // customClass: getDayClass,
            minDate: new Date(),
            showWeeks: true
        };

        $scope.logs = [];

        // END datepicker crap
        $scope.openCalPicker = function (scheduleid) {
            var result = $.grep($scope.schedule, function (e) {
                return e.ID == scheduleid;
            });
            result[0].calOpen = true;
            //	    	$scope.calPicker.opened = true;
        };


        $scope.randomColor = function () {
            var letters = '0123456789ABCDEF';
            var color = '#';
            for (var i = 0; i < 6; i++) {
                color += letters[Math.floor(Math.random() * 16)];
            }
            return color;
        };

        $scope.calEvents = [];
        this.getCalEvents = function () {
            $scope.beatheart = true;
            $http.get('/schedule/all')
                .then(function (response) {
                    $scope.schedule = response.data.stationSchedules;
                    $scope.loadCalEvents();
                    $scope.beatheart = false;
                });
        };

        $scope.StrToDate = function (str) {
            return new Date(str);
        };
        $scope.setTabCookie = function () {
            $cookies.put('lastTab', $scope.currentPage);
        };

        $scope.schedule = [];
        $scope.weatherData = {};
        $scope.settingsData = {};
        $scope.currentPage = 'home'; // history / home / settings / add
        $scope.stations = [];
        $scope.navTitle = "All Stations";
        $scope.gpio_pins = [];
        $scope.searchResults = [];
        $scope.searchText = "";
        $scope.showSearchResults = false;
        $scope.history = [];
        $scope.historyScope = "All Stations";
        $scope.gpio_add_model = {
            default_message: "Select GPIO",
            GPIO: undefined
        };

        $scope.durationIntervals = [1, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60];
        $scope.show_gpio_diagram = false;

        this.resetAddForm = function () {
            $scope.gpio_add_model = {
                default_message: "Select GPIO",
                GPIO: undefined
            };
        };

        this.setGPIO = function (gpio) {
            $scope.gpio_add_model.GPIO = gpio;
        };

        this.setPage = function (pageName) {
            $scope.currentPage = pageName;
        };

        $scope.dripnodes = [];
        $scope.watercost = 0.0021;
        $scope.getWaterUsageStats = function () {
            $http.get('/nodes/usage')
                .then(function (response) {
                    $scope.dripnodes = undefined;
                    $scope.dripnodes = response.data.waterUsage;
                }).then(function () {
                    $scope.calcMonthlyCost();
                })
        };

        this.getChartData = function (chartNum) {
            $http.get('/stats/' + chartNum)
                .then(function (response) {
                    $scope['chartData' + chartNum].labels = response.data.Labels;
                    $scope['chartData' + chartNum].series = response.data.Series;
                    $scope['chartData' + chartNum].data = response.data.Data;
                });
        };

        $scope.monthly_cost = 0;
        $scope.calcMonthlyCost = function () {
            $scope.monthly_cost = 0;
            angular.forEach($scope.dripnodes, function (value, key) {
                $scope.monthly_cost += value.Total30Days * $scope.watercost;
            });
        };


        this.loadStatsData = function () {
            $scope.beatheart = true;
            Chart.defaults.global.defaultFontColor = "#fff";
            this.getChartData(1);
            this.getChartData(2);
            // this.getChartData(3)
            this.getChartData(4);
            $scope.beatheart = false;
        };

        $scope.stationModel = {};
        this.addStationButton = function () {
            $scope.showAddStationButton = false;
            $scope.stationModel = {
                tempID: -1,
                GPIO: $scope.availableGpios[0].GPIO,
                Notes: "Created: " + new Date(),
                new: true
            };
            $scope.stations.unshift($scope.stationModel);

        };
        $scope.showAddStationButton = true;

        this.setEditingStationInfo = function (station) {
            $scope.stationModel = station;
            $scope.stationModel.tempID = station.ID;
            $scope.stationModel.GPIO = 0 + station.GPIO;
        };

        this.submitEditStation = function (newStation) {
            $scope.showAddStationButton = true;
            $scope.stationModel.ID = $scope.stationModel.tempID;
            $http.post('/station/edit', $scope.stationModel)
                .then(function (response) {
                    $scope.stations = response.data.stations;
                    $scope.stationModel = undefined;
                });
            $scope.getGpios();
        };

        this.submitDeleteStation = function (stationID) {
            $scope.showAddStationButton = true;
            $http.post('/station/delete', {
                ID: stationID
            })
                .then(function (response) {
                    $scope.stations = response.data.stations;
                    $scope.stationModel = undefined;
                });
            $scope.getGpios();
        };

        this.submitAddStation = function () {
            $scope.showAddStationButton = true;
            $http.post('/station/add', $scope.stationModel)
                .then(function (response) {
                    $scope.stations = response.data.stations;
                    $scope.stationModel = undefined;
                });
            $scope.getGpios();

        };

        $scope.scheduleModel = {};
        this.addScheduleButton = function () {
            $scope.scheduleModel = undefined;
            var d = new Date();
            var year = d.getFullYear();
            var month = d.getMonth();
            var day = d.getDate();
            var endDate = new Date(year + 10, month, day);
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
                StartDate: new Date(),
                EndDate: endDate,
                Duration: 0,
                timepicker: d
            };
            $scope.schedule.unshift($scope.scheduleModel);
        };

        this.addScheduleButtonFromCalendar = function (startTime) {
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

            $scope.schedule.unshift($scope.scheduleModel);
        };

        this.submitEditSchedule = function () {
            $http.post('/schedule/edit', $scope.scheduleModel)
                .then(function (response) {
                    $scope.schedule = response.data.stationSchedules;
                });
            $scope.scheduleModel = {};
            $scope.scheduleModel = undefined;
            this.refresh();
        };

        this.loadLogs = function () {
            $http.get('/logs/all')
                .then(function (response) {
                    $scope.logs = response.data.logs;
                });
        };

        this.submitAddSchedule = function () {
            this.convertScheduleBoolToInt();
            $http.post('/schedule/add', $scope.scheduleModel)
                .success(function (response) { });
            // cleanup
            $scope.scheduleModel = {};
            $scope.scheduleModel = undefined;
            this.refresh();
        };

        this.mapModelForSchedEdit = function (currentModel) {
            $scope.scheduleModel = currentModel;
        };

        this.mapModelForSchedEditFromCalClick = function (id) {
            $scope.currentPage = 'calendar';
            var sch = $filter('filter')($scope.schedule, {
                id: id
            })[0];
            $scope.scheduleModel = sch;
        };

        this.submitDeleteSchedule = function (schedule_id) {
            $http.post('/schedule/delete', {
                ID: schedule_id
            })
                .then(function (response) { });
            $scope.scheduleModel = undefined;
            this.refresh();
        };

        $scope.startTestRun = function () {
            $scope.stations.forEach(function (station) {
                var r = {
                    StationID: station.ID,
                    Duration: 90
                };
                $http.post('/station/run', r);
            })
        };

        $scope.singleRunModel = {};
        this.submitSingleRun = function () {
            $http.post('/station/run', $scope.singleRunModel)
                .then(function (response) {
                    // console.log($scope.singleRunModel);
                });
            $scope.singleRunModel = undefined;
            $scope.singleRunModel = {};
        };

        this.refresh = function () {
            // this.getSchedule();
            this.getCalEvents();
            this.loadStations();
            $scope.getGpios();
            this.loadStatsData();
            $scope.getWaterUsageStats();
            this.loadSettings();
            this.loadWeather();
        };
        this.loadStations = function () {
            $http.get('/station/all')
                .then(function (response) {
                    $scope.stations = response.data.stations;
                });
        };

        this.loadSettings = function () {
            $http.get('/settings/load')
                .then(function (response) {
                    $scope.settingsData = response.data.data;
                });
        };

        this.loadWeather = function () {
            $scope.beatheart = true;
            $http.get('/weather')
                .then(function (response) {
                    //                    $scope.weatherData = response.data;
                    //                    $scope.weatherData.current.sys.sunrise_t = moment(data.current.sys.sunrise * 1000).fromNow();
                    //                    $scope.weatherData.current.sys.sunset_t = moment(data.current.sys.sunset * 1000).fromNow();

                });
        };

        this.loadHistory = function (station) {
            var query = '?station=' + station + '&earliest=-168';
            $http.get('/history' + query)
                .then(function (response) {
                    $scope.history = response.data.history;
                });
        };

        this.prettyTime = function (uglyTime) {
            if (uglyTime !== undefined && uglyTime !== null) {
                // console.log(uglyTime)
                var pt = moment(uglyTime).calendar();
                return pt;
            } else {
                return "Never";
            }
        };

        $scope.addDays = function (startDate, numberOfDays) {
            var returnDate = new Date(
                startDate.getFullYear(),
                startDate.getMonth(),
                startDate.getDate() + numberOfDays,
                startDate.getHours(),
                startDate.getMinutes(),
                startDate.getSeconds());
            return returnDate;
        };

        $scope.waterNodeEntries = [];
        $scope.waterNodeModel = {};
        $scope.mapWaterNodeModel = function (node) {
            $scope.waterNodeModel = node;
        };

        this.getWaterNodes = function () {
            $http.get('/nodes')
                .then(function (response) {
                    $scope.waterNodeEntries = response.data.nodes;
                });
        };

        this.submitEditNode = function () {
            $http.post('/nodes/edit', $scope.waterNodeModel)
                .then(function (response) {
                    $scope.waterNodeEntries = response.data.nodes
                }).then(function () {
                    $scope.getWaterUsageStats();
                }).then(function () {
                    $scope.calcMonthlyCost();
                });
            // cleanup



            $scope.waterNodeModel = undefined;
            $scope.waterNodeModel = {};
        };

        $scope.waterNodeEntries = [];
        this.addWaterNodeButton = function () {
            $scope.waterNodeModel = undefined;
            $scope.waterNodeModel = {
                StationID: 'Select Station ID',
                GPH: 0,
                Count: 0,
                new: true
            };
            if ($scope.waterNodeEntries.length > 0) {
                $scope.waterNodeEntries.unshift($scope.waterNodeModel);
            } else {
                $scope.waterNodeEntries = [$scope.waterNodeModel];
            }
        };

        this.submitAddNode = function () {
            $http.post('/nodes/add', $scope.waterNodeModel)
                .then(function (response) {
                    $scope.waterNodeEntries = response.data.nodes;
                }).then(function () {
                    $scope.getWaterUsageStats();
                }).then(function () {
                    $scope.calcMonthlyCost();
                });
            // cleanup
            $scope.waterNodeModel = undefined;
            $scope.waterNodeModel = {};
        }

        this.submitDeleteNode = function (id) {
            $http.post('/nodes/delete', {
                ID: id
            })
                .then(function (response) {
                    $scope.waterNodeEntries = response.data.nodes;
                }).then(function () {
                    $scope.getWaterUsageStats();
                }).then(function () {
                    $scope.calcMonthlyCost();
                });
            // cleanup
            $scope.waterNodeModel = undefined;
            $scope.waterNodeModel = {};
        }

        $scope.loader = this.autoLoader;
        // START CAL

        var date = new Date();
        var d = date.getDate();
        var m = date.getMonth();
        var y = date.getFullYear();

        $scope.calOptions = {
            calendarMode: "week"
        };

        function padDigits(number, digits) {
            return Array(Math.max(digits - String(number).length + 1, 0)).join(0) + number;
        }

        $scope.loadCalEvents = function () {
            $scope.eventSource = [];
            for (i = 0; i < $scope.schedule.length; i++) {
                var entry = $scope.schedule[i];
                var startTime = padDigits(entry.StartTime, 4);
                var curDate = new Date();
                var fstartDate = new Date();
                var startDate = new Date(
                    fstartDate.getFullYear(),
                    fstartDate.getMonth(),
                    fstartDate.getDate(),
                    startTime.slice(0, 2),
                    startTime.slice(2, 4),
                    0
                );
                var endDate = new Date(entry.EndDate);

                var dateDiff = Math.abs($scope.addDays(curDate, 31) - $scope.addDays(curDate, -31));
                var diffDays = Math.ceil(dateDiff / (1000 * 3600 * 24));

                while (diffDays > 0) {
                    var newRealDate = $scope.addDays(startDate, diffDays);
                    var shouldShow = false;
                    if (newRealDate < endDate && newRealDate > startDate) {
                        if (newRealDate.getDay() === 0 && entry.Sunday) {
                            shouldShow = true;
                        } else if (newRealDate.getDay() == 1 && entry.Monday) {
                            shouldShow = true;
                        } else if (newRealDate.getDay() == 2 && entry.Tuesday) {
                            shouldShow = true;
                        } else if (newRealDate.getDay() == 3 && entry.Wednesday) {
                            shouldShow = true;
                        } else if (newRealDate.getDay() == 4 && entry.Thursday) {
                            shouldShow = true;
                        } else if (newRealDate.getDay() == 5 && entry.Friday) {
                            shouldShow = true;
                        } else if (newRealDate.getDay() == 6 && entry.Saturday) {
                            shouldShow = true;
                        }

                        if (shouldShow) {
                            var newEntry = {
                                stationCalID: entry.ID,
                                title: "Station Run: " + entry.ID,
                                allDay: false,
                                startTime: new Date(
                                    newRealDate.getFullYear(),
                                    newRealDate.getMonth(),
                                    newRealDate.getDate(),
                                    newRealDate.getHours(),
                                    newRealDate.getMinutes(),
                                    newRealDate.getSeconds()
                                ),
                                endTime: new Date(
                                    newRealDate.getFullYear(),
                                    newRealDate.getMonth(),
                                    newRealDate.getDate(),
                                    newRealDate.getHours(),
                                    newRealDate.getMinutes(),
                                    newRealDate.getSeconds() + entry.Duration
                                )
                            };
                            $scope.eventSource.push(newEntry);
                        }
                    }
                    diffDays--;
                }
            }
            $scope.$broadcast('eventSourceChanged', $scope.eventSource);
        };

        $scope.mode = "week";
        $scope.changeMode = function (mode) {
            $scope.mode = mode;
        };


        $scope.availableGpios = [1];
        $scope.getGpios = function () {
            $http.get('/gpio/available')
                .then(function (response) {
                    $scope.availableGpios = response.data.gpios;
                });
        };

        $scope.commonGpio = "unset";
        $scope.commonGpioModel = undefined;
        $scope.setGpioCommon = function (gpio) {
            $http.post("/gpio/common/set", { GPIO: gpio })
                .then(function (response) {
                    // console.log({ GPIO: gpio })
                    $scope.commonGpio = response.data.gpio.GPIO
                    $scope.commonGpioModel = undefined;
                });
            $scope.getGpios();
        };

        $scope.getGpioCommon = function () {
            $http.get("/gpio/common")
                .then(function (response) {
                    // console.log(response.data)
                    $scope.commonGpio = response.data.gpio.GPIO
                });
            $scope.getGpios();
        };

        this.autoLoader = function () {
            this.getCalEvents();
            this.loadStations();
            $scope.getWaterUsageStats();
            this.getWaterNodes();
            $scope.getRunStatus();
            this.loadStatsData();
            this.loadHistory();
            $scope.calcMonthlyCost();
            $scope.getGpios();
            $scope.getGpioCommon();
            if ($cookies.get('lastTab') !== undefined) {
                $scope.currentPage = $cookies.get('lastTab');
            }
        };
        this.autoLoader();
    });

    // END CAL

})();