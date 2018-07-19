package main

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
	"strings"
)

func connect () (*ldap.Conn) {
	l, err := ldap.Dial("tcp", ldapAddress)
	if err != nil {
		log.Fatalf("Failed to connect. %s", err)
	}
	return l
}

func search (l *ldap.Conn, login string) (string) {

	sr, err := l.Search(ldap.NewSearchRequest(
		"dc=42,dc=fr",
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		fmt.Sprintf("(uid=%s)", login),
		[]string{"status"},
		nil,
	))
	if err != nil {
		log.Fatalf("Search error. %s", err)
	}
	if len(sr.Entries) < 1 {
		log.Fatal("User does not exist on LDAP")
	}
	if len(sr.Entries) > 1 {
		log.Fatal("Too many entries returned")
	}
	return sr.Entries[0].GetAttributeValue("status")
//	fmt.Printf("%s", sr.Entries[0].GetAttributeValue("status"))
}

func storageCheck(infos info, output string) (info) {
	if strings.Contains(output, "student-storage-1") {
		infos.ldapNb = 1
	} else if strings.Contains(output, "student-storage-2") {
		infos.ldapNb = 2
	}
	if strings.Contains(output, "home=nfs-direct") {
		infos.ldapNb = -infos.ldapNb
	}
	return infos
}

func ldapGetInfos (infos info, login string) (info) {
	l := connect()
	defer l.Close()
	output := search(l, login)
	infos = storageCheck(infos, output)
	return infos
}
