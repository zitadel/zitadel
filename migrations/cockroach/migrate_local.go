//+build ignore

package migrations

//go:generate flyway -url=jdbc:postgresql://localhost:26257/defaultdb -user=root -password= -locations=filesystem:./ -placeholders.eventstorepassword=NULL -placeholders.managementpassword=NULL -placeholders.adminapipassword=NULL -placeholders.authpassword=NULL -placeholders.notificationpassword=NULL -placeholders.authzpassword=NULL migrate
