// TODO: Show disabled slots
// TODO: option to start from a different day of week

;(function($, window, document, _, moment, undefined){
	'use strict';

	if( typeof Object.create !== 'function'){
		Object.create = function(obj){
			function F(){};
			F.prototype = obj;
			return new F();
		}
	}

	var widgetProperties = {
		dataKey : 'easycal',
		classnames : {
			widget : 'easycal-widget',
			mainTable : 'easycal',
			headerTable : 'ec-head-table',
			timeGridContainer : 'ec-time-grid-container',
			timeGridTable : 'ec-time-grid-table',

			timeLabel : 'ec-time',
			dayColumn : 'ec-slot-col',
			timeSlot : 'ec-slot',
			minorSlot : 'ec-minor-slot',
			eventContainer : 'ec-event',
			timeRange : 'ec-time-range',
			eventTitle : 'ec-event-title'
		},

		format : {
			dateLong : 'YYYY-MM-DD HH:mm:ss',
			dateShort : 'YYYY-MM-DD',
			timeLong : 'HH:mm:ss',
			timeShort : 'HH:mm'
		}
	};

	var classes = widgetProperties.classnames;
	var format = widgetProperties.format;

	var Easycal = {

		init : function(elem, options){
			var self = this;

			self.elem = elem;
			self.$elem = $( elem );
			self.$elem.addClass(classes.widget);

			self.options = options;

			self.momStartDate = moment(this.options.startDate, format.dateShort);
			self.momMinTime = moment(this.options.minTime, format.timeLong);
			self.momMaxTime = moment(this.options.maxTime, format.timeLong);

			self._launch();
			self._attachEventHandlers();			
		},

		_launch : function(){
			this._display();
			this._calculateDimensions();
			this._inflateMinorSlots();	
			this._showEvents();
		},

		_calculateDimensions : function(){
			var $timeGridTable = this.$elem.find('table.' + classes.timeGridTable);
			var $cols = $timeGridTable.find('td.' + classes.dayColumn);

			self.colWidth = $cols.eq(0).width();
			self.slotHeight = $cols.find('.' + classes.timeSlot).eq(0).height();
		},

		_attachEventHandlers : function(){
			var self = this;
			self.$elem.find('table.' + classes.timeGridTable + ' .' + classes.timeSlot).on('click', function(ev){
				if($(ev.target).closest('.' + classes.eventContainer).length || $(ev.target).hasClass(classes.eventContainer)){
					if(typeof self.options.eventClick === 'function'){
						var $eventContainer = (($(ev.target).closest('.' + classes.eventContainer).length) ? $(ev.target).closest('.' + classes.eventContainer) : $(ev.target));
						var eventId = $eventContainer.attr('data-event-id');
						
						self.options.eventClick.apply(self, [eventId]);
					}
				}else if(typeof self.options.dayClick === 'function'){
					var slotStartTime = $(this).attr('data-time');
					self.options.dayClick.apply(self, [$(this), slotStartTime]);
				}
			});
		},

		_detachEventHandlers : function(){
			this.$elem.find('table.' + classes.timeGridTable + ' .' + classes.timeSlot).off();
		},

		_display : function(){
			var html = this.renderHTML();
			this.$elem.html(html);
		},

		refresh : function(events){
			this._detachEventHandlers();
			this._clearEvents();
			
			if(events){
				this.options.events = events;
			}

			this._calculateDimensions();
			this._inflateMinorSlots();
			this._showEvents();
		},

		destroy : function(){
			this.$elem.removeClass(classes.widget);
			$.data(this.elem, widgetProperties.dataKey, null);
			this._detachEventHandlers();
			this.$elem.children().hide().remove();
		},

		_mapEventsByDate : function(){
			var res = {};
			var events = this.options.events;
			var date = this.momStartDate.clone().isoWeekday(1);
			
			for(var i = 0 ; i < 7 ; i++){
				var dateStr = date.format(format.dateShort);
				var filteredEvents = _.filter(events, function(event){
					if((moment(event.start, format.dateLong).format(format.dateShort)) === dateStr){
						return true;
					}
				});
				res[dateStr] = filteredEvents;
				date.add(1, 'd');
			}
			return res;
		},

		_showEvents : function(){
			var self = this;
			var events = this.options.events;

			var $timeGridTable = this.$elem.find('table.' + classes.timeGridTable);
			var $cols = $timeGridTable.find('td.' + classes.dayColumn);

			var eventDateMap = this._mapEventsByDate(), $col = null, $slots = null, $slot = null, schedule = null, slotTime = null;

			_.each($cols, function(col, i){
				$col = $(col); 
				var colDate = $col.attr('data-date');
				var dayEvents = eventDateMap[colDate];
				
				if(dayEvents.length){
					schedule = self.getDaySchedule(dayEvents);

					if(self.options.timeGranularity === self.options.slotDuration){
						$slots = $col.find('.' + classes.timeSlot);
					}else{
						$slots = $col.find('.' + classes.timeSlot + ' .' + classes.minorSlot);
					}
					
					_.each($slots, function(slot, i){
						$slot = $(slot);
						slotTime = $slot.attr('data-time');
						var scheduleForSlot = schedule[slotTime];
						if(scheduleForSlot.length > 1){
							// Events overlap for a time slot
							$(slot).css({
								'background-color' : self.options.overlapColor,
								color : self.options.overlapTextColor
							}).html(self.renderSlotHTML(scheduleForSlot));

						}else if(scheduleForSlot.length){
							$slot.css({
								'background-color' : scheduleForSlot[0].backgroundColor,
								color : scheduleForSlot[0].textColor
							}).addClass(classes.eventContainer).attr('data-event-id', scheduleForSlot[0].id);

							var slotStartTime = moment(colDate + ' ' + slotTime, format.dateLong);
							var eventStartTime = moment(scheduleForSlot[0].start, format.dateLong);
							if(slotStartTime.isSame(eventStartTime)){
								$slot.html(self.renderSlotHTML(scheduleForSlot));
								if(!self.isSpanMultipleSlots(scheduleForSlot)){
									$slot.find('.' + classes.eventTitle).css({
										width: (self.colWidth - 6),
										'white-space': 'nowrap',
										overflow: 'hidden',
										'text-overflow': 'ellipsis'
									});
								}
							}else{
								$(slot).css({
									'border-top' : '1px solid ' + scheduleForSlot[0].backgroundColor
								});
							}
							if($slot.hasClass(classes.minorSlot)){
								var parentSlotTime = $slot.parent('.' + classes.timeSlot).attr('data-time');
								var momParentSlotTime = moment(colDate + ' ' + parentSlotTime, format.dateLong);
								if(momParentSlotTime.isAfter(eventStartTime)){
									$slot.parent('.' + classes.timeSlot).css({
										'border-top' : '1px solid ' + scheduleForSlot[0].backgroundColor
									});
								}
							}
						}
						
					});
				}
			});
		},

		/*
		 * Redraws the timeGridTable without the events
		 */
		_clearEvents : function(){
			var html = this._renderTimeGridHTML();
			this.$elem.find('.' + classes.timeGridContainer).children().hide().remove();
			this.$elem.find('.' + classes.timeGridContainer).html(html);
		},

		/*
		 * Accepts only a single schedule and returns true if
		 * it spans multiple slotDurations
		 */
		isSpanMultipleSlots : function(schedule){
			if(schedule.length === 1){
				var startTime = moment(schedule[0].start, format.dateLong).add(this.options.slotDuration, 'm');
				var endTime = moment(schedule[0].end, format.dateLong);
				if(!startTime.isSame(endTime)){
					return true;
				}else{
					return false;
				}
			}
		},

		renderSlotHTML : function(scheduleList){
			var html = '';
			if(scheduleList.length > 1){
				html += '<div>' + this.options.overlapTitle + '</div>';
			}else if(scheduleList.length){
				var schedule = scheduleList[0];
				var startTime = moment(schedule.start, format.dateLong).format(format.timeShort);
				var endTime = moment(schedule.end, format.dateLong).format(format.timeShort); 
				html += '' +
					'<div style="padding-top: 4px;">' +
						'<div class="' + classes.timeRange + '">' + startTime + ' - ' + endTime + '</div>' +
						'<div class="' + classes.eventTitle + '">' + schedule.title + '</div>' +
					'</div>';
			}
			return html;
		},

		getDaySchedule : function(dayEvents){
			var date = moment(dayEvents[0].start, format.dateLong).format(format.dateShort);
			var minTime = moment(date + ' ' + this.options.minTime, format.dateLong);
			var maxTime = moment(date + ' ' + this.options.maxTime, format.dateLong);
			var time = minTime.clone();

			var schedule = {};

			var begining = null, end = null;
			for(;time.isBefore(maxTime);){
				begining = time.clone();
				end = begining.clone().add(this.options.timeGranularity, 'm');

				var slotEvents = _.filter(dayEvents, function(event){
					
					var eventStart = moment(event.start, format.dateLong);
					var eventEnd = moment(event.end, format.dateLong);

					if(eventStart.isBefore(end) && eventEnd.isAfter(begining)){
						return true;
					}
				});

				schedule[time.format(format.timeLong)] = slotEvents;
				time.add(this.options.timeGranularity,'m');
			}
			return schedule;
		},

		renderHTML : function(){
			return '<table border="0" cellpadding="0" cellspacing="0" class="easycal">' +
						'<thead>' +
							'<tr>' +
								'<td>' +
									(this.renderHeadHTML()) +
								'</td>' +
							'</tr>' +
						'</thead>' +
						'<tbody>' +
							'<tr>' +
								'<td class="' + classes.timeGridContainer + '">' +
									(this._renderTimeGridHTML()) +
								'</td>' +
							'</tr>' +
						'</tobdy' +
				'</table>';
		},

		renderHeadHTML : function(){
			var date = moment(this.options.startDate, format.dateShort);
			date.isoWeekday(1);

			var html = '<table border="0" cellspacing="0" cellpadding="0" class="ec-head-table"><tbody><tr>';
			for(var i = 0 ; i < 8 ; i++){
				var cellContent = '';
				if(i !== 0){
					cellContent = date.format(this.options.columnDateFormat);
					html += '<td class="ec-day-header">' + cellContent + '</td>';
					date.add(1, 'd');
				}else{
					html += '<td></td>';
				}
			}
			return html + '</tr></tbody></table>';
		},

		_renderTimeGridHTML : function(){
			var minTime = this.momMinTime;
			var maxTime = this.momMaxTime;
			var time = minTime.clone();

			var date = moment(this.options.startDate, format.dateShort);
			date.isoWeekday(1);
			
			var html = '<table border="0" cellspacing="0" cellpadding="0" class="ec-time-grid-table"><tbody><tr>';

			var cellContent = null, timeTag = null, colDate = null;

			for(var i = 0 ; i < 8 ; i++){
				if(i===0){
					html += '<td>';
				}else{
					colDate = date.format(format.dateShort);
					html += '<td class="ec-slot-col" data-date="' + colDate + '">';
					date.add(1, 'd');
				}

				for(;time.isBefore(maxTime);){
					if(i === 0){
						cellContent = time.format(this.options.timeFormat);
						html += '<div class="table-cell ' + classes.timeLabel + '">' + cellContent + '</div>';
					}else{
						timeTag = time.format(format.timeLong);
						html += '<div class="table-cell ' + classes.timeSlot + '" data-time="' + timeTag + '">';
						html += this._getMinorSlotForCell(time);
						html += '</div>';
					}
					time.add(this.options.slotDuration,'m');
				}

				html += '</td>';
				time = minTime.clone();
			}

			return html + '</tr></tbody></table>';
		},

		_getMinorSlotForCell : function(momTime){
			var time = momTime.clone();
			var html = '';
			if(this.options.timeGranularity < this.options.slotDuration){
				for(var i = 0 ; i < (this.options.slotDuration/this.options.timeGranularity); i++ ){
					html += '<div class="' + classes.minorSlot + '" data-time="' + time.format(format.timeLong) + '"></div>';
					time.add(this.options.timeGranularity,'m');
				}
			}
			return html;
		},

		_inflateMinorSlots : function(){
			var granularityLevel = this.options.slotDuration/this.options.timeGranularity;
			var minorSlotHeight = self.slotHeight/granularityLevel;
			this.$elem.find('table.' + classes.timeGridTable + ' td .' + classes.minorSlot).css({
				height: minorSlotHeight,
				'max-height' : minorSlotHeight
			});
		}

	};

	$.fn.easycal = function(options){

		var mergedOptions = $.extend({}, $.fn.easycal.defaults, options);
		var args = Array.prototype.slice.call(arguments, 1);

		if(typeof options === 'undefined' || typeof options === 'object'){
			return this.each(function(){
				if(!$.data(this, widgetProperties.dataKey)){
					var easycal = Object.create(Easycal);
					$.data(this, widgetProperties.dataKey, easycal);
					easycal.init(this, mergedOptions);
				}
			});
		} else if(typeof options === 'string' && options[0] !== '_'){
			var returns;
			this.each(function(){
				var instance = $.data(this, widgetProperties.dataKey);
				if(Easycal.isPrototypeOf(instance) && typeof instance[options] === 'function'){
					returns = instance[options].apply(instance, args);
				}
			});

			return (typeof returns === 'undefined') ? this : returns;
		}
	};

	$.fn.easycal.defaults = {
		startDate : moment().format(format.dateShort),
		columnDateFormat : 'dddd, DD MMM',
		timeFormat : 'HH:mm',
		minTime : '08:00:00',
		maxTime : '19:00:00',
		slotDuration : 30, //in mins
		timeGranularity : 15, // in mins
		dayClick : null,
		eventClick : null,
		events : [],
		overlapColor : '#FF0',
		overlapTextColor : '#000',
		overlapTitle : 'Multiple'
	};

})(jQuery, window, document, _, moment);