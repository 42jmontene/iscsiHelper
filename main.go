package main

import (
	"fmt"
	"os"
	"log"
	"strconv"
)

const (
	FULL_CREATION int = 1
	BIND_LDAP int = 2
	RESET_TARGET int = 3
	DELETE_HOME int = 4
	HARD_RESET int = 5
	CREATE_HOME int = 6
	ldapAddress string = "ldap-master.42.fr:389"
)

type info struct {
	ldapNb int
	storageNb int
}

func initInfos () (info) {
	infos := info{0, 0}
	return infos
}

func printMenu(infos info) ([6]int) {
	var menu [6]int
	i := 1
	fmt.Printf("Choose your option:\n\n")
	fmt.Printf("\nstorage nb %d ldapnb %d\n", infos.storageNb, infos.ldapNb)
	if infos.storageNb == 0 && infos.ldapNb <= 0  {
		fmt.Printf("%d - FULL CREATION\n", i)
		menu[i] = FULL_CREATION
		i++
	}
	if infos.storageNb != 0 && infos.ldapNb <= 0 {
		fmt.Printf("%d - BIND LDAP\n", i)
		menu[i] = BIND_LDAP
		i++
	}
	if infos.storageNb != 0 && infos.ldapNb != 0 {
		fmt.Printf("%d - RESET TARGET\n", i)
		menu[i] = RESET_TARGET
		i++
		fmt.Printf("%d - DELETE HOME\n", i)
		menu[i] = DELETE_HOME
		i++
		fmt.Printf("%d - HARD RESET\n", i)
		menu[i] = HARD_RESET
		i++
	}
	if infos.storageNb == 0 && infos.ldapNb != 0 {
		fmt.Printf("%d - CREATE HOME\n", i)
		menu[i] = CREATE_HOME
		i++
	}
	return menu
}

func waitForInput(menu [6]int) (int) {
	input := "0"
	confirm := ""
	i := 1
	for menu[i] != 0 {
		i++
	}
	inputNb, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	for inputNb > i - 1 || inputNb < 1 {
		fmt.Scanf("%d", &inputNb)
	}
	for confirm != "n" && confirm != "y" {
		fmt.Printf("Choix %d, etes vous sur ? (y/n) ", inputNb)
		fmt.Scanf("%s\n", &confirm)
	}
	if confirm == "n" {
		os.Exit(1)
	}
	return inputNb
}

func performActions(infos info, action int, login string) {
	if action == FULL_CREATION {
		if infos.ldapNb == 0 {
			fmt.Printf("Sur quel storage faut-il créer le home ? (1/2)")
			for infos.storageNb < 1 && infos.storageNb > 2{
				fmt.Scanf("%d", &infos.storageNb)
			}
		} else {
			infos.ldapNb = -infos.ldapNb
			infos.storageNb = infos.ldapNb
		}
		bindLdap(login, infos.storageNb)
		createHome(login, infos.ldapNb)
		resetTarget(login, infos.storageNb)
	} else if action == BIND_LDAP {
		bindLdap(login, infos.storageNb)
	} else if action == RESET_TARGET {
		resetTarget(login, infos.storageNb)
	} else if action == DELETE_HOME {
		deleteHome(login, infos.storageNb)
	} else if action == HARD_RESET {
		deleteHome(login, infos.storageNb)
		createHome(login, infos.ldapNb)
		resetTarget(login, infos.storageNb)
	} else if action == CREATE_HOME {
		createHome(login, infos.ldapNb)
		resetTarget(login, infos.storageNb)
	} else {
	}
}

func main () {
	if len(os.Args) != 2 {
		log.Fatal("usage : ./iscsiHelper [login]")
	}
	fmt.Printf("\nWelcome to ISCSI Manager !\n")
	fmt.Printf("-----------------------------------\n\n")
	login := os.Args[1]
	fmt.Printf("\nSearching home infos for %s...", login)
	infos := initInfos()
	infos = ldapGetInfos(infos, login)
	infos = iscsiGetInfos(infos, login)
	fmt.Printf("Informations found !\n\n")
	menu := printMenu(infos)
	inputNb := waitForInput(menu)
	performActions(infos, menu[inputNb], login)
//	fmt.Println(infos)
//	displayOptions()
//	makeActions()
}


	/* TEST DES CONDITIONS :	
			EST CE QUE L ETUDIANT A UN LDAP ?					(ldapOn)
				OUI : VALIDATION
				NON : ARRET, ON NE PEUT RIEN FAIRE

			EST CE QUE L ETUDIANT A UN BIND SS SUR LE LDAP ?			(bindOn)
				OUI : ON A SON STORAGE DISPO
				NON : ON NOTE QU IL A PAS DE BIND LDAP --> SCRIPT A REMI

			EST CE QUE L ETUDIANT A UN STORAGE ?					(storageOn)
				OUI : ON NOTE LE STORAGE
				NON : ON NOTE QU IL A PAS DE STORAGE
	*/

	/* AFFICHAGE DU MENU SELON LES OPTIONS DISPO :

	FAIT		- BIND LDAP			| ldapOn && storageOn && bindOff
	FAIT		- CREATE ALL			| ldapOn && storageOff && bindOff
	FAIT		- RESET TARGET			| ldapOn && storageOn && bindOn
	FAIT		- CREATE HOME			| ldapOn && storageOff && bindOn
	FAIT		- DELETE HOME			| ldapOn && storageOn && bindOn
	FAIT		- HARD RESET			| ldapOn && storageOn && bindOn

	*/

	/* ACTIONS :

			- BIND LDAP			| ldap_ss(which_storage(login))

			- CREATE ALL			| ldap_ss(choose_storage)
							| create-iscsi-home
							| kill-iscsi-sessions
							| reset-target

			- RESET TARGET			| kill-iscsi-sessions
							| reset-target

			- CREATE HOME			| create-iscsi-home

			- DELETE HOME			| delete-iscsi-home

			- HARD RESET			| delete-iscsi-home
							| create-iscsi-home
							| kill-iscsi-sessions
							| reset-target
	*/
