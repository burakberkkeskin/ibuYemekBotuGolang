package services

var helloMessage string = "-İBU Yemek Listesi Botuna Hoş Geldiniz!\n" +
	"-Her sabah 9'da yemek listesini almak için abone olun.\n" +
	"-Abone olmak için /subscribe\n" +
	"-Abonelikten çıkmak için /unsubscribe\n" +
	"-Bugünün Listesini öğrenmek için /today\n" +
	"-Yarının Listesini öğrenmek için /tomorrow\n" +
	"-Kaynak Kod İçin /source\n" +
	"-Yardım almak için /help\n"

var lunchListToday string = GetLunchList("today")
var lunchListTomorrow string = GetLunchList("tomorrow")
