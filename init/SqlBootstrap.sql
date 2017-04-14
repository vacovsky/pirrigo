
# GPIO pins
use pirri;
insert into gpio_pins (gpio) values (4),(5),(6),(12),(13),(16),(18),(20),(21),(22),(23),(24),(25),(26),(27);
update gpio_pins set notes='common' where gpio=21;

INSERT INTO pirri.station_schedules (sunday, monday, tuesday, wednesday, thursday, friday, saturday, station_id, start_time, duration) values (true, true, true, true, true,true, true, 1, 1235, 60);