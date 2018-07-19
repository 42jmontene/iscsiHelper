#!/usr/bin/env python
import subprocess
import argparse
import sys
from jinja2 import Template

LDAPURL="ldap://ldap-master.42.fr"

def get_dn(login):
    out = subprocess.check_output(["ldapsearch", "-H", LDAPURL, "-LLL", "uid={}".format(login), "dn"], stderr=open("/dev/null")).rstrip()
    if len(out) == 0:
        raise Exception("No such login {}".format(login))
    return out.split(" ")[1]

def render_ldif_homedir(login, homedir):
    dn = get_dn(login)
    tpl = Template("""dn: {{ dn }}
changetype: modify
replace: homeDirectory
homeDirectory: {{ homedir }}
""")
    ldif = tpl.render(homedir=homedir, dn=dn)
    return ldif

def get_status(login):
    status_raw = subprocess.check_output("ldapsearch -H {} -LLL uid={} status 2>/dev/null | perl -p00e 's/\r?\n //g'".format(LDAPURL, login), shell=True).rstrip()
    if len(status_raw) == 0:
        raise Exception("No such login {}".format(login))
    status = {}
    try:
        status_raw = status_raw.split(" ")[2]
    except:
        return {}
    for flag in status_raw.split(','):
        if "=" in flag:
            k, v = flag.split('=')
            status[k] = v
        else:
            status[flag] = None

    return status

def save_status(status, login):
    status_list = []
    for (k, v) in status.iteritems():
        if v is None:
            status_list.append(k)
        else:
            status_list.append("{}={}".format(k, v))

    ldif = render_ldif_status(login, ",".join(status_list))
    ldapmodify(ldif)

def render_ldif_status(login, status):
    dn = get_dn(login)
    ldif ="""dn: {}
changetype: modify
replace: status
status: {}
""".format(dn, status)
    return ldif

def ldapmodify(ldif):
    p = subprocess.Popen("ldapmodify -H {}".format(LDAPURL), shell=True, stdout=subprocess.PIPE, stdin=subprocess.PIPE, stderr=subprocess.PIPE)
    out, err = p.communicate(input=ldif)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("login")
    parser.add_argument("onoff", choices=["on", "off"])

    args = parser.parse_args()

    status = get_status(args.login)
    status['iscsi-portal'] = '10.51.1.229'
    status['iscsi-iqn'] = 'iqn.2016-08.fr.42.homedirs:{}'.format(args.login)
    status['iscsi-ws'] = 'https://student-storage-1.42.fr'
    status['iscsi-warn'] = 'no'
    status['iscsi-voluntary'] = 'no'
    status['migration'] = 'no'
    if args.onoff == 'on':
        #TMP - UNCOMMENT THIS AND REMOVE THE HOMEDIR PART
        print "Switching {} to iscsi mode".format(args.login)
        status['home'] = 'iscsi'
        #print "Creating iscsi home for {}".format(args.login)
        #out = subprocess.check_output("/etc/scripts/iscsi-ws/create_iscsi_home.sh {}".format(args.login), shell=True)
        ldif = render_ldif_homedir(args.login, "/Users/{}".format(args.login))
        print "Changing homedir for {}".format(args.login)
        ldapmodify(ldif)
    else:
        print "Switching {} to nfs mode".format(args.login)
        status['home'] = 'nfs-direct'
    save_status(status, args.login)
