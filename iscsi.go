
package main

import (
	"fmt"
	"strconv"
	"os/exec"
//	"log"
)

func bindLdap(login string, storage int) {
	dir := "./scripts/ldap_ss" + strconv.Itoa(storage) + ".py"
	fmt.Printf("Ici on executera ceci :\nkinit -kt ~/admin.keytab admin@42.FR\n%s %s on\n", dir, login)
	_ = exec.Command("kinit", "-kt", "~/admin.keytab", "admin@42.FR", "&&", dir, login, "on").Run()
//	_ = exec.Command(dir, login, "on").Run()
/*	if err != nil {
		log.Fatal("err %s", err)
	}*/
}

func createHome(login string, ldap int) {
	fmt.Printf("Ici fera le curl de creation de home sur le storage %d\n", ldap)
}

func resetTarget(login string, storage int) {
	fmt.Println("Ici on fera le curl de kill sessions et celui de reset target")
}

func deleteHome(login string, storage int) {
	fmt.Println("Ici on fera le curl de suppression de home")
}
