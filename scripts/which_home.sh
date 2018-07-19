#!/bin/sh

kinit -kt ~/homemaker.keytab homemaker@42.FR
if curl --negotiate -s -u : https://student-storage-1.42.fr/search/$1 | grep -q "success"
then
	echo "1"
elif curl --negotiate -s -u : https://student-storage-2.42.fr/search/$1 | grep -q "success"
then
	echo "2"
else
	echo "0"
fi
